package domain

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type (
	Notice struct {
		entity.Notice

		User     entity.UserSimple `json:"user"`
		Business any               `json:"business"`
	}
)
