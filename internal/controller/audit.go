package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewAuditController() gone.Goner {
	return &auditController{}
}

type auditController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter           `gone:"router-pub"`
	Audit         service.IContentAudit `gone:"*"`
}

func (ctr *auditController) Mount() gin.MountError {
	ctr.
		AuthRouter.
		Group("/audit").
		POST("/text", ctr.textScan)

	return nil
}

func (ctr *auditController) textScan(ctx *gin.Context) (interface{}, error) {
	var param struct {
		Text string `json:"text" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	result, err := ctr.Audit.ScanText(param.Text)
	if err != nil {
		return nil, gin.ToError(err)
	}

	return map[string]any{
		"valid": result.CheckPass(),
	}, nil
}
