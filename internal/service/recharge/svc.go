package recharge

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type svc struct {
	gone.Goner

	pay   service.IWxPay    `gone:"*"`
	user  service.IUser     `gone:"*"`
	goods iGoodsPersistence `gone:"*"`

	order *orderSvc `gone:"*"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s svc) GetRechargeOrderByOrderId(serialNumber string) (*entity.PointsRechargeOrder, bool, error) {
	return s.order.GetByOrderId(serialNumber)
}

func (s svc) ListRechargeOrders(userId int64) ([]*entity.PointsRechargeOrder, error) {
	return s.order.ListPayedOrders(userId)
}

func (s svc) ListRechargeGoods(userId int64) ([]*entity.PointsRechargeGoods, error) {
	list, err := s.goods.list()
	if err != nil {
		return nil, err
	}

	isFirstRecharge, err := s.order.IsFirstRecharge(userId)
	if err != nil {
		return nil, err
	}

	return s.processRechargeDefinitionBatch(isFirstRecharge, list)
}

func (s svc) GetRechargeGoods(userId, id int64) (*entity.PointsRechargeGoods, bool, error) {
	goods, has, err := s.goods.getById(id)
	if err != nil {
		return nil, has, err
	}

	isFirstRecharge, err := s.order.IsFirstRecharge(userId)
	if err != nil {
		return nil, has, err
	}

	g, err := s.processRechargeDefinition(isFirstRecharge, goods)
	if err != nil {
		return nil, has, err
	}

	return g, has, nil
}

func (s svc) processRechargeDefinitionBatch(firstRecharge bool, goods []*entity.PointsRechargeGoodsDefinition) (order []*entity.PointsRechargeGoods, err error) {
	var g []*entity.PointsRechargeGoods

	for _, good := range goods {
		definition, err := s.processRechargeDefinition(firstRecharge, good)
		if err != nil {
			return nil, err
		}

		g = append(g, definition)
	}

	return g, nil
}

func (s svc) processRechargeDefinition(firstRecharge bool, goods *entity.PointsRechargeGoodsDefinition) (*entity.PointsRechargeGoods, error) {
	g := entity.PointsRechargeGoods{
		Id:          goods.Id,
		Price:       goods.Price,
		Points:      goods.Points,
		Description: goods.Description,
	}

	if goods.Tag == entity.RechargeGoodsTagNewDeal {
		payCNY, err := goods.Price.Sub(goods.Coupon)
		if err != nil {
			return nil, err
		}

		if firstRecharge {
			g.Tag = goods.Tag
			g.Price = payCNY
		}
	}

	return &g, nil
}

func (s svc) CreateRechargeOrder(goodsId, userId int64) (*entity.PointsRechargeOrder, error) {
	goods, has, err := s.GetRechargeGoods(userId, goodsId)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, gin.NewParameterError("goods not found")
	}

	user, err := s.user.GetUserById(userId)
	if user == nil {
		return nil, gin.NewParameterError("user not found")
	}

	_ = s.order.CloseUsersOrders(userId)
	return s.order.CreateOrder(goods, user)
}

func (s svc) IsFirstRecharge(userId int64) (bool, error) {
	return s.order.IsFirstRecharge(userId)
}
