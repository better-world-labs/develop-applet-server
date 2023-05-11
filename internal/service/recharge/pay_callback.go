package recharge

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"net/http"
)

type payCallback struct {
	logrus.Logger `gone:"gone-logger"`

	PubRouter gin.IRouter `gone:"router-pub"`
	gone.Goner

	callbackPath string         `gone:"config,wechat.pay.notify-path"`
	pay          service.IWxPay `gone:"*"`
	order        *orderSvc      `gone:"*"`
}

//go:gone
func NewPayCallback() gone.Goner {
	return &payCallback{}
}

func (c *payCallback) Mount() gin.MountError {
	c.PubRouter.POST(c.callbackPath, c.callback)
	return nil
}

func (c *payCallback) callback(ctx *gin.Context) (any, error) {
	notifyInfo, err := c.pay.DecryPayNoticeInfo(ctx.Request)
	if err != nil {
		c.Errorf("wechatPaymentCallback: DecryPayNoticeInfo error: %v\n ", err)
		responseError(ctx)
		return nil, nil
	}

	err = c.order.handlePayNotify(notifyInfo)
	if err != nil {
		c.Errorf("wechatPaymentCallback: HandlePayNotify: error: %v\n ", err)
		responseError(ctx)
		return nil, nil
	}

	responseOk(ctx)
	return nil, nil
}

func responseError(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{"FAIL", "失败"})
}

func responseOk(ctx *gin.Context) {
	ctx.Status(200)
}
