package middleware

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewPlantMiddleware() gone.Goner {
	return &PlantMiddleware{}
}

type PlantMiddleware struct {
	gone.Flag
	CookieJwtKey      string        `gone:"config,cookie.jwt-key"`
	CookieClientIdKey string        `gone:"config,cookie.client-id-key"`
	User              service.IUser `gone:"*"`
	logrus.Logger     `gone:"gone-logger"`
}

func (m *PlantMiddleware) Next(ctx *gin.Context) (interface{}, error) {
	//在clientId不存在时，使jwt也失效

	clientId, _ := ctx.Cookie(m.CookieClientIdKey)
	jwt, _ := ctx.Cookie(m.CookieJwtKey)
	if clientId == "" {
		clientId = m.User.GenClientId()
		m.User.SetCookie(ctx, m.CookieClientIdKey, clientId)

		if jwt != "" {
			jwt = ""
			m.User.SetCookie(ctx, m.CookieJwtKey, jwt)
		}
	}

	ctx.Set(entity.ClientIdKey, clientId)
	ctx.Set(entity.JwtKey, jwt)

	return nil, nil
}
