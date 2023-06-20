package aitool

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type svc struct {
	gone.Goner

	p         iPersistence         `gone:"*"`
	categoryP iCategoryPersistence `gone:"*"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s svc) List() ([]*entity.AiTool, error) {
	return s.p.list()
}

func (s svc) ListCategories() ([]*entity.AiToolCategory, error) {
	return s.categoryP.listCategories()
}
