package strategy

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyInvite struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategyInvite() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeInvite); ok {
		return &strategyInvite{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyInvite) GetType() string {
	return s.definition.Type
}

func (s strategyInvite) GetDescription() string {
	return s.definition.Description
}

func (s strategyInvite) GetPoints(arg entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategyInvite) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	return true, nil
}
