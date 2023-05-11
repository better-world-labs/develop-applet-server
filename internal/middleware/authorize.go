package middleware

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"net/http"
)

//go:gone
func NewAuthorizeMiddleware() gone.Goner {
	return &AuthorizeMiddleware{}
}

type AuthorizeMiddleware struct {
	gone.Flag
	User          service.IUser `gone:"*"`
	logrus.Logger `gone:"gone-logger"`
}

func (m *AuthorizeMiddleware) Next(ctx *gin.Context) (any, error) {
	userId, err := m.User.ParseJwtInfo(ctx)

	if err != nil {
		failed(ctx, err.Code(), err.Msg())
		return nil, nil
	}

	m.Infof("userId:%v", userId)
	ctx.Set(entity.UserIdKey, userId)
	ctx.Next()
	return nil, nil
}

func failed(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
		"code": code,
		"msg":  message,
	})

	ctx.Abort()
}
