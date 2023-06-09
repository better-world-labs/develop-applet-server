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
	return &gptConversation{}
}

type gptConversation struct {
	*Base `gone:"*"`

	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter              `gone:"router-auth"`
	PubRouter     gin.IRouter              `gone:"router-pub"`
	svc           service.IGPTConversation `gone:"*"`
}

func (con *gptConversation) Mount() gin.MountError {
	con.AuthRouter.
		POST("/gpt-messages", con.sendGptMessage).
		POST("/gpt-messages/:messageId/like", con.likeMessage).
		GET("/gpt-conversations", con.listGptChatMessages)

	return nil
}

func (con *gptConversation) listGptChatMessages(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	var query page.StreamQuery
	err := query.BindQuery(ctx)
	if err != nil {
		return nil, err
	}

	apps, err := con.svc.ListMessages(query, int64(userId))
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (con *gptConversation) sendGptMessage(ctx *gin.Context) (any, error) {
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

	userId := utils.CtxMustGetUserId(ctx)
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

	reader, err := con.svc.SendMessage(userId, param.Content)
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

func (con *gptConversation) likeMessage(ctx *gin.Context) (any, error) {
	messageId := ctx.Param("messageId")

	var param struct {
		LikeState entity.LikeState `json:"likeState"`
	}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, con.svc.LikeMessage(messageId, param.LikeState)
}
