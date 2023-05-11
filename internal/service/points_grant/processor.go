package points_grant

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type processor struct {
	gone.Goner
	logrus.Logger `gone:"gone-logger"`

	points  service.IPointStrategy `gone:"*"`
	miniApp service.IMiniApp       `gone:"*"`
}

//go:gone
func NewProcessor() gone.Goner {
	return &processor{}
}

func (e processor) Consume(on emitter.OnEvent) {
	on(e.handleAppCreated)
}

func (e processor) handleAppCreated(evt *entity.AppCreatedEvent) error {
	e.Infof("handleAppCreated: appId=%s\n", evt.AppId)

	if evt.DuplicateFrom == "" {
		_, err := e.points.ApplyPoints(evt.CreatedBy, entity.StrategyArgAppCreated{})
		return err
	}

	duplicateFrom, has, err := e.miniApp.GetAppByUuid(evt.DuplicateFrom)
	if err != nil {
		return err
	}

	if !has {
		e.Errorf("duplicateApp appId=%s not found\n", evt.DuplicateFrom)
		return nil
	}

	_, err = e.points.ApplyPoints(duplicateFrom.CreatedBy, entity.StrategyArgAppDuplicated{})
	return err
}
