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

type MemberUpdateHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender  emitter.Sender   `gone:"gone-emitter"`
	Channel service.IChannel `gone:"*"`
	Notice  service.INotice  `gone:"*"`
}

//go:gone
func ChannelNewEventHandler() gone.Goner {
	return &MemberUpdateHandler{}
}

func (e *MemberUpdateHandler) Consume(on emitter.OnEvent) {
	on(e.handleChannelMemberAdded)
	on(e.handleChannelMemberRemoved)
}

func (e *MemberUpdateHandler) handleChannelMemberAdded(evt *entity.ChannelMemberAdded) error {
	e.Logger.Infof("handleChannelMemberAdded: channelId=%d, userId=%d", evt.ChannelId, evt.UserId)

	return e.Sender.Send(&wsevent.TriggerEvent{
		Scope:      wsevent.TriggerEventScopeByScene,
		Scene:      wsevent.ChannelScene,
		Type:       wsevent.TriggerTypeUpdateUserList,
		SceneParam: strconv.FormatInt(evt.ChannelId, 10),
		Params:     []interface{}{"channel add member", evt.UserId},
	})
}

func (e *MemberUpdateHandler) handleChannelMemberRemoved(evt *entity.ChannelMemberRemoved) error {
	e.Logger.Infof("handleChannelMemberRemoved: channelId=%d, userId=%d", evt.ChannelId, evt.UserId)

	err := e.Sender.Send(&wsevent.TriggerEvent{
		Scope:  wsevent.TriggerEventScopeByUser,
		Type:   TriggerRemovedFromChannel,
		UserId: evt.UserId,
		Params: []interface{}{evt.ChannelId},
	})
	if err != nil {
		return err
	}

	return e.Sender.Send(&wsevent.TriggerEvent{
		Scope:      wsevent.TriggerEventScopeByScene,
		Scene:      wsevent.ChannelScene,
		Type:       wsevent.TriggerTypeUpdateUserList,
		SceneParam: strconv.FormatInt(evt.ChannelId, 10),
		Params:     []interface{}{"channel remove member", evt.UserId},
	})
}

func (e *MemberUpdateHandler) handleUserInfoUpdated(evt *entity.UserInfoUpdated) error {
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
