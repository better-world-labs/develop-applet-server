package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IAITool interface {
	List() ([]*entity.AiTool, error)
	ListCategories() ([]*entity.AiToolCategory, error)
}
