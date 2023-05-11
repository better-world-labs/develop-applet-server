package controller

import (
	"encoding/json"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

//go:gone
func NewConfig() gone.Goner {
	return &config{}
}

type config struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	PubRouter     gin.IRouter `gone:"router-pub"`
	InnerRouter   gin.IRouter `gone:"router-inner"`

	sys service.ISystemConfig `gone:"*"`
}

func (ctr *config) Mount() gin.MountError {
	ctr.PubRouter.
		GET("/system-configs", ctr.get)

	ctr.InnerRouter.
		PUT("/system-configs", ctr.put)

	return nil
}

func (ctr *config) get(ctx *gin.Context) (any, error) {
	key, exists := ctx.GetQuery("key")
	if !exists {
		return nil, gin.NewParameterError("key cannot empty")
	}

	value, err := ctr.sys.Get(key)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"value": value,
	}, nil
}

func (ctr *config) put(ctx *gin.Context) (any, error) {
	var param struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError("key cannot empty")
	}

	return nil, ctr.sys.Put(param.Key, param.Value)
}
