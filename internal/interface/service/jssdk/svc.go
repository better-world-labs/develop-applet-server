package jssdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"strconv"
	"time"
)

type JSTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

type svc struct {
	gone.Goner
	logrus.Logger `gone:"gone-logger"`

	WxAppId  string       `gone:"config,user.login.wechat.appid"`
	WxSecret string       `gone:"config,user.login.wechat.secret"`
	Cache    *TicketCache `gone:"*"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s svc) GetSignature(url string) (entity.JSSDKSignature, error) {
	jsTicket, err := s.getTicket()
	if err != nil {
		return entity.JSSDKSignature{}, err
	}

	signature := entity.JSSDKSignature{
		NonceStr:  uuid.NewString(),
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		AppId:     s.WxAppId,
	}
	signature.Signature = createSignature(signature.NonceStr, jsTicket, signature.Timestamp, url)
	return signature, nil
}

func (s *svc) getTicket() (string, error) {
	ticket, has, err := s.Cache.get()
	if err != nil {
		return "", err
	}

	if !has {
		s.Infof("ticket from cache expired, getting from wechat...\n")
		if ticket, err = s.getJsTicketAndPutCache(); err != nil {
			return "", err
		}
	}

	return ticket, nil
}

func (s *svc) getJsTicketAndPutCache() (string, error) {
	token, err := s.getAccessToken()
	if err != nil {
		return "", err
	}

	ticket, err := s.getJSTicketFromWechat(token)
	if err != nil {
		return "", err
	}

	if err := s.Cache.put(ticket); err != nil {
		return "", err
	}

	return ticket.Ticket, nil
}

func (s *svc) getAccessToken() (string, error) {
	params := req.QueryParam{
		"appid":      s.WxAppId,
		"secret":     s.WxSecret,
		"grant_type": "client_credential",
	}

	resp, err := req.Get("https://api.weixin.qq.com/cgi-bin/token", params)
	if err != nil {
		return "", err
	}

	var res struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
	}

	body := resp.String()
	s.Logger.Infof("getAccessToken: response=%s\n", body)
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return "", err
	}

	if res.ErrCode != 0 {
		s.Logger.Errorf("getAccessToken error: %s\n", res.ErrMsg)
		return "", errors.New(res.ErrMsg)
	}

	return res.AccessToken, nil
}

func (s *svc) getJSTicketFromWechat(accessToken string) (JSTicket, error) {
	resp, err := req.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken))
	if err != nil {
		return JSTicket{}, err
	}

	var res struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		JSTicket
	}

	body := resp.String()
	s.Logger.Infof("getJSTicket: response=%s\n", body)
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return JSTicket{}, err
	}

	if res.ErrCode != 0 {
		s.Logger.Errorf("getJSTicket error: %s\n", res.ErrMsg)
		return JSTicket{}, errors.New(res.ErrMsg)
	}

	return res.JSTicket, nil
}
