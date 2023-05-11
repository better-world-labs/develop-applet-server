package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IResign interface {
	ListTemplates() (*domain.ResignList, error)
}
