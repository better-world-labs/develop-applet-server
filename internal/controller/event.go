package controller

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewEventController() gone.Goner {
	return &eventController{}
}

type eventController struct {
	gone.Flag
	logrus.Logger  `gone:"gone-logger"`
	emitter.Sender `gone:"gone-emitter"`

	AuthRouter gin.IRouter   `gone:"router-auth"`
	User       service.IUser `gone:"*"`
}

func (con *eventController) Mount() gin.MountError {
	con.
		AuthRouter.
		Group("/events").
		POST("", con.triggerEvent)

	return nil
}

func (con *eventController) triggerEvent(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	evt := entity.UserEvent{UserId: userId}

	err := ctx.ShouldBindJSON(&evt)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, con.User.PostEvent(&evt)
}
