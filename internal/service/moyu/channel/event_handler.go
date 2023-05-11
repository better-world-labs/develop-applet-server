package channel

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

type EventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Sender emitter.Sender   `gone:"gone-emitter"`
	Planet service.IPlanet  `gone:"*"`
	Svc    service.IChannel `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleUserEnter)
	on(e.handleUserLeave)
}

func (e *EventHandler) handleUserEnter(evt *event.UserEnterEvent) error {
	e.Logger.Infof("handleUserEnter: userId=%d, channelId=%d", evt.UserId, evt.ChannelId)

	// 1. 如果用户不在该频道下，添加该用户到频道
	channel, exists, err := e.Svc.GetChannelByChannelId(evt.ChannelId)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	role, err := e.Planet.GetPlanetRoles(channel.PlanetId, evt.UserId)
	if err != nil {
		return err
	}

	if channel.Type == entity.ChannelTypeNormal || role > entity.PlanetRoleMember {
		_ = e.Svc.AddChannelMember(evt.ChannelId, evt.UserId)
	}

	// 2. 获取频道下该用户上次阅读的消息id,用于redis缓存获取失败时重新从计算
	lastRead, err := e.Svc.GetLastReadMsgId(evt.UserId, evt.ChannelId)
	if err != nil || lastRead == nil {
		// do nothing
		return nil
	}
	e.Logger.Infof("user:[%d] at channel:[%d] last read messageId: id[%d], ", evt.UserId, evt.ChannelId, lastRead.LastReadMessageId)

	total, err := e.Svc.GetUnreadMsgNum(lastRead.LastReadMessageId, evt.UserId, evt.ChannelId)
	if err != nil {
		// do nothing
		return nil
	}

	// 3. 避免用户多屏操作时未切换频道，最新未读消息未更新
	_, _ = e.Svc.UpdateLastReadMessage(evt.UserId, evt.ChannelId)

	// 4. 频道下有未读消息时，触发消息堆叠时间，提供：最后已读消息id & 未读消息数
	if total > 0 {
		var params []interface{}
		params = append(params, entity.MsgRecordSummary{
			UnreadCount: total,
			LastReadId:  lastRead.LastReadMessageId,
		})

		e.Logger.Infof(" user:[%d] at channel:[%d] has %d messages not read, trying to send message event.", evt.UserId, evt.ChannelId, total)
		err = e.Sender.Send(&event.TriggerEvent{
			Scope:  event.TriggerEventScopeBySid,
			Scene:  event.ClientScene,
			UserId: evt.UserId,
			Sid:    evt.Sid,
			Type:   event.TriggerTypeUnreadMsg,
			Params: params,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EventHandler) handleUserLeave(evt *event.UserLeaveEvent) error {
	e.Logger.Infof("handleUserLeave: userId=%d, channelId=%d", evt.UserId, evt.ChannelId)
	// 1. update user's last read message record
	_, _ = e.Svc.UpdateLastReadMessage(evt.UserId, evt.ChannelId)

	// 2. use MQ event
	//e.Logger.Infof("Trying to send unread event when user [%d] leave channel [%d]. ", evt.CreatedBy, evt.ChannelId)
	//_ = e.Sender.Send(&entity.UpdateUnreadEvent{
	//	UpdateUnreadRange: entity.All,
	//})
	return nil
}
