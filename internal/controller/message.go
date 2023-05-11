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
func NewMessageController() gone.Goner {
	return &messageController{}
}

type messageController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	AuthRouter gin.IRouter            `gone:"router-auth"`
	Message    service.IMessageRecord `gone:"*"`
}

func (ctr *messageController) Mount() gin.MountError {
	ctr.
		AuthRouter.
		Group("/messages").
		POST("get", ctr.getMany).
		POST("/history", ctr.getHistoryMsg).
		POST("/:id/like", ctr.like).
		POST("/like", ctr.getLikes)

	return nil
}

func (ctr *messageController) getMany(ctx *gin.Context) (interface{}, error) {
	var param struct {
		Ids []int64 `json:"ids" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, err
	}

	records, err := ctr.Message.GetRecords(param.Ids)
	if err != nil {
		return nil, gone.NewInnerError(500, err.Error())
	}

	return entity.ListWrap{List: records}, nil
}

func (ctr *messageController) getHistoryMsg(ctx *gin.Context) (interface{}, error) {
	var param struct {
		ChannelId int64 `json:"channelId"`
		FromId    int64 `json:"fromId"`
		Size      int   `json:"size"`
		UpFlag    bool  `json:"upFlag"`
	}

	err := ctx.BindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	history, err := ctr.Message.ListHistory(param.ChannelId, param.FromId, param.Size, param.UpFlag)
	return entity.ListWrap{List: history}, err
}

func (ctr *messageController) like(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	id, err := utils.CtxPathParamInt64(ctx, "id")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	var param struct {
		IsLike bool `json:"isLike"`
	}

	err = ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, ctr.Message.Like(id, userId, param.IsLike)
}

func (ctr *messageController) getLikes(ctx *gin.Context) (any, error) {
	var param struct {
		Ids []int64 `json:"ids" binding:"required"`
	}

	userId := utils.CtxMustGetUserId(ctx)
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	likes, err := ctr.Message.GetMessageLikes(userId, param.Ids)
	return entity.ListWrap{
		List: likes,
	}, err
}
