package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IMessageStat interface {
	HotMessageTop(channelId int64, top int) ([]*entity.HotMessage, error)
}
