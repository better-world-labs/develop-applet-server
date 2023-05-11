package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

// INotice
// mockgen1.6 不支持泛型，在mockgen1.7将支持，特性还在rc阶段，暂时注释起来
// //go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type INotice interface {
	Create(userId int64, messageType entity.NoticeType, businessId int64, read bool) error

	Get(id int64) (*entity.Notice, bool, error)

	Page(userId int64, query page.Query) (*page.Result[*domain.Notice], error)

	MarkRead(userId int64, messageId []int64) error

	CountUnread(userId int64) (int64, error)

	ListUnreadIMMessages(userId, ackMessageId, channelId int64) ([]int64, error)

	ListUnreadNotice(noticeType entity.NoticeType, businessId int64) ([]*entity.Notice, error)
}
