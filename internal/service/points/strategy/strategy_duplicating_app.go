package strategy

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyDuplicatingApp struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategyDuplicatingApp() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeDuplicatingApp); ok {
		return &strategyDuplicatingApp{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyDuplicatingApp) GetType() string {
	return s.definition.Type
}

func (s strategyDuplicatingApp) GetDescription() string {
	return s.definition.Description
}

func (s strategyDuplicatingApp) GetPoints(arg entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategyDuplicatingApp) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	return true, nil
}
