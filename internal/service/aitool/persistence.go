package aitool

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type persistence struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewPersistence() gone.Goner {
	return &persistence{}
}

func (p persistence) list() ([]*entity.AiTool, error) {
	var res []*entity.AiTool
	return res, p.Desc("id").Find(&res)
}

func (p persistence) listCategories() ([]*entity.AiToolCategory, error) {
	var res []*entity.AiToolCategory
	return res, p.Asc("sort").Find(&res)
}
