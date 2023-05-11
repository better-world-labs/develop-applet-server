package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type IMiniAppUserExtra interface {
	CompleteGuidance(userId int64) error
	GetByUserId(userId int64) (entity.MiniAppUserExtra, bool, error)
}
