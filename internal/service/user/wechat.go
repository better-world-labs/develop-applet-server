package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req"
	businesserrors "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/error"
)

type accessToken struct {
	Openid         string `json:"openid"`
	AccessToken    string `json:"access_token"`
	IsSnapshotUser int    `json:"is_snapshotuser"`
}

type wechatRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (w wechatRes) isSucceed() bool {
	return w.ErrCode == 0
}

type wxUserInfo struct {
	Openid       string `json:"openid"`
	Nickname     string `json:"nickname"`
	HeadImageUrl string `json:"headimgurl"`
	Unionid      string `json:"unionid"`
}

func (u accessToken) isSnapshotUser() bool {
	return u.IsSnapshotUser == 1
}

func (s *svc) parseWxCode(code string) (*wxUserInfo, error) {
	token, err := s.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	if token.isSnapshotUser() {
		return nil, businesserrors.ErrorInvalidLoginUser
	}

	user, err := s.getUserInfo(token.AccessToken, token.Openid)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *svc) getAccessToken(code string) (*accessToken, error) {
	s.Logger.Debugf("getAccessToken: code=%s\n", code)

	params := req.QueryParam{
		"appid":      s.WxAppId,
		"secret":     s.WxSecret,
		"grant_type": "authorization_code",
		"code":       code,
	}

	resp, err := req.Get("https://api.weixin.qq.com/sns/oauth2/access_token", params)
	if err != nil {
		return nil, err
	}

	var res struct {
		wechatRes
		accessToken
	}

	body := resp.String()
	s.Logger.Infof("getAccessToken: response=%s\n", body)
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return nil, err
	}

	if !res.isSucceed() {
		s.Logger.Errorf("getAccessToken error: %s\n", res.ErrMsg)
		return nil, errors.New(res.ErrMsg)
	}

	//mark 可以通过 接口（https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId）
	// 换取 unionID
	return &res.accessToken, nil
}

func (s *svc) getUserInfo(accessToken, openid string) (*wxUserInfo, error) {
	s.Logger.Debugf("getUserInfo: accessToken=%s, openid=%s\n", accessToken, openid)

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", accessToken, openid)
	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	var res struct {
		wechatRes
		wxUserInfo
	}

	body := resp.String()
	s.Logger.Infof("getAccessToken: response=%s\n", body)
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return nil, err
	}

	if !res.isSucceed() {
		s.Logger.Errorf("getAccessToken error: %s\n", res.ErrMsg)
		return nil, errors.New(res.ErrMsg)
	}

	//mark 可以通过 接口（https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId）
	// 换取 unionID
	return &res.wxUserInfo, nil
}
