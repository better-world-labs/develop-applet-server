package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewRechargeController() gone.Goner {
	return &recharge{}
}

type recharge struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter       `gone:"router-auth"`
	Recharge      service.IRecharge `gone:"*"`
}

func (con *recharge) Mount() gin.MountError {
	con.
		AuthRouter.
		GET("/points-goods", con.listGoods).
		GET("/points-orders/:orderId/state", con.getOrderState).
		GET("/points-orders", con.listOrders).
		POST("/points-goods/:goodsId/points-orders", con.recharge)

	return nil
}

func (con *recharge) listGoods(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	goods, err := con.Recharge.ListRechargeGoods(userId)
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{
		List: goods,
	}, nil
}

func (con *recharge) recharge(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	goodsId, err := utils.CtxPathParamInt64(ctx, "goodsId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	order, err := con.Recharge.CreateRechargeOrder(goodsId, userId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (con *recharge) listOrders(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	orders, err := con.Recharge.ListRechargeOrders(userId)
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{List: orders}, nil
}

func (con *recharge) getOrderState(ctx *gin.Context) (any, error) {
	orderId := ctx.Param("orderId")
	order, has, err := con.Recharge.GetRechargeOrderByOrderId(orderId)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return map[string]any{
		"state": order.State,
	}, nil
}
