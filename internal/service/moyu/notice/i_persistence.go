package notice

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

type iPersistence interface {
	pageByUserId(query page.Query, userId int64) (*page.Result[*entity.Notice], error)

	create(message *entity.Notice) error

	markRead(userId int64, messageId []int64) error

	countUnread(userId int64) (int64, error)

	listUnreadIMMessage(userId, ackMessageId, channelId int64) ([]int64, error)

	listUnreadNotice(noticeType entity.NoticeType, businessId int64) ([]*entity.Notice, error)

	getById(id int64) (*entity.Notice, bool, error)
}
