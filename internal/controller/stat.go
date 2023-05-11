package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewStatController() gone.Goner {
	return &statController{}
}

type statController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	gin.IRouter   `gone:"router-auth"`
	stat          service.IMessageStat `gone:"*"`
}

func (ctr *statController) Mount() gin.MountError {
	ctr.
		Group("/stat").
		GET("/hot-messages", ctr.getHotMessages)
	return nil
}

func (ctr *statController) getHotMessages(ctx *gin.Context) (any, error) {
	type Req struct {
		Top      int   `form:"top"`
		ChanelId int64 `form:"channel" binding:"required"`
	}

	var req Req
	if err := ctx.BindQuery(&req); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}
	return utils.List(ctr.stat.HotMessageTop(req.ChanelId, req.Top))
}
