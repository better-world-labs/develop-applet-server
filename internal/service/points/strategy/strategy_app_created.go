package strategy

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	_interface "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyAppCreated struct {
	gone.Goner
	definition     entity.PointsDefinition
	points         service.IPoints `gone:"*"`
	emitter.Sender `gone:"gone-emitter"`
}

//go:gone
func NewStrategyAppCreated() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeAppCreated); ok {
		return &strategyAppCreated{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyAppCreated) GetType() string {
	return s.definition.Type
}

func (s strategyAppCreated) GetDescription() string {
	return s.definition.Description
}

func (s strategyAppCreated) GetPoints(arg entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategyAppCreated) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	return s.checkTodayIncomeOK(userId)
}

func (s strategyAppCreated) checkTodayIncomeOK(userId int64) (bool, error) {
	pointsForType, err := s.points.GetTodayIncomeByType(s.GetType(), userId)
	return pointsForType < _interface.LimitMaxAppCreatedEarnPointsEveryday, err
}
