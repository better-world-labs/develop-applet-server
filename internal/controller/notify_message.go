package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewNotifyMessageController() gone.Goner {
	return &notifyMessageController{}
}

type notifyMessageController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter            `gone:"router-auth"`
	Svc           service.INotifyMessage `gone:"*"`
}

func (ctr *notifyMessageController) Mount() gin.MountError {
	ctr.AuthRouter.
		GET("/notify-messages", ctr.pageNotifyMessages).
		GET("/notify-messages/unread-count", ctr.countUnreadMessages).
		PUT("/notify-messages/:id/read", ctr.markMessageRead).
		PUT("/notify-messages/read-all", ctr.markMessageReadAll)

	return nil
}

func (ctr *notifyMessageController) pageNotifyMessages(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)

	var query page.StreamQuery
	var filter entity.NotifyMessageListFilter

	err := query.BindQuery(ctx)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	err = ctx.ShouldBindQuery(&filter)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.Svc.PageNotifyMessages(userId, query, filter)
}

func (ctr *notifyMessageController) countUnreadMessages(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	count, err := ctr.Svc.CountUnread(userId)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return map[string]int64{
		"count": count,
	}, nil
}

func (ctr *notifyMessageController) markMessageRead(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	notifyMessageId, err := utils.CtxPathParamInt64(ctx, "id")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, ctr.Svc.MarkRead(userId, notifyMessageId)
}

func (ctr *notifyMessageController) markMessageReadAll(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	return nil, ctr.Svc.MarkReadAll(userId)
}
