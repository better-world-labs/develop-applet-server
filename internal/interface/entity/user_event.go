package entity

import wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"

type (

	// UserInfoUpdated 用户信息更新时触发
	UserInfoUpdated struct {
		UserId int64 `json:"userId"`
	}

	// UserMentioned 消息 @ 用户时触发
	UserMentioned struct {
		Id          int64        `json:"id"`
		Sender      int64        `json:"sender"`
		TargetUsers []int64      `json:"targetUsers"`
		ChannelId   int64        `json:"channelId"`
		Msg         *wsevent.Msg `json:"message"`
		Text        string       `json:"text"`
	}

	FirstLoginEvent struct {
		User
	}
)
