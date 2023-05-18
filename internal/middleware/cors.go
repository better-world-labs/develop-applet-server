package middleware

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"net/http"
)

//go:gone
func NewCorsMiddleware() gone.Goner {
	return &CorsMiddleware{}
}

type CorsMiddleware struct {
	gone.Flag
	CookieJwtKey      string        `gone:"config,cookie.jwt-key"`
	CookieClientIdKey string        `gone:"config,cookie.client-id-key"`
	User              service.IUser `gone:"*"`
	logrus.Logger     `gone:"gone-logger"`
}

func (m *CorsMiddleware) Next(c *gin.Context) (interface{}, error) {
	method := c.Request.Method

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Max-Age", "1800")

	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	// 处理请求
	c.Next()
	return nil, nil
}
