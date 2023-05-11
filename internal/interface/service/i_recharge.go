package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type OrderListener interface {
	OnOrderPayed(order *entity.PointsRechargeOrder) error
	OnOrderClosed(order *entity.PointsRechargeOrder) error
}

type IRecharge interface {
	CreateRechargeOrder(goodsId, userId int64) (*entity.PointsRechargeOrder, error)
	GetRechargeOrderByOrderId(serialNumber string) (*entity.PointsRechargeOrder, bool, error)
	ListRechargeOrders(userId int64) ([]*entity.PointsRechargeOrder, error)
	IsFirstRecharge(userId int64) (bool, error)
	ListRechargeGoods(userId int64) ([]*entity.PointsRechargeGoods, error)
}
