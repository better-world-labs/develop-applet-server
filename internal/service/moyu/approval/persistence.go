package approval

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

const TableName = "approval"

type persistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewPersistence() gone.Goner {
	return &persistence{}
}

func (p *persistence) listByIds(ids []int64) ([]*entity.Approval, error) {
	var res []*entity.Approval
	return res, p.Table(TableName).In("id", ids).Find(&res)
}

func (p *persistence) getById(id int64) (*entity.Approval, bool, error) {
	var a entity.Approval
	exists, err := p.Table(TableName).Where("id = ?", id).Get(&a)
	return &a, exists, err
}

func (p *persistence) create(approval *entity.Approval) error {
	_, err := p.Table(TableName).Insert(approval)
	return err
}

func (p *persistence) update(approval *entity.Approval) error {
	_, err := p.Table(TableName).ID(approval.Id).Update(approval)
	return err
}

func (p *persistence) delete(id int64) error {
	_, err := p.Table(TableName).Where("id = ?", id).Delete()
	return err
}
