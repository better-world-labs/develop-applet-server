package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"

type IRetainMessage interface {
	ListRetainMessages(userId int64) ([]*entity.RetainMessage, error)
	SendRetainMessage(message *entity.RetainMessage) error
}
