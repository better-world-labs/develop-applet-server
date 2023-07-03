package strategy

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyGptConversation struct {
	gone.Goner
	definition entity.PointsDefinition
	points     service.IPoints `gone:"*"`
}

//go:gone
func NewStrategyGptConversation() gone.Goner {
	if def, ok := getDefinition(entity.PointsTypeGptConversation); ok {
		return &strategyGptConversation{
			definition: def,
		}
	}

	panic("definition not found")
}

func (s strategyGptConversation) GetType() string {
	return s.definition.Type
}

func (s strategyGptConversation) GetDescription() string {
	return s.definition.Description
}

func (s strategyGptConversation) GetPoints(arg entity.IStrategyArg) (int, error) {
	return s.definition.Points, nil
}

func (s strategyGptConversation) LimitValidate(userId int64, arg entity.IStrategyArg) (bool, error) {
	return true, nil
}
