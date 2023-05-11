package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewRetainMessageController() gone.Goner {
	return &retainMessageController{}
}

type retainMessageController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter            `gone:"router-auth"`
	retainMessage service.IRetainMessage `gone:"*"`
}

func (con *retainMessageController) Mount() gin.MountError {
	con.
		AuthRouter.
		GET("/retain-messages", con.listRetainMessages)

	return nil
}

func (con *retainMessageController) listRetainMessages(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	messages, err := con.retainMessage.ListRetainMessages(userId)
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{List: messages}, nil
}
