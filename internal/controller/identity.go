package controller

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type IdentityController struct {
	gone.Flag
	R   gin.IRouter                `gone:"router-pub"`
	Svc service.IAnonymousIdentity `gone:"*"`
}

//go:gone
func NewIdentityController() gone.Goner {
	return &IdentityController{}
}

func (ctr *IdentityController) Mount() gin.MountError {
	ctr.R.Group("/anonymous-identities").
		GET("", ctr.list)

	return nil
}

func (ctr *IdentityController) list(*gin.Context) (interface{}, error) {
	list, err := ctr.Svc.List()
	if err != nil {
		return nil, gin.ToError(err)
	}

	return entity.ListWrap{List: list}, nil
}
