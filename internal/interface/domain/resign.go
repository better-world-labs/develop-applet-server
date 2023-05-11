package domain

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type (
	ResignList struct {
		ResignList []*entity.ResignTemplate `json:"list"`
	}
)
