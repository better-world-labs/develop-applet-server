package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type IPointStrategy interface {
	ApplyPoints(userId int64, arg entity.IStrategyArg) (int, error)
	GetStrategyPoints(arg entity.IStrategyArg) (int, error)
}
