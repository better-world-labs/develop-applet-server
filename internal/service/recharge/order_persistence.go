package recharge

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type orderPersistence struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewOrderPersistence() gone.Goner {
	return &orderPersistence{}
}

func (p orderPersistence) getBySerialNumber(serialNumber string) (*entity.PointsRechargeOrder, bool, error) {
	var res entity.PointsRechargeOrder
	has, err := p.Where("order_id = ?", serialNumber).Get(&res)
	return &res, has, err
}

func (p orderPersistence) checkUserPayedOrderExists(userId int64) (bool, error) {
	return p.Table(entity.PointsRechargeOrder{}).Where("user_id = ? and state = ?", userId, entity.RechargeOrderStatePayed).Exist()
}

func (p orderPersistence) getById(id int64) (*entity.PointsRechargeOrder, bool, error) {
	var res entity.PointsRechargeOrder
	has, err := p.Where("id = ?", id).Get(&res)
	return &res, has, err
}

func (p orderPersistence) create(order *entity.PointsRechargeOrder) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := session.Insert(order)
		if err != nil {
			return err
		}

		return nil
	})
}

func (p orderPersistence) listByUserId(userId int64, state int) ([]*entity.PointsRechargeOrder, error) {
	var res []*entity.PointsRechargeOrder
	return res, p.Where("user_id = ? and state = ?", userId, state).MustCols("state").Desc("id").Find(&res)
}

func (p orderPersistence) update(order *entity.PointsRechargeOrder) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := session.ID(order.Id).Omit("order_id").Update(order)
		if err != nil {
			return err
		}

		return nil
	})
}
