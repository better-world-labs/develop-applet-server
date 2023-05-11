package wechat

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"net/http"
	"time"
)

//go:gone
func NewPaySvc() gone.Goner {
	return &paySvc{}
}

type paySvc struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	mchID      string `gone:"config,wechat.pay.mchID"`
	certSerial string `gone:"config,wechat.pay.mch.cert-serial"`
	keyPath    string `gone:"config,wechat.pay.mch.key-path"`
	certPath   string `gone:"config,wechat.pay.mch.cert-path"`
	v3Secret   string `gone:"config,wechat.pay.secret.v3"`
	minaAppId  string `gone:"config,wechat.mina.appId"`

	client  *core.Client
	disable bool `gone:"config,wechat.pay.disable"`
	handler *notify.Handler
}

func (s *paySvc) AfterRevive() gone.AfterReviveError {
	if !s.disable {
		s.initClient()
		s.initHandler()
	}
	return nil
}

func (s *paySvc) initClient() {
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(s.keyPath)
	if err != nil {
		s.Fatal("load merchant private key error")
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(s.mchID, s.certSerial, mchPrivateKey, s.v3Secret),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		s.Fatalf("new wechat pay client err:%s", err)
	}

	s.client = client
}

func (s *paySvc) initHandler() {
	wechatPayCert, err := utils.LoadCertificateWithPath(s.certPath)
	if err != nil {
		s.Fatalf("utils.LoadCertificateWithPath:%v", err)
	}

	certificateVisitor := core.NewCertificateMapWithList([]*x509.Certificate{wechatPayCert})

	handler, err := notify.NewRSANotifyHandler(s.v3Secret, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))

	if err != nil {
		s.Fatalf("notify.NewRSANotifyHandler:%v", err)
	}
	s.handler = handler
}

func (s *paySvc) CreateOrderForNativeAPI(request *entity.WxPayOrderRequest) (*string, error) {
	svc := native.NativeApiService{Client: s.client}
	resp, result, err := svc.Prepay(context.Background(), native.PrepayRequest{
		Appid:       core.String(s.minaAppId),
		Mchid:       core.String(s.mchID),
		OutTradeNo:  core.String(request.OrderNo),
		Description: core.String(request.Description),
		TimeExpire:  core.Time(request.TimeExpire),
		Attach:      core.String(request.Attach),
		NotifyUrl:   core.String(request.NotifyUrl),
		Amount: &native.Amount{
			Total: core.Int64(request.Total),
		},
	})

	if err != nil {
		s.Errorf("wechat Prepay err:%v", err)
		return nil, gin.NewInnerError(err.Error(), 500)
	}

	if resp.CodeUrl == nil {
		s.Infof("wechat Prepay result: %v", result)
		return nil, gin.NewInnerError("PrepayId is empty", 500)
	}

	return resp.CodeUrl, nil
}

// CreateOrderForJSAPI 文档参考 https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_1.shtml
func (s *paySvc) CreateOrderForJSAPI(request *entity.WxPayOrderRequest) (*entity.WxPayOrderResponse, error) {
	s.Infof("CreateOrderForJSAPI: orderNo=%s, openid=%s", request.OrderNo, request.Openid)
	svc := jsapi.JsapiApiService{Client: s.client}
	ctx := context.Background()

	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, result, err := svc.Prepay(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(s.minaAppId),
			Mchid:       core.String(s.mchID),
			OutTradeNo:  core.String(request.OrderNo),
			Description: core.String(request.Description),
			TimeExpire:  core.Time(request.TimeExpire),
			Attach:      core.String(request.Attach),
			NotifyUrl:   core.String(request.NotifyUrl),
			Amount: &jsapi.Amount{
				Total: core.Int64(request.Total),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(request.Openid),
			},
		},
	)
	if err != nil {
		s.Errorf("wechat Prepay err:%v", err)
		return nil, gin.NewInnerError(err.Error(), 500)
	}

	if resp.PrepayId == nil {
		s.Infof("wechat Prepay result: %v", result)
		return nil, gin.NewInnerError("PrepayId is empty", 500)
	}

	response := entity.WxPayOrderResponse{
		PrepayId: *resp.PrepayId,
	}
	nonce, err := utils.GenerateNonce()
	if err != nil {
		s.Errorf("utils.GenerateNonce() err:%v", err)
		return nil, gin.NewInnerError(err.Error(), 500)
	}

	response.NonceStr = nonce
	response.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	response.Package = fmt.Sprintf("prepay_id=%s", response.PrepayId)
	response.SignType = "RSA"

	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", s.minaAppId, response.TimeStamp, response.NonceStr, response.Package)
	sign, err := s.client.Sign(ctx, message)
	if err != nil {
		s.Errorf("client.Sign(ctx, message) err:%v", err)
		return nil, gin.NewInnerError(err.Error(), 500)
	}
	response.PaySign = sign.Signature
	return &response, nil
}

