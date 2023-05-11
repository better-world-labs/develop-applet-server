package strategy

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyBeInvited struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategyBeInvited() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeBeInvited); ok {
		return &strategyBeInvited{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyBeInvited) GetType() string {
	return s.definition.Type
}

func (s strategyBeInvited) GetDescription() string {
	return s.definition.Description
}

func (s strategyBeInvited) GetPoints(arg entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategyBeInvited) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	exists, err := s.points.CheckExistsTodayByType(userId, s.GetType())
	return !exists, err
}
