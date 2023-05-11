package controller

import (
	"encoding/json"
	"errors"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"io"
	"strconv"
)

//go:gone
func NewRobotController() gone.Goner {
	return &robotController{}
}

type robotController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	gin.IRouter   `gone:"router-pub"`

	Robot service.IRobot `gone:"*"`
	env   string         `gone:"config,server.env"`
}

func (ctr *robotController) Mount() gin.MountError {
	ctr.
		Group("/robots").
		POST("/messages", ctr.postMessage)

	return nil
}

func (ctr *robotController) postMessage(ctx *gin.Context) (interface{}, error) {
	content, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}

	env := ctx.GetHeader("Env")
	robotId := ctx.GetHeader("App-Id")
	channelId, _ := strconv.ParseInt(ctx.GetHeader("Channel-Id"), 10, 64)
	if !json.Valid(content) {
		return nil, errors.New("invalid request body")
	}

	err = ctr.Robot.SendMessage(env, robotId, channelId, content)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
