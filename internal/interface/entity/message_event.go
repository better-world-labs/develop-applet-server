package entity

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

// MessageReferenced 回复消息时触发
type MessageReferenced struct {
	*wsevent.Msg

	ReferenceId int64 `json:"referenceId"`
}

// MessageLikeEvent 用户对消息进行点赞时 触发
type MessageLikeEvent struct {
	MessageLike
}

// ImageMessageSavedEvent 图片消息被保存时 触发
type ImageMessageSavedEvent struct {
	message.ImageContent

	Origin *wsevent.Msg `json:"origin"`
}
