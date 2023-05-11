package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewResignController() gone.Goner {
	return &resignController{}
}

type resignController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter     `gone:"router-auth"`
	ResignService service.IResign `gone:"*"`
}

func (ctr *resignController) Mount() gin.MountError {

	ctr.
		AuthRouter.
		Group("/resign").
		GET("/templates", ctr.listTemplates)

	return nil
}

func (ctr *resignController) listTemplates(*gin.Context) (interface{}, error) {
	return ctr.ResignService.ListTemplates()
}
