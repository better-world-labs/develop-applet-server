package middleware

import (
	"errors"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"net/http"
	"strings"
)

//go:gone
func NewAdminMiddleware() gone.Goner {
	return &AdminMiddleware{}
}

type AdminMiddleware struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	user          service.IUser `gone:"*"`
	secret        string        `gone:"config,auth.admin.secret"`
}

func (m *AdminMiddleware) CheckToken(ctx *gin.Context) (int64, error) {
	authorization := ctx.GetHeader("Authorization")
	bearerToken := strings.Split(authorization, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return 0, errors.New("invalid authorization")
	}

	id, err := m.user.ParseJwt(bearerToken[1], m.secret)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *AdminMiddleware) Next(ctx *gin.Context) (any, error) {
	_, err := m.CheckToken(ctx)
	if err != nil {
		m.Abort(ctx, -1, err.Error())
		return nil, nil
	}

	ctx.Next()
	return nil, nil
}

func (m *AdminMiddleware) Abort(ctx *gin.Context, code int, msg string) {
	ctx.JSON(http.StatusUnauthorized, map[string]any{
		"code": code,
		"msg":  msg,
	})
	ctx.Abort()
}
