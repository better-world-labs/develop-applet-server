package middleware

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
)

//go:gone
func NewPubAuthorizeMiddleware() gone.Goner {
	return &PubAuthorizeMiddleware{}
}

type PubAuthorizeMiddleware struct {
	gone.Flag
	User          service.IUser `gone:"*"`
	logrus.Logger `gone:"gone-logger"`
}

func (m *PubAuthorizeMiddleware) Next(ctx *gin.Context) (any, error) {
	if utils.CtxGetString(ctx, entity.JwtKey) == "" {
		ctx.Next()
		return nil, nil
	}

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
