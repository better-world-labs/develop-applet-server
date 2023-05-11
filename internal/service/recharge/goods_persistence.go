package recharge

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type goodsPersistence struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewGoodPersistence() gone.Goner {
	return &goodsPersistence{}
}

func (g goodsPersistence) list() ([]*entity.PointsRechargeGoodsDefinition, error) {
	var slice []*entity.PointsRechargeGoodsDefinition
	return slice, g.Find(&slice)
}

func (g goodsPersistence) getById(id int64) (*entity.PointsRechargeGoodsDefinition, bool, error) {
	var res entity.PointsRechargeGoodsDefinition
	has, err := g.Where("id = ?", id).Get(&res)
	return &res, has, err
}
