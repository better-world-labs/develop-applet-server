package retain_message

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/service/push"
)

type svc struct {
	gone.Goner
	Sender      emitter.Sender `gone:"gone-emitter"`
	xorm.Engine `gone:"gone-xorm"`

	p    iPersistence  `gone:"*"`
	push *push.PushSvc `gone:"*"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s svc) ListRetainMessages(userId int64) ([]*entity.RetainMessage, error) {
	offset, err := s.p.getReadOffset(userId)
	if err != nil {
		return nil, err
	}

	messages, err := s.p.listRetainMessages(userId, offset)
	if err != nil {
		return nil, err
	}

	if len(messages) > 0 {
		if err := s.p.markOffset(userId, messages[0].Id); err != nil {
			return nil, err
		}
	}

	return messages, err
}

func (s svc) SendRetainMessage(message *entity.RetainMessage) error {
	return s.Transaction(func(session xorm.Interface) error {
		if err := s.p.create(message); err != nil {
			return err
		}

		return s.Sender.Send(entity.RetainMessageCreatedEvent{
			RetainMessage: *message,
		})
	})
}
