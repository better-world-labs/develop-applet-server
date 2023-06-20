package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewAITool() gone.Goner {
	return &aiTool{}
}

type aiTool struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	PubRouter     gin.IRouter     `gone:"router-pub"`
	svc           service.IAITool `gone:"*"`
}

func (con *aiTool) Mount() gin.MountError {
	con.PubRouter.
		GET("/ai-tools", con.listAITools).
		GET("/ai-tool-categories", con.listAIToolCategories)

	return nil
}

func (con *aiTool) listAITools(ctx *gin.Context) (any, error) {
	tools, err := con.svc.List()
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{
		List: tools,
	}, nil
}

func (con *aiTool) listAIToolCategories(ctx *gin.Context) (any, error) {
	categories, err := con.svc.ListCategories()
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{
		List: categories,
	}, nil
}
