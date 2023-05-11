package domain

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type (
	MessageLike struct {
		entity.MessageLikeCount

		IsLike bool `json:"isLike"`
	}
)
