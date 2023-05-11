package limit

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/google/uuid"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type EventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	marker *marker                `gone:"*"`
	notify service.INotifyMessage `gone:"*"`
	app    service.IMiniApp       `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handlePointsStrategyLimitEvent)
}

func (e *EventHandler) handlePointsStrategyLimitEvent(evt *entity.PointsStrategyLimitEvent) error {
	e.Logger.Infof("handlePointsStrategyLimitEvent: userId=%d, type=%s\n", evt.UserId, evt.Type)

	switch evt.Type {
	case entity.PointsTypeAppCreated:
		return e.handleAppCreatePointsLimit(evt.UserId)
	}

	return nil
}

func (e *EventHandler) handleAppCreatePointsLimit(userId int64) error {
	isSent, err := e.marker.IsAppCreateNotifySent(userId)
	if err != nil {
		return err
	}

	if !isSent {
		return e.sendNotifyAndMarkSent(userId)
	}

	return nil

}

func (e *EventHandler) sendNotifyAndMarkSent(userId int64) error {
	if err := e.marker.MarkAppCreatePointsLimitNotifySent(userId); err != nil {
		return err
	}

	return e.notify.SendNotify(entity.NotifyMessageCreateAppPointsLimited, userId, "创建积分上限", "哇呜，你已经创建超过了10个小程序了，精神可嘉！不过奖励的积分达到了上限哦", uuid.NewString())
}
