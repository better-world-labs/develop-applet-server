package middleware

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"net/http"
)

//go:gone
func NewCorsMiddleware() gone.Goner {
	return &CorsMiddleware{}
}

type CorsMiddleware struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
}

func (m *CorsMiddleware) Next(c *gin.Context) (interface{}, error) {
	m.Infof("handle CorsMiddleware request=%s\n", c.Request.URL)
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
		return nil, nil
	}
	// 处理请求
	c.Next()
	return nil, nil
}
