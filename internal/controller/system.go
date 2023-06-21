package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"strconv"
)

//go:gone
func NewSystemController() gone.Goner {
	return &systemController{}
}

type systemController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter     `gone:"router-auth"`
	PubRouter     gin.IRouter     `gone:"router-pub"`
	System        service.ISystem `gone:"*"`
	Jssdk         service.IJSSDK  `gone:"*"`
}

func (ctr *systemController) Mount() gin.MountError {
	ctr.AuthRouter.
		Group("/system").
		GET("/oss-token", ctr.genOssToken).
		GET("/emoticons", ctr.getEmoticonList)

	ctr.PubRouter.
		Group("/system").
		POST("/js-sdk-signature", ctr.getJSSignature)

	return nil
}

func (ctr *systemController) genOssToken(context *gin.Context) (interface{}, error) {
	fileExt := context.Query("ext")
	return ctr.System.GenOssToken(utils.CtxMustGetUserId(context), fileExt)
}

func (ctr *systemController) getEmoticonList(ctx *gin.Context) (data any, err error) {
	group := ctx.Query("group")
	sort := ctx.Query("sort")

	var groupId int64
	if group != "" {
		groupId, err = strconv.ParseInt(group, 10, 64)
		if err != nil {
			return nil, gin.NewParameterError(err.Error())
		}
	}
	return utils.List(ctr.System.GetEmoticonList(int(groupId), sort != "none"))
}

func (ctr *systemController) getJSSignature(ctx *gin.Context) (any, error) {
	var param struct {
		Url string `json:"url" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.Jssdk.GetSignature(param.Url)
}
