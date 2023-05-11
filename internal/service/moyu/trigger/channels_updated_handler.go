package trigger

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

type ChannelsUpdateHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender emitter.Sender `gone:"gone-emitter"`
	User   service.IUser  `gone:"*"`
}

//go:gone
func NewChannelsUpdateHandler() gone.Goner {
	return &ChannelsUpdateHandler{}
}

func (e *ChannelsUpdateHandler) Consume(on emitter.OnEvent) {
	on(e.handleChannelsUpdated)
	on(e.handleChannelsDeleted)
	on(e.handleChannelGroupsUpdated)
}

func (e *ChannelsUpdateHandler) handleChannelGroupsUpdated(*entity.ChannelGroupsUpdated) error {
	e.Infof("handleChannelGroupsUpdated:")

	return e.sendTriggers(wsevent.TriggerTypeUpdateChannelList, nil)
}

func (e *ChannelsUpdateHandler) handleChannelsDeleted(evt *entity.ChannelDeleted) error {
	e.Infof("handleChannelsDeleted:")

	err := e.sendTriggers(wsevent.TriggerTypeUpdateChannelList, nil)
	if err != nil {
		return err
	}

	return e.sendTriggers(TriggerEventChannelDeleted, []any{evt.Id})
}

func (e *ChannelsUpdateHandler) handleChannelsUpdated(*entity.ChannelsUpdated) error {
	e.Infof("handleChannelsUpdated:")

	return e.sendTriggers(wsevent.TriggerTypeUpdateChannelList, nil)
}

func (e *ChannelsUpdateHandler) sendTriggers(triggerType string, params []any) error {
	users, err := e.User.ListOnlineUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		err := e.Sender.Send(&wsevent.TriggerEvent{
			Scope:  wsevent.TriggerEventScopeByUser,
			Type:   triggerType,
			UserId: u.Id,
			Params: params,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
