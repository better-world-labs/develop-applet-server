package entity

import (
	"time"
)

type WxPayOrderRequest struct {
	OrderNo     string    //外部订单号
	Description string    //商品描述
	TimeExpire  time.Time //订单过期时间
	Total       int64     //金额，单位 分
	Openid      string    //用户openID
	Attach      string    //附加信息；附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用，实际情况下只有支付完成状态才会返回该字段。
	NotifyUrl   string    //异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。
}

type WxPayOrderResponse struct {
	PrepayId  string `json:"prepayId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type WxTransferRequest struct {
	OrderNo     string //单号  商户系统内部的商家批次单号，要求此参数只能由数字、大小写字母组成，在商户系统内部唯一
	OrderName   string //单名  该笔批量转账的名称
	OrderRemark string //备注  转账说明，UTF8编码，最多允许32个字符
	Amount      int64  //转账金额，单位 分
	Openid      string //收款用户对应openID
}
