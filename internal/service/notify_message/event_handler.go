package notify_message

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type EventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	svc *svc `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleCreateNotifyMessage)
}

func (e *EventHandler) handleCreateNotifyMessage(evt *entity.CreateNotifyMessageEvent) error {
	e.Infof("handleCreateNotifyMessage: userId=%d, title=%s, content=%s", evt.UserId, evt.Title, evt.Content)

	has, err := e.svc.checkNotifyMessageExistsByOperationId(evt.OperationId)
	if err != nil {
		return err
	}

	if !has {
		return e.svc.sendNotifySync(evt.NotifyMessageInfo)
	}

	return nil
}
