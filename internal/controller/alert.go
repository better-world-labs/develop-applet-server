package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewAlertController() gone.Goner {
	return &alertController{}
}

type alertController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter    `gone:"router-auth"`
	Alert         service.IAlert `gone:"*"`
}

func (ctr *alertController) Mount() gin.MountError {
	ctr.
		AuthRouter.
		Group("/alerts").
		GET("/version", ctr.alertVersion)

	return nil
}

func (ctr *alertController) alertVersion(*gin.Context) (interface{}, error) {
	return ctr.Alert.GetAvailableAlert()
}
