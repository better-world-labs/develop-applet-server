package wechat

import "github.com/wechatpay-apiv3/wechatpay-go/services/payments"

type (
	Response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
)

const (
	TransactionStateSuccess    = "SUCCESS"    //支付成功
	TransactionStateRefund     = "REFUND"     //转入退款
	TransactionStateNotPay     = "NOTPAY"     //未支付
	TransactionStateClosed     = "CLOSED"     //关闭
	TransactionStateRevoked    = "REVOKED"    //已撤销（仅付款码支付会返回）
	TransactionStateUserPaying = "USERPAYING" //UserPaying
	TransactionStatePayError   = "PAYERROR"   //支出失败
)

func IsTransactionStateSuccess(transaction *payments.Transaction) bool {
	return *transaction.TradeState == TransactionStateSuccess
}

func IsTransactionStateNotPay(transaction *payments.Transaction) bool {
	return *transaction.TradeState == TransactionStateNotPay
}

func IsTransactionStateClosed(transaction *payments.Transaction) bool {
	return *transaction.TradeState == TransactionStateClosed
}
