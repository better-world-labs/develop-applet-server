package trigger

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
	"sync"
	"time"
)

type NoticeEventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	app    gone.Heaven     `gone:"gone-heaven"`
	Sender emitter.Sender  `gone:"gone-emitter"`
	Notice service.INotice `gone:"*"`

	timeoutHandlersByUser sync.Map
}

//go:gone
func NewNoticeEventHandler() gone.Goner {
	return &NoticeEventHandler{timeoutHandlersByUser: sync.Map{}}
}

func (e *NoticeEventHandler) Consume(on emitter.OnEvent) {
	utils.ListenBatch[*entity.NoticeMessageCreatedEvent](on, 100, 5*time.Second, e.app.GetHeavenStopSignal(), e.handleNoticeBatch)
	utils.ListenBatch[*entity.NoticeMessageRead](on, 100, 2*time.Second, e.app.GetHeavenStopSignal(), e.handleNoticeReadBatch)
}

func (e *NoticeEventHandler) handleNoticeBatch(events []*entity.NoticeMessageCreatedEvent) {
	e.Logger.Infof("handleNoticeBatch: size=%d", len(events))

	linq.From(events).SortT(func(i, j *entity.NoticeMessageCreatedEvent) bool {
		return i.CreatedAt.After(j.CreatedAt)
	}).GroupByT(func(i *entity.NoticeMessageCreatedEvent) interface{} {
		return i.UserId
	}, func(i interface{}) interface{} {
		return i
	}).ForEachT(func(group linq.Group) {
		i := group.Group[0].(*entity.NoticeMessageCreatedEvent)

		e.Logger.Infof("sendTrigger for notice id=%d", i.Id)
		err := e.Sender.Send(&wsevent.TriggerEvent{
			Scope:  wsevent.TriggerEventScopeByUser,
			Type:   wsevent.TriggerTypeNotice,
			UserId: i.UserId,
			Params: []interface{}{i.Id},
		})

		if err != nil {
			e.Logger.Error("sendTimerTrigger failed: %v", err)
		}
	})
}

func (e *NoticeEventHandler) handleNoticeReadBatch(events []*entity.NoticeMessageRead) {
	e.Logger.Infof("handleNoticeReadBatch: size=%d", len(events))

	var userIds []int64

	linq.From(events).SelectT(func(e *entity.NoticeMessageRead) int64 {
		return e.UserId
	}).Distinct().ToSlice(&userIds)

	for _, u := range userIds {
		err := e.Sender.Send(&wsevent.TriggerEvent{
			Scope:  wsevent.TriggerEventScopeByUser,
			Type:   wsevent.TriggerTypeNotice,
			UserId: u,
		})

		if err != nil {
			e.Errorf("send trigger failed: %v", err)
		}
	}
}

func (e *NoticeEventHandler) handleNotice(evt *entity.NoticeMessageCreatedEvent) error {
	e.Logger.Infof("handleNotice: type=%s userId=%d, businessId=%d", evt.Type, evt.UserId, evt.BusinessId)

	return e.SendTimedTrigger(&wsevent.TriggerEvent{
		Scope:  wsevent.TriggerEventScopeByUser,
		Type:   wsevent.TriggerTypeNotice,
		UserId: evt.UserId,
		Params: []interface{}{evt.Id},
	})
}

func (e *NoticeEventHandler) SendTimedTrigger(trigger *wsevent.TriggerEvent) error {
	h, ok := e.timeoutHandlersByUser.Load(trigger.UserId)
	if !ok {
		h = NewTimeoutHandler(time.Second * 2)
		e.timeoutHandlersByUser.Store(trigger.UserId, h)
	}

	h.(*TimeoutHandler).Handle(func() {
		err := e.Sender.Send(trigger)
		if err != nil {
			e.Error("trigger send failed")
		}

		e.timeoutHandlersByUser.Delete(trigger.UserId)
	})

	return nil
}
