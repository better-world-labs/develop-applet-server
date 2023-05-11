package resign

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

const TableName = "resign_template"

//go:gone
func NewResignPersistence() gone.Goner {
	return &persistence{}
}

type persistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

func (p *persistence) listResignTemplates() ([]*entity.ResignTemplate, error) {
	templates := make([]*entity.ResignTemplate, 0)
	err := p.Table(TableName).Select("id, title, content").Find(&templates)
	return templates, err
}
