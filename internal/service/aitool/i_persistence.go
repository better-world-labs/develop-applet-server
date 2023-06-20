package aitool

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type iPersistence interface {
	list() ([]*entity.AiTool, error)
}

type iCategoryPersistence interface {
	listCategories() ([]*entity.AiToolCategory, error)
}
