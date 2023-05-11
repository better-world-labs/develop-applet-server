package system

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type iPersistence interface {
	listByGroupId(group int, sort bool) ([]*entity.Emoticon, error)
	updateRefStat() error
}