func (s *paySvc) CloseOrderByOrderNo(orderNo string) error {
	svc := jsapi.JsapiApiService{Client: s.client}
	resp, err := svc.CloseOrder(context.Background(), jsapi.CloseOrderRequest{
		OutTradeNo: core.String(orderNo),
		Mchid:      core.String(s.mchID),
	})
	if err != nil {
		s.Errorf("wechat CloseOrderByOrderNo err:%v", err)
		return gin.NewInnerError(err.Error(), 500)
	}
	if resp.Response.StatusCode == http.StatusNoContent {
		return nil
	}

	return gin.NewInnerError(fmt.Sprintf("CloseOrder %s error with status code %d", orderNo, resp.Response.StatusCode), 500)
}

// QueryOrderInfoByOrderNo 订单查询
func (s *paySvc) QueryOrderInfoByOrderNo(orderNo string) (*payments.Transaction, bool, error) {
	svc := jsapi.JsapiApiService{Client: s.client}

	resp, result, err := svc.QueryOrderByOutTradeNo(context.Background(), jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(orderNo),
		Mchid:      core.String(s.mchID),
	})
	if err != nil {
		if apiErr, ok := err.(*core.APIError); ok {
			if apiErr.StatusCode == 404 {
				return nil, false, nil
			}
		}

		s.Errorf("wechat QueryOrderByOutTradeNo err:%v", err)
		return nil, false, gin.NewInnerError(err.Error(), 500)
	}
	if result.Response.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}

	if resp.TradeState == nil {
		s.Infof("wechat QueryOrderByOutTradeNo result: %v", result)
		return nil, false, gin.NewInnerError("TradeState is empty", 500)
	}
	return resp, true, nil
}

func (s *paySvc) QueryOrderInfoByWxTransactionId(wxTransactionId string) (*payments.Transaction, error) {
	svc := jsapi.JsapiApiService{Client: s.client}
	resp, result, err := svc.QueryOrderById(context.Background(), jsapi.QueryOrderByIdRequest{
		TransactionId: core.String(wxTransactionId),
		Mchid:         core.String(s.mchID),
	})

	if err != nil {
		s.Errorf("wechat QueryOrderById err:%v", err)
		return nil, gin.NewInnerError(err.Error(), 500)
	}
	if resp.TransactionId == nil {
		s.Infof("wechat QueryOrderById result: %v", result)
		return nil, gin.NewInnerError("TransactionId is empty", 500)
	}
	return resp, nil
}

func (s *paySvc) DecryPayNoticeInfo(request *http.Request) (*payments.Transaction, error) {
	transaction := new(payments.Transaction)
	_, err := s.handler.ParseNotifyRequest(context.Background(), request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		s.Errorf("ParseNotifyRequest err:%v", err)
		return nil, err
	}
	return transaction, nil
}

// TransferByOpenId 通过openID转账
func (s *paySvc) TransferByOpenId(req *entity.WxTransferRequest) (*transferbatch.InitiateBatchTransferResponse, error) {
	svc := transferbatch.TransferBatchApiService{Client: s.client}
	resp, result, err := svc.InitiateBatchTransfer(context.Background(),
		transferbatch.InitiateBatchTransferRequest{
			Appid:       core.String(s.minaAppId),
			OutBatchNo:  core.String(req.OrderNo),
			BatchName:   core.String(req.OrderName),
			BatchRemark: core.String(req.OrderRemark),
			TotalAmount: core.Int64(req.Amount),
			TotalNum:    core.Int64(1),
			TransferDetailList: []transferbatch.TransferDetailInput{{
				OutDetailNo:    core.String(req.OrderNo),
				TransferAmount: core.Int64(req.Amount),
				TransferRemark: core.String(req.OrderRemark),
				Openid:         core.String(req.Openid),
			}},
		},
	)

	if err != nil {
		// 处理错误
		s.Errorf("call InitiateBatchTransfer err:%s", err)
		return nil, err
	}
	// 处理返回结果
	s.Infof("status=%d resp=%s", result.Response.StatusCode, resp)
	return resp, nil
}

// GetTransferByBatchId 转账单查询
func (s *paySvc) GetTransferByBatchId(batchId string) (*transferbatch.TransferBatchEntity, error) {
	svc := transferbatch.TransferBatchApiService{Client: s.client}
	resp, result, err := svc.GetTransferBatchByNo(context.Background(),
		transferbatch.GetTransferBatchByNoRequest{
			BatchId:         core.String(batchId),
			NeedQueryDetail: core.Bool(true),
			Offset:          core.Int64(0),
			Limit:           core.Int64(20),
			DetailStatus:    core.String("FAIL"),
		},
	)

	if err != nil {
		// 处理错误
		s.Errorf("call GetTransferBatchByNo err:%s", err)
		return nil, err
	}
	// 处理返回结果
	s.Infof("status=%d resp=%s", result.Response.StatusCode, resp)
	return resp, nil
}
