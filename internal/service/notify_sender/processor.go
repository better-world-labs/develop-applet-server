package notify_sender

import (
	"fmt"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type processor struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	svc service.INotifyMessage `gone:"*"`
}

//go:gone
func NewProcessor() gone.Goner {
	return &processor{}
}

func (e *processor) Consume(on emitter.OnEvent) {
	on(e.handlePointsChanged)
}

func (e *processor) handlePointsChanged(evt *entity.PointsChangeEvent) error {
	e.Infof("handlePointsChanged: userId=%d, points=%s, description=%s", evt.UserId, evt.Points, evt.Description)

	if mapping, ok := getNotifyMapping(evt.Type); ok {
		return e.svc.SendNotify(mapping.NotifyType, evt.UserId, evt.Description, fmt.Sprintf(mapping.NotifyContent, evt.AbsPoints()), evt.OperationId)
	}

	return nil
}
