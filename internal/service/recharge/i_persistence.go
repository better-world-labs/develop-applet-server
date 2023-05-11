package recharge

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type (
	iGoodsPersistence interface {
		list() ([]*entity.PointsRechargeGoodsDefinition, error)
		getById(id int64) (*entity.PointsRechargeGoodsDefinition, bool, error)
	}

	iOrderPersistence interface {
		getBySerialNumber(serialNumber string) (*entity.PointsRechargeOrder, bool, error)
		checkUserPayedOrderExists(userId int64) (bool, error)
		listByUserId(userId int64, state int) ([]*entity.PointsRechargeOrder, error)
		getById(id int64) (*entity.PointsRechargeOrder, bool, error)
		create(order *entity.PointsRechargeOrder) error
		update(order *entity.PointsRechargeOrder) error
	}
)
