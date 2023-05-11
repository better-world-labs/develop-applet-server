package statistic

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

	miniApp service.IStatisticMiniApp `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleAppCreate)
	on(e.handleAppRunDone)
	on(e.handleAppComment)
	on(e.handleAppView)
	on(e.handleAppLikeChanged)
	on(e.handleAppOutputLikeChanged)
}

func (e *EventHandler) handleAppCreate(evt *entity.AppCreatedEvent) error {
	e.Logger.Infof("handleAppCreate: appId=%s, duplicateFrom=%s", evt.AppId, evt.DuplicateFrom)

	if len(evt.DuplicateFrom) > 0 {
		return e.miniApp.IncrementAppUseTimes(evt.DuplicateFrom)
	}

	return nil
}

func (e *EventHandler) handleAppRunDone(evt *entity.AppRunDoneEvent) error {
	e.Logger.Infof("handleAppRunDone: appId=%s", evt.AppId)

	return e.miniApp.IncrementAppRuntimes(evt.AppId)
}

func (e *EventHandler) handleAppView(evt *entity.AppViewEvent) error {
	e.Logger.Infof("handleAppView: appId=%s", evt.AppId)

	return e.miniApp.IncrementAppViewTimes(evt.AppId)
}

func (e *EventHandler) handleAppComment(evt *entity.MiniAppCommentedEvent) error {
	e.Logger.Infof("handleAppComment: appId=%s", evt.AppId)

	return e.miniApp.IncrementAppCommentTimes(evt.AppId)
}

func (e *EventHandler) handleAppLikeChanged(evt *entity.MiniAppLikeChangedEvent) error {
	e.Logger.Infof("handleAppLikeChanged: appId=%s", evt.AppId)

	if evt.Like {
		return e.miniApp.IncrementAppLikeTimes(evt.AppId)
	}

	return e.miniApp.DecrementAppLikeTimes(evt.AppId)
}

func (e *EventHandler) handleAppOutputLikeChanged(evt *entity.MiniAppOutputLikeChangedEvent) error {
	e.Logger.Infof("handleAppOutputLikeChanged: outputId=%s", evt.OutputId)

	if evt.BecomeLike() {
		if err := e.miniApp.IncrementAppOutputLikeTimes(evt.OutputId); err != nil {
			return err
		}
	}
	if evt.BecomeHate() {
		if err := e.miniApp.IncrementAppOutputHateTimes(evt.OutputId); err != nil {
			return err
		}
	}
	if evt.FromLike() {
		if err := e.miniApp.DecrementAppOutputLikeTimes(evt.OutputId); err != nil {
			return err
		}
	}
	if evt.FromHate() {
		if err := e.miniApp.DecrementAppOutputHateTimes(evt.OutputId); err != nil {
			return err
		}
	}

	return nil
}
