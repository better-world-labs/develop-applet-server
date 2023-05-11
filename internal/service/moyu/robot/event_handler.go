package robot

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/imroc/req"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"
)

type EventHandler struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`

	Svc     service.IRobot         `gone:"*"`
	Channel service.IChannel       `gone:"*"`
	Message service.IMessageRecord `gone:"*"`
}

//go:gone
func NewEventHandler() gone.Goner {
	return &EventHandler{}
}

func (e *EventHandler) Consume(on emitter.OnEvent) {
	on(e.handleMentioned)
	on(e.handleImage)
}

func (e *EventHandler) handleImage(evt *entity.ImageMessageSavedEvent) error {
	e.Logger.Info("[robot] handleImage: evt=", evt)

	configs, err := e.Svc.ListConfigsByTrigger(entity.TriggerTypeImage)
	if err != nil {
		return err
	}

	return e.sendRobotMessage(configs, evt.Origin)
}

func (e *EventHandler) handleMentioned(evt *entity.UserMentioned) error {
	e.Logger.Info("[robot] handleMentioned: evt=", evt)

	configs, err := e.Svc.ListConfigsByUserIds(evt.TargetUsers)
	if err != nil {
		return err
	}

	return e.sendRobotMessage(configs, evt.Msg)
}

func (e *EventHandler) postMessage(c *entity.RobotConfig, context *entity.RobotContext) error {
	b, _ := json.Marshal(context)
	e.Info("postMessage: appId=", c.RobotId, "message=", string(b))

	resp, err := req.Post(c.MessageReceiveUrl, req.Header{
		"App-Id": c.RobotId,
		"Token":  e.createToken(c),
	}, req.BodyJSON(context))
	if err != nil {
		return err
	}

	var r struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	e.Infof("postMessage response: ", resp.String())
	err = resp.ToJSON(&r)
	if err != nil {
		return err
	}

	if r.Code != 0 {
		return errors.New(r.Msg)
	}
	return nil
}

func (e *EventHandler) createToken(c *entity.RobotConfig) string {
	b := bytes.Buffer{}
	b.WriteString(c.RobotId)
	b.WriteString("|")

	h := md5.New()
	h.Write(b.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

func (e *EventHandler) sendRobotMessage(configs []*entity.RobotConfig, evt *wsevent.Msg) error {
	context, err := e.Message.GetContext(evt.Id, evt.ChannelId)
	if err != nil {
		return nil
	}

	e.Logger.Debugf("GetContent: len=%d\n", len(context))

	context = append(context, evt)
	for _, c := range configs {
		c := c
		if c.UserId == evt.UserId {
			continue
		}

		if c.ContextLength <= 0 || len(context) <= 0 {
			continue
		}

		go func() {
			err := e.postMessage(c, &entity.RobotContext{Robot: c.UserId, ChannelId: evt.ChannelId, Context: utils.CutTail(context, c.ContextLength)})
			if err != nil {
				e.Error("postMessage error: ", err)
			}

			_, err = e.Channel.UpdateLastReadMessage(c.UserId, evt.ChannelId)
			if err != nil {
				e.Error("UpdateLastReadMessage error: ", err)
			}
		}()
	}

	return nil
}
