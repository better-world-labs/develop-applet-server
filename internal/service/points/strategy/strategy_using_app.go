package strategy

import (
	"errors"
	"github.com/gone-io/gone"
	_interface "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyUsingApp struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategyUsingApp() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeUsingApp); ok {
		return &strategyUsingApp{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyUsingApp) GetType() string {
	return s.definition.Type
}

func (s strategyUsingApp) GetDescription() string {
	return s.definition.Description
}

func (s strategyUsingApp) GetPoints(arg entity.IStrategyArg) (int, error) {
	if argUsingApp, ok := arg.(entity.StrategyArgUsingApp); ok {
		if argUsingApp.Form == nil {
			return 0, errors.New("form is nil")
		}

		if len(argUsingApp.Form) <= 3 {
			return -_interface.PointsStrategyUsingAppMinCost, nil
		}

		return -_interface.PointsStrategyUsingAppMaxCost, nil
	}

	return 0, errors.New("invalid IStrategyArg typ")
}

func (s strategyUsingApp) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	return true, nil
}
