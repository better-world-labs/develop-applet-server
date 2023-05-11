package recharge

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type rechargeListener struct {
	gone.Goner
	points service.IPointStrategy `gone:"*"`
}

//go:gone
func NewListener() gone.Goner {
	return &rechargeListener{}
}

func (s rechargeListener) OnOrderPayed(order *entity.PointsRechargeOrder) error {
	_, err := s.points.ApplyPoints(order.UserId, entity.StrategyArgRecharge{
		Points: int(order.Goods.Points),
	})
	return err
}

func (s rechargeListener) OnOrderClosed(order *entity.PointsRechargeOrder) error {
	return nil
}
