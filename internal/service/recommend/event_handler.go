package recommend

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

	Svc service.IRecommendMiniApp `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleRecommendApp)
}

func (e *EventHandler) handleRecommendApp(evt *entity.MiniAppRecommendEvent) error {
	e.Logger.Infof("handleRecommendApp: userId=%d, appId=%s", evt.CreatedBy, evt.AppId)

	return e.Svc.DoRecommendApp(evt.MiniAppRecommend)
}
