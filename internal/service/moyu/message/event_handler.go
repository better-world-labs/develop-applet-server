package message

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
	"regexp"
	"strconv"
	"strings"
)

type EventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender   emitter.Sender         `gone:"gone-emitter"`
	Svc      service.IMessageRecord `gone:"*"`
	PChannel service.IChannel       `gone:"*"`
	Message  service.IMessageRecord `gone:"*"`
	Audit    service.IContentAudit  `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleMsgSend)
	on(e.handleHistoryEvent)
	on(e.handleMsgSaved)
	on(e.handleMessageLike)
}

func (e *EventHandler) handleMessageLike(evt *entity.MessageLikeEvent) error {
	e.Logger.Debugf("handleMessageLike: messageId=%d, userId=%d, isLike=%v\n", evt.MessageId, evt.UserId, evt.IsLike)

	return e.Svc.CreateLikeMessage(&evt.MessageLike)
}

func (e *EventHandler) listBeforeHistory(evt *wsevent.HistoryEvent) error {
	for {
		var size int
		if evt.Size > 100 {
			size = 100
		} else {
			size = evt.Size
		}
		history, err := e.Svc.ListHistory(evt.ChannelId, evt.MsgId, size, true)
		if err != nil {
			return err
		}

		l := len(history)

		err = e.Sender.Send(&wsevent.SyncHistoryEvent{
			Sid:       evt.Sid,
			ChannelId: evt.ChannelId,
			UserId:    evt.UserId,
			MsgList:   messagesToMsgs(history),
			CreatedAt: evt.CreatedAt,
		})
		if err != nil {
			return err
		}

		if l < 100 {
			return nil
		}

		evt.MsgId = history[len(history)-1].Id
		evt.Size -= l
		if evt.Size == 0 {
			return nil
		}
	}
}

func (e *EventHandler) listAfterHistory(evt *wsevent.HistoryEvent) error {
	for {
		var size int
		if evt.Size > 100 {
			size = 100
		} else {
			size = evt.Size
		}
		history, err := e.Svc.ListHistory(evt.ChannelId, evt.MsgId, size, false)
		if err != nil {
			return err
		}

		l := len(history)

		err = e.Sender.Send(&wsevent.SyncHistoryEvent{
			Sid:       evt.Sid,
			ChannelId: evt.ChannelId,
			UserId:    evt.UserId,
			MsgList:   messagesToMsgs(history),
			CreatedAt: evt.CreatedAt,
		})
		if err != nil {
			return err
		}

		if l < 100 {
			return nil
		}

		evt.MsgId = history[0].Id
		evt.Size -= l
		if evt.Size == 0 {
			return nil
		}
	}
}

func (e *EventHandler) handleHistoryEvent(evt *wsevent.HistoryEvent) error {
	e.Logger.Infof("handleHistoryEvent: msgId=%d, userId=%d", evt.MsgId, evt.UserId)
	if evt.Size > 0 {
		return e.listBeforeHistory(evt)
	} else {
		evt.Size = -evt.Size
		return e.listAfterHistory(evt)
	}
}

func (e *EventHandler) scanMessageAndContinue(m message.Content) (bool, error) {
	if text, ok := m.(*message.TextContent); ok {
		e.Logger.Infof("textMessageScan: text=%s\n", text)
		if len(text.Text) > 0 {
			result, err := e.Audit.ScanText(text.Text)
			if err != nil {
				return false, err
			}

			e.Logger.Infof("textMessageScan: result=%v\n", result)
			return result.CheckPass(), nil
		}
	}

	return true, nil
}

//	func (e *EventHandler) handleMsgSend(evt *wsevent.MsgSendEvent) error {
//		e.Logger.Infof("handleMsgSend: id=%s, seqId=%d", evt.Id, evt.SeqId)
//		m, err := e.createMessage(evt)
//		if err != nil {
//			e.Logger.Errorf("invalid message: %v", err)
//			return nil
//		}
//
//		valid, err := e.PChannel.IsChannelValid(evt.Msg.ChannelId)
//		if err != nil {
//			return err
//		}
//
//		isMember, err := e.PChannel.IsChannelMember(evt.Msg.ChannelId, evt.Msg.CreatedBy)
//		if err != nil {
//			return err
//		}
//
//		if !valid || !isMember {
//			e.Logger.Infof("channel valid=%v, isMember=%v, skip", valid, isMember)
//			return nil
//		}
//
// }
func (e *EventHandler) handleMsgSend(evt *wsevent.MsgSendEvent) error {
	e.Logger.Infof("handleMsgSend: id=%s, seqId=%d", evt.Id, evt.SeqId)

	m, err := e.sendEventToMessage(evt)
	if err != nil {
		return nil
	}

	ok, err := e.scanMessageAndContinue(m.Content)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	valid, err := e.PChannel.IsChannelValid(m.ChannelId)
	if err != nil {
		return err
	}

	isMember, err := e.PChannel.IsChannelMember(m.ChannelId, m.UserId)
	if err != nil {
		return err
	}

	if !valid || !isMember {
		e.Logger.Infof("channel valid=%v, isMember=%v, skip", valid, isMember)
		return nil
	}

	err = e.Svc.SaveMessage(m)
	if err != nil {
		return err
	}

	event, err := messageToSavedEvent(m)
	if err != nil {
		return nil
	}

	return e.Sender.Send(event)
}

func (e *EventHandler) handleMsgSaved(evt *wsevent.MsgSavedEvent) error {
	content, err := message.DecodeContent(evt.Msg.Content)
	if err != nil || content == nil {
		e.Errorf("decode message error: %v", err)
		return nil
	}

	if content.GetReference() != 0 {
		err = e.processReference(&evt.Msg, content.GetReference())
		if err != nil {
			return err
		}
	}

	switch content.(type) {
	case *message.TextContent:
		return e.handleTextMessageSaved(&evt.Msg, content.(*message.TextContent))

	case *message.ImageContent:
		return e.handleImageMessageSaved(&evt.Msg, content.(*message.ImageContent))

	case *message.EmoticonReplyContent:
		if content.GetReference() == 0 {
			e.Errorf("err EmoticonReply Message(%v) without reference", content)
		} else {
			e.processEmoticonReplyMessage(content.(*message.EmoticonReplyContent), evt.UserId)
		}

	}

	return nil
}

func (e *EventHandler) processEmoticonReplyMessage(message *message.EmoticonReplyContent, userId int64) {
	if message.EmoticonId == 0 {
		e.Errorf("msg(%v) message.EmoticonId cannot be zero", message)
	}

	e.Svc.MsgReplyRecord(message.GetReference(), message.EmoticonId, userId)
}

func (e *EventHandler) handleImageMessageSaved(evt *wsevent.Msg, message *message.ImageContent) error {
	e.Infof("handleImageMessage: id=%d, url=%s", evt.Id, message.Url)

	return e.Sender.Send(&entity.ImageMessageSavedEvent{
		ImageContent: *message,
		Origin:       evt,
	})
}

func (e *EventHandler) handleTextMessageSaved(evt *wsevent.Msg, message *message.TextContent) error {
	e.Infof("handleTextMessageSaved: id=%d, text=%s", evt.Id, message.Text)

	err := e.processMention(evt, message)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventHandler) processMention(evt *wsevent.Msg, message *message.TextContent) error {
	mentionedIds, originMessage, err := parseMentionedUserIds(message.Text)
	if err != nil {
		return err
	}

	linq.From(mentionedIds).WhereT(func(userId int64) bool {
		return userId != evt.UserId
	}).ToSlice(&mentionedIds)

	e.Infof("processMention: mention %d user", len(mentionedIds))

	if len(mentionedIds) > 0 {
		return e.Sender.Send(entity.UserMentioned{
			Id:          evt.Id,
			Sender:      evt.UserId,
			TargetUsers: mentionedIds,
			ChannelId:   evt.ChannelId,
			Msg:         evt,
			Text:        originMessage,
		})
	}

	return nil
}

func (e *EventHandler) processReference(evt *wsevent.Msg, reference int64) error {
	e.Infof("processReference: reference=%d", reference)

	record, exists, err := e.Message.GetRecord(reference)
	if err != nil {
		return err
	}

	if !exists {
		e.Errorf("message not found: id=%d", reference)
		return nil
	}

	if record.UserId == evt.UserId {
		return nil
	}

	return e.Sender.Send(&entity.MessageReferenced{
		Msg:         evt,
		ReferenceId: reference,
	})
}

func parseMentionedUserIds(text string) ([]int64, string, error) {
	set := collection.NewSet[int64]()
	r, err := regexp.Compile("@\\[\\d+]")
	if err != nil {
		return nil, "", err
	}

	mentionPlaceholders := r.FindAllString(text, -1)
	originMessage := r.ReplaceAll([]byte(text), []byte(""))
	for _, m := range mentionPlaceholders {
		userIdPrefix := strings.ReplaceAll(m, "@[", "")
		userId, err := strconv.ParseInt(userIdPrefix[:len(userIdPrefix)-1], 10, 64)
		if err != nil {
			return nil, "", err
		}

		set.Add(userId)
	}

	return set.ToSlice(), strings.TrimSpace(string(originMessage)), nil
}
