package miniapp

import (
	"errors"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/imroc/req"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type statisticNotify struct {
	gone.Goner
	logrus.Logger `gone:"gone-logger"`

	miniApp   service.IMiniApp `gone:"*"`
	user      service.IUser    `gone:"*"`
	notifyUrl string           `gone:"config,statistic.notify.url"`
	env       string           `gone:"config,server.env"`
}

//go:gone
func NewStatisticNotify() gone.Goner {
	return &statisticNotify{}
}

func (e statisticNotify) Consume(on emitter.OnEvent) {
	on(e.handleAppCreated)
}

func (e statisticNotify) handleAppCreated(evt *entity.AppCreatedEvent) error {
	e.Infof("handleAppCreated: appId=%s\n", evt.AppId)
	app, has, err := e.miniApp.GetAppDetailByUuid(evt.AppId)
	if err != nil {
		return err
	}

	if !has {
		return nil
	}

	return e.SendNotifyAppCreated(app)
}

func (e statisticNotify) SendNotifyAppCreated(app *entity.MiniAppDetailDto) error {
	e.Infof("SendNotifyAppCreated: appId=%s\n", app.Uuid)
	template := AppCreatedNotifyTemplate{
		AppId:           app.Uuid,
		Name:            app.Name,
		Description:     app.Description,
		CreatedNickname: app.CreatedBy.Nickname,
		Env:             e.env,
	}

	if len(app.DuplicateFrom) != 0 {
		duplicateFrom, has, err := e.miniApp.GetAppDetailByUuid(app.DuplicateFrom)
		if err != nil {
			return nil
		}

		if !has {
			return nil
		}

		template.DuplicateFrom = duplicateFrom.Name
	}

	resp, err := req.Post(e.notifyUrl, req.BodyJSON(map[string]any{
		"msg_type": "text",
		"content": map[string]any{
			"text": template.String(),
		},
	}))

	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	err = resp.ToJSON(&res)
	if err != nil {
		return err
	}

	if res.Code != 0 {
		e.Errorf("Send notify failed: %s", res.Msg)
		return errors.New(res.Msg)
	}

	return err
}
