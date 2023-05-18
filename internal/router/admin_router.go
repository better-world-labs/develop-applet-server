package router

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/tracer"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/middleware"
)

const IdAdminRouter = "router-admin"

//go:gone
func NewAdmin() (gone.Goner, gone.GonerId) {
	return &adminRouter{}, IdAdminRouter
}

type adminRouter struct {
	gone.Flag
	gin.IRouter
	logrus.Logger `gone:"gone-logger"`
	root          gin.IRouter `gone:"gone-gin-router"`

	Auth   *middleware.AdminMiddleware `gone:"*"`
	tracer tracer.Tracer               `gone:"gone-tracer"`
}

func (r *adminRouter) AfterRevive() gone.AfterReviveError {
	r.IRouter = r.root.Group("/admin-api", r.Auth.Next)
	return nil
}
