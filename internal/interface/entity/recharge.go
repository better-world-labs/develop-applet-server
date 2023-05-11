package entity

import (
	"errors"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/helper"
	"time"
)

const (
	RechargeGoodsTagNewDeal = "new-deal"
)

const (
	RechargeOrderStateNotPay = 0
	RechargeOrderStatePayed  = 1
	RechargeOrderStateClosed = 2
)
const (
	DefaultOrderExpiresIn = time.Hour * 2
)

var (
	ErrorOrderAlreadyClosed = errors.New("order already closed")
	ErrorOrderAlreadyPayed  = errors.New("order already Payed")
)

type PointsRechargeGoods struct {
	Id          int64  `json:"id"`
	Price       CNY    `json:"price"`
	Tag         string `json:"tag"`
	Points      int64  `json:"points"`
	Description string `json:"description"`
}

type PointsRechargeGoodsDefinition struct {
	Id          int64  `json:"id"`
	Price       CNY    `json:"price"`
	Tag         string `json:"tag"`
	Points      int64  `json:"points"`
	Coupon      CNY    `json:"-"`
	Description string `json:"description"`
}

type PointsRechargeTransaction struct {
	PayedAt             *string `json:"payedAt"`
	PayTradeType        *string `json:"-"`
	PayTradeState       *string `json:"-"`
	PayBankType         *string `json:"-"`
	PayOpenid           *string `json:"-"`
	PayAmountTotal      *int64  `json:"-"`
	PayAmountPayerTotal *int64  `json:"-"`
	PayTransactionId    *string `json:"-"`
}

func CreatePointsRechargeTransaction(transaction *payments.Transaction) PointsRechargeTransaction {
	orderTransaction := PointsRechargeTransaction{
		PayTransactionId: transaction.TransactionId,
		PayBankType:      transaction.BankType,
		PayTradeState:    transaction.TradeState,
		PayedAt:          transaction.SuccessTime,
	}

	if transaction.Payer != nil {
		orderTransaction.PayOpenid = transaction.Payer.Openid
	}

	if transaction.Amount != nil {
		orderTransaction.PayAmountTotal = transaction.Amount.Total
		orderTransaction.PayAmountPayerTotal = transaction.Amount.PayerTotal
	}

	return orderTransaction
}

type PointsRechargeOrder struct {
	Id           int64               `json:"-"`
	OrderId      string              `json:"orderId"`
	Goods        PointsRechargeGoods `xorm:"json" json:"goods"`
	GoodsId      int64               `json:"goodsId"`
	Price        CNY                 `json:"price"`
	UserId       int64               `json:"userId"`
	CreatedAt    time.Time           `json:"createdAt"`
	State        int                 `json:"state"`
	CodeUrl      string              `xorm:"-" json:"codeUrl"`
	PayExpiresAt time.Time           `json:"payExpiresAt"`

	PointsRechargeTransaction `xorm:"extends"`
}

func (o PointsRechargeGoods) Buy(user User) (*PointsRechargeOrder, error) {
	return &PointsRechargeOrder{
		OrderId:      helper.GenerateSerialNumber(),
		Goods:        o,
		GoodsId:      o.Id,
		Price:        o.Price,
		UserId:       user.Id,
		CreatedAt:    time.Now(),
		State:        RechargeOrderStateNotPay,
		PayExpiresAt: time.Now().Add(DefaultOrderExpiresIn),
	}, nil
}

func (o PointsRechargeOrder) IsPayed() bool {
	return o.State == RechargeOrderStatePayed
}

func (o PointsRechargeOrder) IsClosed() bool {
	return o.State == RechargeOrderStateClosed
}

func (o *PointsRechargeOrder) MarkStatePayed(transaction PointsRechargeTransaction) error {
	if o.IsPayed() {
		return ErrorOrderAlreadyPayed
	}

	if o.IsClosed() {
		return ErrorOrderAlreadyClosed
	}

	o.PointsRechargeTransaction = transaction
	o.State = RechargeOrderStatePayed

	return nil
}

func (o *PointsRechargeOrder) MarkStateClosed() error {
	if o.IsPayed() {
		return ErrorOrderAlreadyPayed
	}

	if o.IsClosed() {
		return ErrorOrderAlreadyClosed
	}

	o.State = RechargeOrderStateClosed

	return nil
}
