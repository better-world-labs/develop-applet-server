package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewUserExtra() gone.Goner {
	return &userExtraController{}
}

type userExtraController struct {
	gone.Flag
	logrus.Logger    `gone:"gone-logger"`
	AuthRouter       gin.IRouter               `gone:"router-auth"`
	MiniAppUserExtra service.IMiniAppUserExtra `gone:"*"`
}

func (con *userExtraController) Mount() gin.MountError {
	con.
		AuthRouter.
		Group("/users/me").
		POST("/guidance/completion", con.completeGuidance).
		GET("/guidance", con.getGuidanceStatus)

	return nil
}

func (con *userExtraController) completeGuidance(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	return nil, con.MiniAppUserExtra.CompleteGuidance(userId)
}

func (con *userExtraController) getGuidanceStatus(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	extra, _, err := con.MiniAppUserExtra.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"completed": extra.CompleteGuidance,
	}, nil
}
