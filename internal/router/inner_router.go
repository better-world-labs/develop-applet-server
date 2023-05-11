package router

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/middleware"
)

const IdRouterInner = "router-inner"

//go:gone
func NewInnerRouter() (gone.Goner, gone.GonerId) {
	return &innerRouter{}, IdRouterInner
}

type innerRouter struct {
	gone.Flag
	gin.IRouter
	root  gin.IRouter                 `gone:"gone-gin-router"`
	Plant *middleware.PlantMiddleware `gone:"*"`
}

func (r *innerRouter) AfterRevive() gone.AfterReviveError {
	r.IRouter = r.root.Group("/inner", r.Plant.Next)
	return nil
}
