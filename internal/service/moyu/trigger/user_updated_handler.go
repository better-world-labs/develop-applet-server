package trigger

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
	"strconv"
)

type UserEventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender  emitter.Sender   `gone:"gone-emitter"`
	Channel service.IChannel `gone:"*"`
	User    service.IUser    `gone:"*"`
	Notice  service.INotice  `gone:"*"`
}

//go:gone
func NewUserEventHandler() gone.Goner {
	return &UserEventHandler{}
}

func (e *UserEventHandler) Consume(on emitter.OnEvent) {
	on(e.handleUserInfoUpdated)
	on(e.handleUserOnline)
	on(e.handleUserOffline)
}

func (e *UserEventHandler) handleUserInfoUpdated(evt *entity.UserInfoUpdated) error {
	e.Logger.Infof("handleUserInfoUpdated: userId=%d", evt.UserId)

	joinedChannelIds, err := e.Channel.ListJoinedChannelIds(evt.UserId)
	if err != nil {
		return err
	}

	for _, channelId := range joinedChannelIds {
		err := e.Sender.Send(&wsevent.TriggerEvent{
			Scope:      wsevent.TriggerEventScopeByScene,
			Scene:      wsevent.ChannelScene,
			Type:       wsevent.TriggerTypeUpdateUserList,
			SceneParam: strconv.FormatInt(channelId, 10),
			Params:     []interface{}{"user update info", evt.UserId},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *UserEventHandler) handleUserOnline(evt *wsevent.UserOnlineEvent) error {
	e.Logger.Infof(" handleUserOnline: userId=%d", evt.UserId)

	joinedChannelIds, err := e.Channel.ListJoinedChannelIds(evt.UserId)
	if err != nil {
		return err
	}

	for _, channelId := range joinedChannelIds {
		err := e.Sender.Send(&wsevent.TriggerEvent{
			Scope:      wsevent.TriggerEventScopeByScene,
			Scene:      wsevent.ChannelScene,
			Type:       wsevent.TriggerTypeUpdateUserList,
			SceneParam: strconv.FormatInt(channelId, 10),
			Params:     []interface{}{"user online", evt.UserId},
		})
		if err != nil {
			return err
		}
	}

	// update user's connect time
	err = e.User.UpdateConnectTime(evt.UserId)
	if err != nil {
		e.Logger.Errorf("update user [%d] connect time failed, err: %s", evt.UserId, err.Error())
	}

	return nil
}

func (e *UserEventHandler) handleUserOffline(evt *wsevent.UserOfflineEvent) error {
	e.Logger.Infof(" handleUserOffline: userId=%d", evt.UserId)

	joinedChannelIds, err := e.Channel.ListJoinedChannelIds(evt.UserId)
	if err != nil {
		return err
	}

	for _, channelId := range joinedChannelIds {
		err := e.Sender.Send(&wsevent.TriggerEvent{
			Scope:      wsevent.TriggerEventScopeByScene,
			Scene:      wsevent.ChannelScene,
			Type:       wsevent.TriggerTypeUpdateUserList,
			SceneParam: strconv.FormatInt(channelId, 10),
			Params:     []interface{}{"user offline", evt.UserId},
		})
		if err != nil {
			return err
		}
	}

	// update user's total access duration by connect time
	err = e.User.AccumulateTotalAccessDuration(evt.UserId)
	if err != nil {
		e.Logger.Errorf("accumulate user [%d] access duration time failed, err: %s", evt.UserId, err.Error())
	}

	return nil
}
