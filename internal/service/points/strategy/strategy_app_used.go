package strategy

import (
	"errors"
	"github.com/gone-io/gone"
	_interface "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyAppUsed struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints  `gone:"*"`
	miniApp    service.IMiniApp `gone:"*"`
}

//go:gone
func NewStrategyAppUsed() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeAppUsed); ok {
		return &strategyAppUsed{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyAppUsed) GetType() string {
	return s.definition.Type
}

func (s strategyAppUsed) GetDescription() string {
	return s.definition.Description
}

func (s strategyAppUsed) GetPoints(arg entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategyAppUsed) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	ok, err := s.checkTodayIncomeOK(userId)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return s.checkSameUserUsedOK(arg)
}

func (s strategyAppUsed) checkTodayIncomeOK(userId int64) (bool, error) {
	pointsForType, err := s.points.GetTodayIncomeByType(s.GetType(), userId)
	return pointsForType < _interface.LimitMaxAppUsedEarnPointsEveryDay, err
}

func (s strategyAppUsed) checkSameUserUsedOK(arg entity.IStrategyArg) (bool, error) {
	if strategy, ok := arg.(entity.StrategyArgAppUsed); ok {
		ran, err := s.miniApp.CheckAppRanByUser(strategy.App.Uuid, strategy.RunUserId)
		if err != nil {
			return false, err
		}

		return !ran, nil
	}

	return false, errors.New("invalid strategyArg type")
}
