package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewSignController() gone.Goner {
	return &signIn{}
}

type signIn struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter     `gone:"router-auth"`
	Sign          service.ISignIn `gone:"*"`
}

func (con *signIn) Mount() gin.MountError {
	con.
		AuthRouter.
		GET("/sign-in", con.getSignIn).
		POST("/sign-in", con.signIn)

	return nil
}

func (con *signIn) getSignIn(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	isSignIn, err := con.Sign.GetSignInStatus(userId)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"signIn": isSignIn,
	}, nil
}

func (con *signIn) signIn(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	return nil, con.Sign.SignIn(userId)
}
