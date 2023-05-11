package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"time"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IUserSettings interface {
	SetWorkingTime(userId int64, startWorkingTime, endWorkingTime time.Time) error

	SetBossKey(userId int64, shortcut string) error

	GetUserSettings(id int64) (*entity.UserSettings, error)
}
