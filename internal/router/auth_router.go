package router

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/tracer"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/middleware"
)

const IdAuthRouter = "router-auth"

//go:gone
func NewAuth() (gone.Goner, gone.GonerId) {
	return &authRouter{}, IdAuthRouter
}

type authRouter struct {
	gone.Flag
	gin.IRouter
	logrus.Logger `gone:"gone-logger"`
	root          gin.IRouter `gone:"gone-gin-router"`

	Auth   *middleware.AuthorizeMiddleware `gone:"*"`
	Plant  *middleware.PlantMiddleware     `gone:"*"`
	tracer tracer.Tracer                   `gone:"gone-tracer"`
}

func (r *authRouter) AfterRevive() gone.AfterReviveError {
	r.IRouter = r.root.Group("/api", r.Plant.Next, r.Auth.Next)
	return nil
}
