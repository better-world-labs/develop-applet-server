package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/core/message"
)

type IContextManager interface {
	WriteMessage(record *message.Message) error
	GetContext(channelId int64) ([]*message.Message, error)
}
