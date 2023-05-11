package domain

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type (
	ChannelIdRes struct {
		ChannelId int64 `json:"channelId"`
	}

	GroupIdRes struct {
		GroupId int64 `json:"groupId"`
	}

	ChannelMemberRes struct {
		Total          int64            `json:"total"`
		Online         int64            `json:"online"`
		ChannelMembers []*ChannelMember `json:"list"`
	}

	ChannelMember struct {
		Id       int64             `json:"id"`
		Nickname string            `json:"nickname"`
		Avatar   string            `json:"avatar"`
		Online   bool              `json:"online"`
		Role     entity.PlanetRole `json:"role"`
	}

	ChannelMemberLastReadMessageId struct {
		LastReadMessageId int64 `json:"lastReadMessageId"`
	}

	UserChannelUnreadNum struct {
		ChannelId int64 `json:"channelId"`
		UnreadNum int64 `json:"unreadNum"`
	}
)
