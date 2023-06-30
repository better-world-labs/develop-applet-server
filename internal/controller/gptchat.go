package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"io"
	"net/http"
)

//go:gone
func NewGPTChatController() gone.Goner {
	return &gptChat{}
}

type gptChat struct {
	*Base `gone:"*"`

	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter      `gone:"router-pub"`
	svc           service.IGPTChat `gone:"*"`
}

func (con *gptChat) Mount() gin.MountError {
	con.AuthRouter.
		GET("/gpt-conversations", con.listGptChatMessages).
		POST("gpt-messages", con.sendGptMessage)
	return nil
}

func (con *gptChat) listGptChatMessages(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)

	var query page.StreamQuery
	err := query.BindQuery(ctx)
	if err != nil {
		return nil, err
	}

	apps, err := con.svc.ListMessages(query, userId)
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{List: apps}, nil
}

func (con *gptChat) sendGptMessage(ctx *gin.Context) (any, error) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")

	var done struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
	}

	f, ok := ctx.Writer.(http.Flusher)
	if !ok {
		return nil, nil
	}

	//userId := utils.CtxMustGetUserId(ctx)
	userId := 10
	var param struct {
		Content string `json:"content" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		done.Code = -1
		done.Msg = err.Error()
		_, _ = io.WriteString(ctx.Writer, "event: done\n")
		doneJson, _ := json.Marshal(done)
		_, _ = io.WriteString(ctx.Writer, fmt.Sprintf("data: %s\n\n", doneJson))
		f.Flush()
		return nil, nil
	}

	reader, err := con.svc.SendMessage(int64(userId), param.Content)
	if err != nil {
		if goneErr, ok := err.(gone.Error); ok {
			done.Code = goneErr.Code()
			done.Msg = goneErr.Error()
		} else {
			done.Code = -1
			done.Msg = err.Error()
		}

		_, _ = io.WriteString(ctx.Writer, "event: done\n")
		doneJson, _ := json.Marshal(done)
		_, _ = io.WriteString(ctx.Writer, fmt.Sprintf("data: %s\n\n", doneJson))
		f.Flush()
		return nil, nil
	}

	for {
		read, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			done.Code = -1
			done.Msg = err.Error()
			break
		}

		json, err := json.Marshal(read)
		if err != nil {
			done.Code = -1
			done.Msg = err.Error()
			break
		}

		_, _ = io.WriteString(ctx.Writer, "event: data\n")
		_, _ = io.WriteString(ctx.Writer, fmt.Sprintf("data: %s\n\n", json))
		f.Flush()
	}

	_, _ = io.WriteString(ctx.Writer, "event: done\n")
	doneJson, err := json.Marshal(done)
	_, _ = io.WriteString(ctx.Writer, fmt.Sprintf("data: %s\n\n", doneJson))
	f.Flush()
	reader.Close()

	return nil, nil
}
