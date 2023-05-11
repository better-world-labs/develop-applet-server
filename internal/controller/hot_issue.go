package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewHotIssueController() gone.Goner {
	return &hotIssueController{}
}

type hotIssueController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter       `gone:"router-pub"`
	Svc           service.IHotIssue `gone:"*"`
}

func (ctr *hotIssueController) Mount() gin.MountError {
	ctr.AuthRouter.
		GET("/hot-issues", ctr.listIssues)

	return nil
}

func (ctr *hotIssueController) listIssues(*gin.Context) (interface{}, error) {
	list, err := ctr.Svc.ListIssues()
	if err != nil {
		return nil, gin.ToError(err)
	}

	return entity.ListWrap{List: list}, nil
}
