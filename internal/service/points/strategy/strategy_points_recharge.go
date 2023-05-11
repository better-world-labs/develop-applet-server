package strategy

import (
	"errors"
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyPointsRecharge struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategyPointsRecharge() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypePointsRecharge); ok {
		return &strategyPointsRecharge{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyPointsRecharge) GetType() string {
	return s.definition.Type
}

func (s strategyPointsRecharge) GetDescription() string {
	return s.definition.Description
}

func (s strategyPointsRecharge) GetPoints(arg entity.IStrategyArg) (int, error) {
	if argRecharge, ok := arg.(entity.StrategyArgRecharge); ok {
		return argRecharge.Points, nil
	}

	return 0, errors.New("invalid IStrategyArg typ")
}

func (s strategyPointsRecharge) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	return true, nil
}
