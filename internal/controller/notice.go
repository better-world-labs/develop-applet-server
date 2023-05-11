package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	businesserrors "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/error"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewNoticeController() gone.Goner {
	return &noticeController{}
}

type noticeController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	AuthRouter gin.IRouter     `gone:"router-auth"`
	Notice     service.INotice `gone:"*"`
}

func (ctr *noticeController) Mount() gin.MountError {
	ctr.AuthRouter.
		Group("/notices").
		GET("", ctr.pageNotices).
		GET("/:id", ctr.getNotice).
		GET("/unread-count", ctr.countUnread).
		POST("/read", ctr.markRead)

	return nil
}

func (ctr *noticeController) pageNotices(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	query, err := page.ParseQuery(ctx)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.Notice.Page(userId, query)
}

func (ctr *noticeController) getNotice(ctx *gin.Context) (interface{}, error) {
	id, err := utils.CtxPathParamInt64(ctx, "id")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	res, exists, err := ctr.Notice.Get(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, businesserrors.ErrorNoticeNotFound
	}

	return res, nil

}

func (ctr *noticeController) countUnread(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	unread, err := ctr.Notice.CountUnread(userId)
	if err != nil {
		return nil, err
	}

	return map[string]any{"count": unread}, nil
}

func (ctr *noticeController) markRead(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	var param struct {
		Ids []int64 `json:"ids" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, ctr.Notice.MarkRead(userId, param.Ids)
}
