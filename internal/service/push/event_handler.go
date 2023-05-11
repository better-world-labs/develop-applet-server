package push

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type EventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	svc    *PushSvc               `gone:"*"`
	notify service.INotifyMessage `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleNoticeMessageCreated)
	on(e.handleRetainMessageCreated)
	on(e.handleUserPointsChanged)
}

func (e *EventHandler) handleNoticeMessageCreated(evt *entity.NotifyMessageCreatedEvent) error {
	e.Logger.Infof("NotifyMessageCreatedEvent: userId=%d, id=%d", evt.UserId, evt.Id)

	unread, err := e.notify.CountUnread(evt.UserId)
	if err != nil {
		return err
	}

	if unread > 0 {
		return e.svc.PushMessage(evt.UserId, map[string]any{
			"type": "notify-message-changed",
			"payload": map[string]any{
				"unread": unread,
			},
		})
	}

	return nil
}

func (e *EventHandler) handleRetainMessageCreated(evt *entity.RetainMessageCreatedEvent) error {
	e.Logger.Infof("handleRetainMessageCreated: userId=%d, id=%d", evt.UserId, evt.Id)

	return e.svc.PushMessage(evt.UserId, map[string]any{
		"type": "retain-message-changed",
	})
}

func (e *EventHandler) handleUserPointsChanged(evt *entity.PointsChangeEvent) error {
	e.Logger.Infof("handleUserPointsChanged: userId=%d", evt.UserId)

	return e.svc.PushMessage(evt.UserId, map[string]any{
		"type": "user-points-changed",
	})
}
