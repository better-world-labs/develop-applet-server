package service

import (
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"net/http"
)

type IWxPay interface {
	CreateOrderForNativeAPI(request *entity.WxPayOrderRequest) (*string, error)
	CreateOrderForJSAPI(request *entity.WxPayOrderRequest) (*entity.WxPayOrderResponse, error)
	QueryOrderInfoByOrderNo(orderId string) (*payments.Transaction, bool, error)
	QueryOrderInfoByWxTransactionId(wxTransactionId string) (*payments.Transaction, error)
	CloseOrderByOrderNo(orderNo string) error

	//DecryPayNoticeInfo 使用该接口解密微信回调信息
	DecryPayNoticeInfo(request *http.Request) (*payments.Transaction, error)

	//TransferByOpenId 给用户转账
	TransferByOpenId(req *entity.WxTransferRequest) (*transferbatch.InitiateBatchTransferResponse, error)

	//GetTransferByBatchId 转账单查询
	GetTransferByBatchId(batchId string) (*transferbatch.TransferBatchEntity, error)
}
