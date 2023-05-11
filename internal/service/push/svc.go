package push

import (
	"errors"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/imroc/req"
	"strconv"
)

type PushSvc struct {
	gone.Goner
	logrus.Logger `gone:"gone-logger"`

	pushUrl string `gone:"config,push.url"`
}

//go:gone
func NewSvc() gone.Goner {
	return &PushSvc{}
}

func (s PushSvc) PushMessage(userId int64, payload any) error {
	var param = struct {
		Topic   string `json:"topic"`
		Message any    `json:"message"`
	}{
		Topic:   strconv.FormatInt(userId, 10),
		Message: payload,
	}

	s.Infof("PushMessage: %v\n", param)
	resp, err := req.Post(s.pushUrl, req.BodyJSON(param))
	if err != nil {
		return err
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := resp.ToJSON(&response); err != nil {
		return err
	}

	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	return nil
}
