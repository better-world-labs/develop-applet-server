package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

type INotifyMessage interface {
	SendNotify(_type entity.NotifyMessageType, userId int64, Title, content, operationId string) error

	PageNotifyMessages(userId int64, query page.StreamQuery, filter entity.NotifyMessageListFilter) (*page.StreamResult[*entity.NotifyMessage], error)

	CountUnread(userId int64) (int64, error)

	MarkRead(userId, id int64) error

	MarkReadAll(userId int64) error
}
