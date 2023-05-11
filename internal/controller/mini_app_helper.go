package controller

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

func getCategory(categories []*entity.MiniAppAiModelCategory, category int64) (*entity.MiniAppAiModelCategory, bool) {
	for _, c := range categories {
		if category == c.Id {
			return c, true
		}
	}

	return nil, false
}
