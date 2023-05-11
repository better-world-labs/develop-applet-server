package router

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/middleware"
)

const IdRouterPub = "router-pub"

//go:gone
func NewPubRouter() (gone.Goner, gone.GonerId) {
	return &pubRouter{}, IdRouterPub
}

type pubRouter struct {
	gone.Flag
	gin.IRouter
	root  gin.IRouter                 `gone:"gone-gin-router"`
	Plant *middleware.PlantMiddleware `gone:"*"`
}

func (r *pubRouter) AfterRevive() gone.AfterReviveError {
	r.IRouter = r.root.Group("/api", r.Plant.Next)
	return nil
}
