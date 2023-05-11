package strategy

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategySignIn struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategySignIn() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeSignIn); ok {
		return &strategySignIn{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategySignIn) GetType() string {
	return s.definition.Type
}

func (s strategySignIn) GetDescription() string {
	return s.definition.Description
}

func (s strategySignIn) GetPoints(arg entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategySignIn) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	exists, err := s.points.CheckExistsTodayByType(userId, s.GetType())
	return !exists, err
}
