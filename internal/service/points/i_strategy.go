package points

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type IStrategy interface {
	GetType() string
	GetDescription() string
	GetPoints(arg entity.IStrategyArg) (int, error)
	LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error)
}
