package recharge

import (
	"errors"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/xorm"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/service/wechat"
)

type orderSvc struct {
	gone.Goner
	xorm.Engine   `gone:"gone-xorm"`
	logrus.Logger `gone:"gone-logger"`

	order     iOrderPersistence `gone:"*"`
	pay       service.IWxPay    `gone:"*"`
	notifyUrl string            `gone:"config,wechat.pay.notify-url"`

	recharge service.OrderListener `gone:"*"`
}

//go:gone
func NewOrderSvc() gone.Goner {
	return &orderSvc{}
}

func (o orderSvc) GetByOrderId(serialNumber string) (*entity.PointsRechargeOrder, bool, error) {
	return o.order.getBySerialNumber(serialNumber)
}

func (o orderSvc) ListPayedOrders(userId int64) ([]*entity.PointsRechargeOrder, error) {
	return o.order.listByUserId(userId, entity.RechargeOrderStatePayed)
}

func (o orderSvc) CreateOrder(goods *entity.PointsRechargeGoods, user *entity.User) (order *entity.PointsRechargeOrder, err error) {
	err = o.Transaction(func(session xorm.Interface) error {
		order, err = o.localCreateOrder(goods, user)
		if err != nil {
			return err
		}

		return o.remoteCreateOrder(goods, order, user)
	})

	return
}

func (o orderSvc) IsFirstRecharge(userId int64) (bool, error) {
	exists, err := o.order.checkUserPayedOrderExists(userId)
	return !exists, err
}

func (o orderSvc) CloseUsersOrders(userId int64) error {
	notpayed, err := o.order.listByUserId(userId, entity.RechargeOrderStateNotPay)
	if err != nil {
		return err
	}

	for _, order := range notpayed {
		if err := o.closeOrder(order.OrderId); err != nil {
			o.Errorf("closeOrder error: %v\n", err)
		}
	}

	return nil
}

func (o orderSvc) closeOrder(orderId string) error {
	order, has, err := o.order.getBySerialNumber(orderId)
	if err != nil {
		return err
	}

	if !has {
		return errors.New("order not found")
	}

	if order.IsPayed() {
		return nil
	}

	return o.Transaction(func(session xorm.Interface) error {
		if err := o.markOrderStateClosed(order); err != nil {
			return err
		}

		return o.pay.CloseOrderByOrderNo(order.OrderId)
	})
}

func (o orderSvc) remoteCreateOrder(goods *entity.PointsRechargeGoods, order *entity.PointsRechargeOrder, user *entity.User) error {
	codeUrl, err := o.pay.CreateOrderForNativeAPI(&entity.WxPayOrderRequest{
		OrderNo:     order.OrderId,
		Description: goods.Description,
		TimeExpire:  order.PayExpiresAt,
		Total:       goods.Price.Fen(),
		Openid:      user.WxOpenId,
		NotifyUrl:   o.notifyUrl,
	})
	if err != nil {
		return err
	}

	order.CodeUrl = *codeUrl
	return nil
}

func (o orderSvc) localCreateOrder(goods *entity.PointsRechargeGoods, user *entity.User) (*entity.PointsRechargeOrder, error) {
	order, err := goods.Buy(*user)
	if err != nil {
		return nil, err
	}

	return order, o.order.create(order)
}

func (o orderSvc) handleTransactionSuccess(transaction *payments.Transaction) error {
	order, has, err := o.GetByOrderId(*transaction.OutTradeNo)
	if err != nil {
		return err
	}

	if !has {
		o.Warnf("handleTransactionSuccess: order not found")
		return nil
	}

	return o.Transaction(func(session xorm.Interface) error {
		err = o.markOrderStatePayed(order, entity.CreatePointsRechargeTransaction(transaction))
		if err != nil {
			return err
		}

		return o.recharge.OnOrderPayed(order)
	})
}

func (o orderSvc) markOrderStatePayed(order *entity.PointsRechargeOrder, transaction entity.PointsRechargeTransaction) error {
	err := order.MarkStatePayed(transaction)
	if err == entity.ErrorOrderAlreadyPayed {
		return nil
	}

	if err != nil {
		return err
	}

	return o.order.update(order)
}

func (o orderSvc) markOrderStateClosed(order *entity.PointsRechargeOrder) error {
	err := order.MarkStateClosed()
	if err == entity.ErrorOrderAlreadyClosed {
		return nil
	}

	if err != nil {
		return err
	}

	return o.order.update(order)
}

func (o orderSvc) handlePayNotify(transaction *payments.Transaction) error {
	o.Infof("handlePayNotify: transactionId=%s, serialNumber=%s, state=%s\n", transaction.TransactionId, transaction.OutTradeNo, transaction.TradeState)
	if wechat.IsTransactionStateSuccess(transaction) {
		return o.handleTransactionSuccess(transaction)
	}

	return nil
}
