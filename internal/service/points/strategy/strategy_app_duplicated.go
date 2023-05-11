package strategy

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyAppDuplicated struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategyAppDuplicated() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeAppDuplicated); ok {
		return &strategyAppDuplicated{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyAppDuplicated) GetType() string {
	return s.definition.Type
}

func (s strategyAppDuplicated) GetDescription() string {
	return s.definition.Description
}

func (s strategyAppDuplicated) GetPoints(arg entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategyAppDuplicated) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	exists, err := s.points.CheckExistsTodayByType(userId, s.GetType())
	return !exists, err
}
