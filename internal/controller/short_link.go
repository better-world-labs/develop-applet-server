package controller

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewShortLinkController() gone.Goner {
	return &shortLinkController{}
}

type shortLinkController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	R             gin.IRouter        `gone:"router-pub"`
	Svc           service.IShortLink `gone:"*"`
	Host          string             `gone:"config,server.domain"`
}

func (ctr *shortLinkController) Mount() gin.MountError {

	ctr.
		R.
		Group("/l").
		POST("/short-link", ctr.generateShortLink).
		GET("/:linkCode", ctr.proxyLink)

	return nil
}

func (ctr *shortLinkController) generateShortLink(ctx *gin.Context) (interface{}, error) {
	var param struct {
		Url string `json:"url" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	code, err := ctr.Svc.Create(param.Url)
	if err != nil {
		return nil, gin.ToError(err)
	}

	return fmt.Sprintf("https://%s/api/l/%s", ctr.Host, code), nil
}

func (ctr *shortLinkController) proxyLink(ctx *gin.Context) (interface{}, error) {
	linkCode := ctx.Param("linkCode")

	origin, exists, err := ctr.Svc.GetOrigin(linkCode)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if !exists {
		return "链接失效", nil
	}

	ctx.Redirect(302, origin)
	return nil, nil
}
