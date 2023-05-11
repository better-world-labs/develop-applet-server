package retain_message

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type iPersistence interface {
	listRetainMessages(userId, offsetId int64) ([]*entity.RetainMessage, error)
	getReadOffset(userId int64) (int64, error)
	create(message *entity.RetainMessage) error
	markOffset(userId, id int64) error
}
