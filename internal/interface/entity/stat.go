package entity

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
)

type HotMessage struct {
	MsgId      int64           `json:"msgId"`
	UserId     int64           `json:"userId"`
	Content    message.Content `json:"content"`
	ReplyCount int             `json:"replyCount"`
}
