package user

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

type EventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender emitter.Sender `gone:"gone-emitter"`
	Svc    service.IUser  `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleUserOnline)
	on(e.handleUserOffline)
}

func (e *EventHandler) handleUserOnline(evt *event.UserOnlineEvent) error {
	e.Logger.Infof("user online: userId=%d", evt.UserId)

	return e.Svc.UpdateOnline(evt.UserId, true, evt.CreatedAt)
}

func (e *EventHandler) handleUserOffline(evt *event.UserOfflineEvent) error {
	e.Logger.Infof("user offline: userId=%d", evt.UserId)

	return e.Svc.UpdateOnline(evt.UserId, false, evt.CreatedAt)
}
