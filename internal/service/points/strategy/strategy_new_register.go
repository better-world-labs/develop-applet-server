package strategy

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyNewRegister struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategyNewRegister() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeNewRegister); ok {
		return &strategyNewRegister{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyNewRegister) GetType() string {
	return s.definition.Type
}

func (s strategyNewRegister) GetDescription() string {
	return s.definition.Description
}

func (s strategyNewRegister) GetPoints(entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategyNewRegister) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	exists, err := s.points.CheckExistsTodayByType(userId, s.GetType())
	return !exists, err
}
