package miniapp

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type pCategory struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewPCategory() gone.Goner {
	return &pCategory{}
}

func (p pCategory) list() ([]*entity.MiniAppCategory, error) {
	var arr []*entity.MiniAppCategory
	return arr, p.Find(&arr)
}

func (p pCategory) checkExists(id int64) (bool, error) {
	return p.Table(entity.MiniAppCategory{}).Where("id = ?", id).Exist()
}
