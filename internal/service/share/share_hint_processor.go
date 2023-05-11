package share

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	_interface "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/service/push"
)

type HintProcessor struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	miniApp service.IMiniApp `gone:"*"`
	push    *push.PushSvc    `gone:"*"`
}

//go:gone
func NewShareHintProcessor() gone.Goner {
	return &HintProcessor{}
}

func (e *HintProcessor) Consume(on emitter.OnEvent) {
	on(e.handleAppCreated)
	on(e.handleAppRan)
}
func (e *HintProcessor) handleAppCreated(evt *entity.AppCreatedEvent) error {
	e.Logger.Infof("handleAppCreated: appId=%s\n", evt.AppId)

	count, err := e.miniApp.CountUserCreatedApps(evt.CreatedBy)
	if err != nil {
		return err
	}

	e.Logger.Infof("CountUserCreatedApps: %d\n", count)
	if count == _interface.ShareHintCreateAppThreshold {
		return e.pushAppCreatedShareHint(evt.CreatedBy, count, _interface.PointsStrategyCreateAppEarn)
	}

	return nil
}

func (e *HintProcessor) handleAppRan(evt *entity.AppOutputCreatedEvent) error {
	e.Logger.Infof("handleAppRan: outputId=%s, userId=%d\n", evt.OutputId, evt.UserId)

	count, err := e.miniApp.CountUserRanApps(evt.UserId)
	if err != nil {
		return err
	}

	if count == _interface.ShareHintUseAppThreshold {
		return e.pushAppRanShareHint(evt.UserId, count, _interface.PointsStrategyUsingAppMaxCost*count)
	}

	return nil
}

func (e *HintProcessor) pushAppCreatedShareHint(userId, createdApps, earnPoints int64) error {
	return e.push.PushMessage(userId, map[string]any{
		"type": "share-hint-create-app",
		"payload": map[string]any{
			"createdApps": createdApps,
			"earnPoints":  earnPoints,
		},
	})
}

func (e *HintProcessor) pushAppRanShareHint(userId, usedApps, costPoints int64) error {
	return e.push.PushMessage(userId, map[string]any{
		"type": "share-hint-use-app",
		"payload": map[string]any{
			"usedApps":   usedApps,
			"costPoints": costPoints,
		},
	})
}
