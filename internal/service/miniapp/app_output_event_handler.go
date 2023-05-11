package miniapp

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"time"
)

type AppOutputEventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender emitter.Sender   `gone:"gone-emitter"`
	appSvc service.IMiniApp `gone:"*"`
}

//go:gone
func NewAppOutputEventHandler() gone.Goner {
	return &AppOutputEventHandler{}
}

func (e *AppOutputEventHandler) Consume(on emitter.OnEvent) {
	on(e.handleAppRunDone)
}

func (e *AppOutputEventHandler) handleAppRunDone(evt *entity.AppRunDoneEvent) error {
	e.Logger.Infof("handleAppRunDone: userId=%d, appId=%s", evt.User, evt.AppId)

	has, err := e.appSvc.CheckOutputExists(evt.Output.OutputId)
	if err != nil {
		return err
	}

	if !has {
		if err := e.appSvc.CreateOutput(&entity.MiniAppOutput{
			InputArgs:         evt.Param.Values,
			AppId:             evt.AppId,
			Open:              evt.Param.Open,
			CreatedAt:         time.Now(),
			CreatedBy:         evt.User,
			MiniAppOutputCore: evt.Output,
		}); err != nil {
			return err
		}

		return e.Sender.Send(entity.AppOutputCreatedEvent{
			OutputId: evt.Output.OutputId,
			UserId:   evt.User,
			AppId:    evt.AppId,
		})
	}

	return nil
}
