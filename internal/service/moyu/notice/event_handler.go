package notice

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

type MessageEventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender   emitter.Sender         `gone:"gone-emitter"`
	Channel  service.IChannel       `gone:"*"`
	Planet   service.IPlanet        `gone:"*"`
	Notice   service.INotice        `gone:"*"`
	Message  service.IMessageRecord `gone:"*"`
	Approval service.IApproval      `gone:"*"`
}

//go:gone
func MessageNewEventHandler() gone.Goner {
	return &MessageEventHandler{}
}

func (e *MessageEventHandler) Consume(on emitter.OnEvent) {
	on(e.handleMentioned)
	on(e.handleReferenced)
	on(e.handleMsgAck)
	on(e.handleApproveStarted)
	on(e.handleApproveAudited)
}

func (e *MessageEventHandler) handleApproveAudited(evt *entity.ApprovalAudited) error {
	e.Infof("handleApproveAudited: id=%d", evt.Id)

	notices, err := e.Notice.ListUnreadNotice(entity.NoticeTypeApproval, evt.Id)
	if err != nil {
		return err
	}

	group := collection.GroupingBy(notices,
		func(notice *entity.Notice) int64 {
			return notice.UserId
		},
		func(notice *entity.Notice) int64 {
			return notice.Id
		},
	)

	for userId, noticeIds := range group {
		err := e.Notice.MarkRead(userId, noticeIds)
		if err != nil {
			e.Errorf("MarkRead error: %v, userId = ?, noticeIds = ?", err, userId, noticeIds)
		}
	}

	return nil
}

func (e *MessageEventHandler) handleApproveStarted(evt *entity.ApprovalStarted) error {
	e.Infof("handleApproveStarted: id=%d", evt.Id)

	switch evt.ApprovalType {
	case entity.ApprovalTypeChannelJoin:
		channel, exists, err := e.Channel.GetChannelByChannelId(evt.BusinessId)
		if err != nil {
			return err
		}

		if !exists {
			return gin.NewParameterError("channel not found")
		}

		admins, err := e.Planet.ListPlanetAdmins(channel.PlanetId)
		if err != nil {
			return err
		}

		for _, admin := range admins {
			err := e.Notice.Create(admin.UserId, entity.NoticeTypeApproval, evt.Id, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *MessageEventHandler) handleMentioned(evt *entity.UserMentioned) error {
	offsets, err := e.Channel.ListMembersOffset(evt.ChannelId, evt.TargetUsers)
	if err != nil {
		return err
	}

	e.Infof("handleMentioned: mention %d user", len(evt.TargetUsers))

	for _, userId := range evt.TargetUsers {
		offset := offsets[userId]
		var read bool
		if offset >= evt.Id {
			read = true
		}

		err := e.Notice.Create(userId, entity.NoticeTypeMention, evt.Id, read)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *MessageEventHandler) handleReferenced(evt *entity.MessageReferenced) error {
	e.Infof("handleReferenced: reference=%d", evt.ReferenceId)

	record, exists, err := e.Message.GetRecord(evt.ReferenceId)
	if err != nil {
		return err
	}

	if !exists {
		e.Errorf("message not found: id=%d", evt.ReferenceId)
		return nil
	}

	if record.UserId == evt.UserId {
		return nil
	}

	offsets, err := e.Channel.ListMembersOffset(evt.ChannelId, []int64{record.UserId})
	if err != nil {
		return err
	}

	var read bool
	if offsets[record.UserId] >= evt.Id {
		read = true
	}

	return e.Notice.Create(record.UserId, entity.NoticeTypeReference, evt.Id, read)
}

func (e *MessageEventHandler) handleMsgAck(evt *wsevent.MsgAckEvent) error {
	e.Infof("handleMsgAck: userId=%d, messageId=%d", evt.UserId, evt.MsgId)

	toRead, err := e.Notice.ListUnreadIMMessages(evt.UserId, evt.MsgId, evt.ChannelId)
	if err != nil {
		return err
	}

	if len(toRead) == 0 {
		e.Infof("no notice should be mark read")
		return nil
	}

	return e.Notice.MarkRead(evt.UserId, toRead)
}
