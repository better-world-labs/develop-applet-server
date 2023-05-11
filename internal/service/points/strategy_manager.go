package points

import (
	"errors"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type strategyManager struct {
	gone.Goner
	strategyMap map[string]IStrategy

	strategies     []IStrategy            `gone:"*"`
	svc            service.IPoints        `gone:"*"`
	notify         service.INotifyMessage `gone:"*"`
	emitter.Sender `gone:"gone-emitter"`
}

//go:gone
func NewStrategySvc() gone.Goner {
	return &strategyManager{}
}

func (s *strategyManager) Start(cemetery gone.Cemetery) error {
	s.strategyMap = make(map[string]IStrategy)

	for _, strategy := range s.strategies {
		s.strategyMap[strategy.GetType()] = strategy
	}

	return nil
}

func (s strategyManager) Stop(cemetery gone.Cemetery) error {
	return nil
}

func (s strategyManager) ApplyPoints(userId int64, arg entity.IStrategyArg) (int, error) {
	if strategy, ok := s.strategyMap[arg.Type()]; ok {
		valid, err := strategy.LimitValidate(userId, arg)
		if err != nil {
			return 0, err
		}

		if !valid {
			return 0, s.Send(&entity.PointsStrategyLimitEvent{
				UserId: userId,
				Type:   arg.Type(),
			})
		}

		points, err := strategy.GetPoints(arg)
		if err != nil {
			return 0, err
		}

		return points, s.svc.Add(userId, int64(points), strategy.GetType(), strategy.GetDescription())
	}

	return 0, errors.New("strategy for type " + arg.Type() + " not found")
}

func (s strategyManager) GetStrategyPoints(arg entity.IStrategyArg) (int, error) {
	if strategy, ok := s.strategyMap[arg.Type()]; ok {
		return strategy.GetPoints(arg)
	}

	return 0, errors.New("strategy for type " + arg.Type() + " not found")
}
