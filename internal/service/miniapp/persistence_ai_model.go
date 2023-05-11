package miniapp

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type pAIModel struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewPAIModel() gone.Goner {
	return &pAIModel{}
}

func (p pAIModel) list() ([]*entity.MiniAppAiModel, error) {
	var arr []*entity.MiniAppAiModel
	return arr, p.Find(&arr)
}

func (p pAIModel) listCategory() ([]*entity.MiniAppAiModelCategory, error) {
	var arr []*entity.MiniAppAiModelCategory
	return arr, p.Find(&arr)
}

func (p pAIModel) checkExistsByName(name string) (bool, error) {
	return p.Where("name = ?", name).Exist()
}
