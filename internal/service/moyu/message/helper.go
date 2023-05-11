package message

import (
	"encoding/json"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

func recordsToMessages(ms []*entity.MessageRecord) []*message.Message {
	res := make([]*message.Message, 0, len(ms))

	for _, record := range ms {
		if m, err := recordToMessage(record); err == nil {
			res = append(res, m)
		}
	}

	return res
}

func messagesToMessageRecords(ms []*message.Message) []*entity.MessageRecord {
	res := make([]*entity.MessageRecord, 0, len(ms))

	for _, do := range ms {
		if record, err := messageToMessageRecord(do); err == nil {
			res = append(res, record)
		}
	}

	return res
}

func messageToMessageRecord(m *message.Message) (*entity.MessageRecord, error) {
	content, err := json.Marshal(m.Content)
	if err != nil {
		return nil, err
	}

	return &entity.MessageRecord{
		Id:        m.Id,
		SendId:    m.SendId,
		CreatedAt: m.CreatedAt,
		SendAt:    m.SendAt,
		UserId:    m.UserId,
		SeqId:     m.SeqId,
		ChannelId: m.ChannelId,
		Content:   content,
	}, nil
}

func recordToMessage(record *entity.MessageRecord) (*message.Message, error) {
	content, err := message.DecodeContent(record.Content)
	if err != nil {
		return nil, err
	}

	return &message.Message{
		Header: message.Header{
			Id:        record.Id,
			SendId:    record.SendId,
			CreatedAt: record.CreatedAt,
			SendAt:    record.SendAt,
			UserId:    record.UserId,
			SeqId:     record.SeqId,
			ChannelId: record.ChannelId,
		},
		Content: content,
	}, nil
}

func messageToSavedEvent(message *message.Message) (*wsevent.MsgSavedEvent, error) {
	msg, err := messageToMsg(message)
	if err != nil {
		return nil, err
	}

	return &wsevent.MsgSavedEvent{
		SeqId:     message.SeqId,
		UserId:    message.UserId,
		Msg:       *msg,
		CreatedAt: message.CreatedAt,
	}, nil
}

func messageToMsg(message *message.Message) (*wsevent.Msg, error) {
	content, err := json.Marshal(message.Content)
	if err != nil {
		return nil, err
	}

	return &wsevent.Msg{
		Id:        message.Id,
		CreatedAt: message.CreatedAt,
		UserId:    message.UserId,
		ChannelId: message.ChannelId,
		Content:   content,
	}, nil
}

func messagesToMsgs(messages []*message.Message) []*wsevent.Msg {
	res := make([]*wsevent.Msg, 0, len(messages))

	for _, do := range messages {
		if msg, err := messageToMsg(do); err == nil {
			res = append(res, msg)
		}
	}

	return res
}

func (e *EventHandler) sendEventToMessage(evt *wsevent.MsgSendEvent) (*message.Message, error) {
	content, err := message.DecodeContent(evt.Msg.Content)
	if err != nil {
		return nil, err
	}
	return &message.Message{
		Header: message.Header{
			Id:        evt.Msg.Id,
			SendId:    evt.Id,
			SendAt:    evt.CreatedAt,
			UserId:    evt.UserId,
			ChannelId: evt.Msg.ChannelId,
			SeqId:     evt.SeqId,
		},
		Content: content,
	}, nil
}
