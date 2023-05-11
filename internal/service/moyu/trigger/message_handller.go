package trigger

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
	"time"
)

type MessageEventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender  emitter.Sender         `gone:"gone-emitter"`
	Channel service.IChannel       `gone:"*"`
	Notice  service.INotice        `gone:"*"`
	Message service.IMessageRecord `gone:"*"`
	User    service.IUser          `gone:"*"`
	app     gone.Heaven            `gone:"gone-heaven"`
}

//go:gone
func MessageNewEventHandler() gone.Goner {
	return &MessageEventHandler{}
}

func (e *MessageEventHandler) Consume(on emitter.OnEvent) {
	on(e.HandleUpdateUnreadEvent)
	on(e.HandleMsgAckEvent)
	utils.ListenBatch[*wsevent.MsgSavedEvent](on, 10, 2*time.Second, e.app.GetHeavenStopSignal(), e.handleMsgSavedBatch)
}

func (e *MessageEventHandler) handleMsgSavedBatch(evts []*wsevent.MsgSavedEvent) {
	e.Logger.Infof("Trying to handle msg in batches. ")
	channelMsgMap := make(map[int64]int64, 0)
	channelUserMap := make(map[int64][]int64, 0)
	channelMaxMsgIdMap := make(map[int64]int64, 0)
	linq.From(evts).ForEachT(func(aEvt *wsevent.MsgSavedEvent) {
		cId := aEvt.Msg.ChannelId
		// if people in this channel, all msg is read
		channelMsgMap[cId] += 1
		channelUserMap[cId] = append(channelUserMap[cId], aEvt.Msg.UserId)
		if channelMaxMsgIdMap[cId] < aEvt.Msg.Id {
			channelMaxMsgIdMap[cId] = aEvt.Msg.Id
		}
	})

	for cId, incr := range channelMsgMap {
		err := e.Channel.IncreaseUnreadNum(cId, incr, channelMaxMsgIdMap[cId], channelUserMap[cId])
		if err != nil {
			continue
		}
	}

	_ = e.Sender.Send(&entity.UpdateUnreadEvent{
		UpdateUnreadRange: entity.All,
	})
}

func (e *MessageEventHandler) HandleMsgAckEvent(evt *wsevent.MsgAckEvent) error {
	e.Logger.Infof("handleMsgAckEvent: userId=%d, channelId=%d", evt.UserId, evt.ChannelId)
	//1. update user last read message id at current channel
	_, err := e.Channel.UpdateLastReadMessage(evt.UserId, evt.ChannelId)
	if err != nil {
		e.Logger.Errorf("update last read msg id failed at msg ack event handle, user id [%d], channel id [%d], msg id [%d], err: %s",
			evt.UserId, evt.ChannelId, evt.MsgId, err.Error())
		return nil
	}

	//2. send a UnreadCntTriggerEvent to this user (not all)
	//_ = e.Sender.Send(&entity.UpdateUnreadEvent{
	//	UpdateUnreadRange: entity.Partial,
	//	UserIds:           []int64{evt.CreatedBy},
	//})

	return nil
}

func (e *MessageEventHandler) HandleUpdateUnreadEvent(evt *entity.UpdateUnreadEvent) error {
	onlineIds := evt.UserIds

	if len(onlineIds) == 0 && evt.UpdateUnreadRange == entity.All {
		e.Logger.Infof("Send Unread Trigger Event to all online users. ")
		users, err := e.User.ListOnlineUsers()
		if err != nil {
			e.Logger.Errorf("List online users failed. err: %s", err.Error())
			return nil
		}

		linq.From(users).SelectT(func(user *entity.User) int64 {
			return user.Id
		}).ToSlice(&onlineIds)
	}

	e.Logger.Infof("found online users [%v]", onlineIds)

	for _, id := range onlineIds {
		unreadNums, err := e.Channel.GetAllUnReadCountByUser(id, 1)
		if err != nil {
			e.Logger.Errorf("get all channel unread num for user [%d] failed and skip it. err: %s ", err.Error())
			continue
		}

		var params []interface{}
		linq.From(unreadNums).SelectT(func(anUnread *domain.UserChannelUnreadNum) domain.UserChannelUnreadNum {
			return *anUnread
		}).ToSlice(&params)

		e.Logger.Infof("get user [%d] unread nums success, trying to send it. ", id)
		triggerEvent := &wsevent.TriggerEvent{
			Scope:  wsevent.TriggerEventScopeByUser,
			Scene:  wsevent.ClientScene,
			UserId: id,
			Type:   "updateUnreadMsgCnt",
			Params: params,
		}

		_ = e.Sender.Send(triggerEvent)
	}

	return nil
}
