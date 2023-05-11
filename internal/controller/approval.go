package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	businesserrors "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/error"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewApprovalController() gone.Goner {
	return &approvalController{}
}

type approvalController struct {
	*Base `gone:"*"`

	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter       `gone:"router-auth"`
	svc           service.IApproval `gone:"*"`
}

func (ctr *approvalController) Mount() gin.MountError {
	ctr.AuthRouter.
		Group("/approvals").
		GET("/:id", ctr.getOne).
		POST("/:id/audit", ctr.WithRole(entity.PlanetRoleAdmin, ctr.audit))

	return nil
}

func (ctr *approvalController) audit(ctx *gin.Context) (any, error) {
	id, err := utils.CtxPathParamInt64(ctx, "id")
	userId := utils.CtxMustGetUserId(ctx)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	param := struct {
		Pass bool `json:"pass"`
	}{}

	err = ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, ctr.svc.Audit(id, userId, param.Pass)
}

func (ctr *approvalController) getOne(ctx *gin.Context) (any, error) {
	id, err := utils.CtxPathParamInt64(ctx, "id")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	approval, exists, err := ctr.svc.GetOne(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, businesserrors.ErrorApprovalNotFound
	}

	return approval, nil
}
