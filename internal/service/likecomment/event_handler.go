package likecomment

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type EventHandler struct {
	gone.Goner
	logrus.Logger `gone:"gone-logger"`

	Svc service.ILikeCommentMiniApp `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleLikeApp)
	on(e.handleLikeAppOutput)
}

func (e *EventHandler) handleLikeApp(evt *entity.MiniAppLikeEvent) error {
	e.Logger.Infof("handleLikeApp: userId=%d, appId=%s", evt.CreatedBy, evt.AppId)

	return e.Svc.DoLikeApp(evt.MiniAppLike)
}

func (e *EventHandler) handleLikeAppOutput(evt *entity.MiniAppOutputLikeEvent) error {
	e.Logger.Infof("handleLikeAppOutput: userId=%d, outputId=%s", evt.CreatedBy, evt.OutputId)

	return e.Svc.DoLikeAppOutput(evt.MiniAppOutputLike)
}
