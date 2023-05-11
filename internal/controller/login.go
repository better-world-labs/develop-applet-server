package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"time"
)

//go:gone
func NewLoginController() gone.Goner {
	return &loginController{}
}

type loginController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	PubRouter     gin.IRouter   `gone:"router-pub"`
	AuthRouter    gin.IRouter   `gone:"router-auth"`
	User          service.IUser `gone:"*"`

	CookieJwtKey string `gone:"config,cookie.jwt-key"`
	Env          string `gone:"config,server.env"`
}

func (ctr *loginController) Mount() gin.MountError {
	group := ctr.PubRouter.Group("/users/login")
	ctr.AuthRouter.Group("/users/logout").POST("", ctr.logout)

	group.
		POST("", ctr.loginByToken).
		GET("/qr", ctr.getQrToken).
		GET("/auth-params", ctr.getAuthParams).
		POST("/code", ctr.authByCode)

	if ctr.Env == "dev" {
		group.
			GET("/code-test", ctr.codeTest).
			POST("/test", ctr.testLogin).
			GET("/test", func(ctx *gin.Context) (any, error) {
				ctr.Infof("headers:%v", ctx.Request.Header)

				select {
				case <-time.After(10 * time.Second):
					ctr.Infof("End:::::time end")
				case <-ctx.Request.Context().Done():
					ctr.Infof("XXXX::::Cancel")
				}
				ctr.Infof("end")
				return "ok", nil
			})
	}
	return nil
}

func (ctr *loginController) getQrToken(*gin.Context) (interface{}, error) {
	return ctr.User.GenQrToken()
}

func (ctr *loginController) loginByToken(context *gin.Context) (interface{}, error) {
	type Req struct {
		LoginToken string     `json:"loginToken"  binding:"required"`
		NickName   string     `json:"nickname"`
		Avatar     entity.Url `json:"avatar"`
		InvitedBy  *int64     `json:"invitedBy"`
		FromApp    string     `json:"fromApp"`
		Source     string     `json:"source"`
	}

	var req Req
	if err := context.BindJSON(&req); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	user, err := ctr.User.SetUserWithLoginToken(req.LoginToken, req.NickName, req.Avatar, req.InvitedBy, req.FromApp, req.Source)
	if err != nil {
		return nil, err
	}
	return ctr.User.Login(context, user)
}

func (ctr *loginController) getAuthParams(context *gin.Context) (interface{}, error) {
	type Req struct {
		QrToken     string `form:"qrToken"  binding:"required"`
		RedirectUrl string `form:"redirectUrl"  binding:"required"`
	}

	var req Req
	if err := context.BindQuery(&req); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	expired, url, err := ctr.User.GetQrTokenParams(req.QrToken, req.RedirectUrl)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tokenExpired": expired,
		"authUrl":      url,
	}, nil
}

func (ctr *loginController) authByCode(context *gin.Context) (interface{}, error) {
	type Req struct {
		QrToken string `json:"qrToken"  binding:"required"`
		Code    string `json:"code"  binding:"required"`
	}

	var req Req
	if err := context.BindJSON(&req); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}
	return ctr.User.AuthByCode(req.QrToken, req.Code)
}

func (ctr *loginController) codeTest(context *gin.Context) (interface{}, error) {
	_, _ = context.Writer.WriteString(context.Request.RequestURI)
	return nil, nil
}

func (ctr *loginController) testLogin(context *gin.Context) (interface{}, error) {
	type Req struct {
		OpenId string `json:"openId"  binding:"required"`
	}
	var req Req
	if err := context.BindJSON(&req); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	user, err := ctr.User.GetUserByOpenId(req.OpenId)
	if err != nil {
		return nil, err
	}
	return ctr.User.Login(context, user)
}

func (ctr *loginController) logout(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	clientId := utils.CtxGetString(ctx, entity.ClientIdKey)
	return nil, ctr.User.Logout(ctx, userId, clientId)
}
