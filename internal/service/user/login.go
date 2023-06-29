package user

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/redis"
	"gitlab.openviewtech.com/gone/gone-lib/token"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"gitlab.openviewtech.com/moyu-chat/ws-server/event"
	"net/http"
	"sort"
	"time"
)

func (s *svc) GenClientId() string {
	return token.GenShortUUID()
}

// SetCookie 设置cookie，在content为空字符串时移除cookie
func (s *svc) SetCookie(ctx *gin.Context, key, content string) {
	var maxAge int
	if content == "" {
		maxAge = -1
	} else {
		maxAge = int(s.CookieExpiresIn.Seconds())
	}
	ctx.SetCookie(key, content, maxAge, "/", s.mainDomain, s.CookieSecure, true)

	//清理子域名cookie
	ctx.SetCookie(key, content, -1, "/", "", s.CookieSecure, true)
}

func (s *svc) ParseJwt(token, secret string) (int64, error) {
	t, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, gin.NewParameterError("invalid token", http.StatusUnauthorized) //无效token
	}

	tokenInfo := t.Claims.(*TokenClaims)
	if tokenInfo.ExpiresAt < time.Now().UnixMilli() {
		return 0, gin.NewParameterError("expired token", http.StatusUnauthorized) //token过期
	}

	return tokenInfo.UserId, nil
}

func (s *svc) ParseJwtInfo(ctx *gin.Context) (userId int64, err gone.Error) {
	jwtToken := utils.CtxGetString(ctx, entity.JwtKey)
	if jwtToken == "" {
		return 0, gin.NewParameterError("not login", http.StatusUnauthorized) //无token
	}

	t, e := jwt.ParseWithClaims(jwtToken, &TokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(s.JwtSecret), nil
	})
	if e != nil {
		err = gin.NewParameterError("invalid token", http.StatusUnauthorized) //无效token
		return
	}
	tokenInfo := t.Claims.(*TokenClaims)
	if tokenInfo.ExpiresAt < time.Now().UnixMilli() {
		err = gin.NewParameterError("expired token", http.StatusUnauthorized) //token过期
		return
	}
	userId = tokenInfo.UserId
	csrfToken := ctx.GetHeader(entity.CsrfTokenKey)
	if csrfToken == "" {
		err = gin.NewParameterError("need csrf-token") //缺少csrf-token
		return
	}

	if csrfToken != tokenInfo.CsrfToken {
		err = gin.NewParameterError("expired token", http.StatusUnauthorized) //token过期
		return
	}

	exists, e := s.CheckUserExists(userId)
	if e != nil || !exists {
		err = gin.NewParameterError("user not found", http.StatusUnauthorized) //token过期
		return
	}

	clientId := utils.CtxGetString(ctx, entity.ClientIdKey)
	if !token.IsValidShortUUID(clientId) {
		err = gin.NewParameterError("invalid token", http.StatusUnauthorized) //无效token
		return
	}

	sessionKey := buildSessionKey(userId, clientId)
	//redis校验逻辑
	{
		var session Session
		e := s.Get(sessionKey, &session)
		if e != nil {
			if e == redis.ErrNil {
				err = gin.NewParameterError("expired token", http.StatusUnauthorized) //token过期
				return
			} else {
				//读取redis异常时，不踢用户下线
				s.Errorf("get session from redis err:%v", e)
			}
		} else {
			//redis中存储的不再是原来的jwt，jwt即过期
			if session.Jwt != jwtToken {
				err = gin.NewParameterError("expired token", http.StatusUnauthorized) //token过期
				return
			}
		}
	}

	//token续期 逻辑
	if tokenInfo.ExpiresAt < time.Now().Add(s.JwtRenewal).UnixMilli() {
		_, _ = s.Login(ctx, &entity.UserSimple{Id: userId})
	}
	return
}

func (s *svc) Login(ctx *gin.Context, user *entity.UserSimple) (rst *entity.LoginRst, err error) {
	userId := user.Id
	clientId := utils.CtxGetString(ctx, entity.ClientIdKey)
	key := buildSessionKey(userId, clientId)

	csrfToken := token.GenShortUUID()
	jwtToken, _ := s.genJwtToken(userId, csrfToken)
	session := Session{
		UserId:    userId,
		ClientId:  clientId,
		Jwt:       jwtToken,
		CreatedAt: time.Now(),
	}

	//信息保存到redis
	err = s.Put(key, session, s.JwtExpiresIn)
	if err != nil {
		s.Errorf("put session to redis err:%v", err)
	}

	//设置cookie
	s.SetCookie(ctx, s.CookieJwtKey, jwtToken)
	s.SetCookie(ctx, s.CookieClientIdKey, clientId)

	rst = new(entity.LoginRst)
	rst.User = user
	rst.CsrfToken = csrfToken

	//一个用户只允许一个设备登录
	s.tracer.Go(func() {
		s.kickMoreLogin(userId)
	})
	return
}

func (s *svc) Logout(ctx *gin.Context, userId int64, clientId string) error {
	key := buildSessionKey(userId, clientId)
	err := s.Remove(key)

	if err != nil {
		return err
	}

	s.tracer.Go(func() {
		s.KickUserById(userId)
	})

	//移除cookie
	s.SetCookie(ctx, s.CookieJwtKey, "")
	return nil
}

func (s *svc) GenQrToken() (*entity.QrTokenWarp, error) {
	expiredAt := time.Now().Add(s.QrTokenExpiresIn)

	return &entity.QrTokenWarp{
		QrToken: token.GenOfflineToken(expiredAt, s.LoginTokenSecret),
		Expired: expiredAt,
	}, nil
}

func (s *svc) GetQrTokenParams(qrToken string, redirectUrl string) (tokenExpired time.Time, authUrl string, err error) {
	if qrToken == "" {
		authUrl = fmt.Sprintf(
			"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&"+
				"scope=snsapi_userinfo&forcePopup=true&forceSnapShot=true#wechat_redirect",
			s.WxAppId,
			redirectUrl,
		)

		return
	}

	tokenExpired, err = token.DecodeOfflineToken(qrToken, s.LoginTokenSecret)
	if err != nil {
		err = gin.NewParameterError(err.Error())
	}
	authUrl = fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&"+
			"scope=snsapi_userinfo&forcePopup=true&state=%s&forceSnapShot=true#wechat_redirect",
		s.WxAppId,
		redirectUrl,
		qrToken,
	)
	return
}

func (s *svc) LoginMobile(code string, invitedBy *int64, fromApp, source string) (userSimple *entity.UserSimple, err error) {
	userInfo, err := s.parseWxCode(code)
	if err != nil {
		return nil, err
	}

	userSimple, err = s.findOrCreate(userInfo)
	if err != nil {
		return
	}

	loginToken := token.GenOfflineToken(time.Now().Add(s.QrTokenExpiresIn), s.LoginTokenSecret, userSimple.Id)
	userSimple, err = s.SetUserWithLoginToken(loginToken, userSimple.Nickname, userSimple.Avatar, invitedBy, fromApp, source)
	return
}

func (s *svc) AuthByCode(qrToken, code string) (userSimple *entity.UserSimple, err error) {
	_, e := token.DecodeOfflineToken(qrToken, s.LoginTokenSecret)
	if e != nil {
		err = gin.NewParameterError(e.Error())
	}

	user, err := s.parseWxCode(code)
	if err != nil {
		return nil, err
	}

	userSimple, err = s.findOrCreate(user)
	if err != nil {
		return
	}

	loginToken := token.GenOfflineToken(time.Now().Add(s.QrTokenExpiresIn), s.LoginTokenSecret, userSimple.Id)

	err = s.Send(&event.TriggerEvent{
		Scope:      event.TriggerEventScopeByScene,
		Scene:      event.LoginScene,
		SceneParam: qrToken,

		Type: event.TriggerTypeLoginSuc,
		Params: []interface{}{
			loginToken,
			userSimple.Id,
		},
	})
	return
}

func (s *svc) moyuPlanetCheck(userId int64) error {
	member, err := s.Planet.GetPlanetMemberByUserId(entity.MoyuPlanetId, userId)
	if err != nil {
		return err
	}

	if member == nil {
		err := s.Planet.MemberJoin(userId, entity.MoyuPlanetId)
		if err != nil {
			return err
		}
	} else if member.Status == entity.PlanetMemberStatusBlack {
		return gin.NewParameterError("用户被禁用")
	}
	return nil
}

func (s *svc) SetUserWithLoginToken(loginToken, nickName string, avatar entity.Url, invitedBy *int64, fromApp string, source string) (rst *entity.UserSimple, err error) {
	s.Infof("SetUserWithLoginToken:  fromApp=%s, source=%s\n", fromApp, source)
	var userId int64
	_, e := token.DecodeOfflineToken(loginToken, s.LoginTokenSecret, &userId)
	if e != nil {
		err = gin.NewParameterError(e.Error())
	}

	if userId == 0 {
		err = gin.NewParameterError("loginToken 无效")
		return
	}

	user, err := s.PUser.getUserById(userId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	err = s.moyuPlanetCheck(user.Id)
	if err != nil {
		return nil, err
	}

	s.Infof("SetUserWithLoginToken: user.fromApp=%s, user.source=%s\n", user.FromApp, user.Source)
	if nickName != "" {
		user.Nickname = nickName
	}

	if avatar != "" {
		user.Avatar = avatar
	}

	if user.InvitedBy == nil {
		if invitedBy != nil {
			tmpUser, err := s.GetUserById(*invitedBy)
			if err != nil {
				return nil, err
			}

			if tmpUser != nil {
				user.InvitedBy = invitedBy
			}
		} else {
			user.InvitedBy = new(int64)
		}
	}

	if user.FromApp == "" {
		if fromApp != "" {
			user.FromApp = fromApp
		}
	}

	if user.Source == "" {
		if source != "" {
			user.Source = source
		}
	}

	now := time.Now()
	user.LastLoginAt = user.LoginAt
	user.LoginAt = &now

	err = s.PUser.updateUserInfo(user)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if user.IsFirstLogin() {
		if err := s.Send(&entity.FirstLoginEvent{
			User: *user,
		}); err != nil {
			return nil, err
		}
	}

	return &entity.UserSimple{
		Id:       user.Id,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

func (s *svc) getJwtHashValue(userId int64, jwt string) (*Session, error) {
	var saved Session
	err := s.Get(buildSessionKey(userId, jwt), &saved)
	if err != nil {
		return nil, err
	}
	return &saved, nil
}

type sessionWrap struct {
	*Session
	k string
}

type sessionList []*sessionWrap

func (l sessionList) Len() int {
	return len(l)
}

func (l sessionList) Less(i, j int) bool {
	return l[i].CreatedAt.Before(l[j].CreatedAt)
}

func (l sessionList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

const kickUserDelaySecond = 2

// 延时${kickUserDelaySecond}秒踢出多点登录用户，保留最后一次登录
func (s *svc) kickMoreLogin(userId int64) {
	<-time.After(kickUserDelaySecond * time.Second)

	keys, err := s.Keys(buildSessionKey(userId, "*"))
	if err != nil {
		s.Errorf("redis: get session(%d) err:%v", userId, err)
		return
	}

	if len(keys) == 0 {
		return
	}

	list := make([]*sessionWrap, 0)
	for _, k := range keys {
		session := new(Session)
		err := s.Get(k, session)
		if err != nil {
			if err != redis.ErrNil {
				s.Errorf("redis: get session(k=%s) err:%v", k, err)
			}
		} else {
			if session.ClientId != s.TestClientId {
				list = append(list, &sessionWrap{Session: session, k: k})
			}
		}
	}
	sort.Sort(sessionList(list))

	list = list[0 : len(list)-1]
	for _, session := range list {
		s.Debugf("kick user(id=%s, clientId=%s)", session.UserId, session.ClientId)
		err := s.Remove(session.k)
		if err != nil {
			s.Errorf("redis: remove %s err:%v", session.k, err)
		}
		s.sendClientKickTrigger(session.ClientId)
	}
}

func (s *svc) sendClientKickTrigger(clientId string) {
	if clientId == "" {
		return
	}

	err := s.Send(&event.TriggerEvent{
		Scope:      event.TriggerEventScopeByScene,
		Scene:      event.ClientScene,
		SceneParam: clientId,

		Type: event.TriggerTypeKick,
		Params: []interface{}{
			"user login in other device",
		},
	})

	if err != nil {
		s.Errorf("send TriggerEvent(kick, clientId=%s) err:%v", clientId, err)
	}
}

func (s *svc) KickUserById(userId int64, msg ...string) {
	if userId == 0 {
		return
	}

	keys, err := s.Keys(buildSessionKey(userId, "*"))
	if err != nil {
		s.Errorf("redis: get session(%d) err:%v", userId, err)
		return
	}

	for _, k := range keys {
		err = s.Remove(k)
		if err != nil {
			s.Errorf("redis: remove %s err:%v", k, err)
		}
	}

	if len(keys) == 0 {
		return
	}

	err = s.Send(&event.TriggerEvent{
		Scope:  event.TriggerEventScopeByUser,
		UserId: userId,

		Type:   event.TriggerTypeKick,
		Params: utils.ToInterfaceSlice(msg...),
	})

	if err != nil {
		s.Errorf("send TriggerEvent(kick, userId=%s) err:%v", userId, err)
	}
}

// genJwtToken 生成jwtToken
// jwtToken
// expiredAt 过期时间，毫秒时间戳
func (s *svc) genJwtToken(userId int64, csrfToken string) (jwtToken string, expiredAt int64) {
	expiredAt = time.Now().Add(s.JwtExpiresIn).UnixMilli()
	claims := &TokenClaims{
		UserId:    userId,
		CsrfToken: csrfToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken, err := t.SignedString([]byte(s.JwtSecret))
	if err != nil {
		panic(err)
	}
	return
}

type Session struct {
	UserId    int64     `json:"u"`
	ClientId  string    `json:"c"`
	Jwt       string    `json:"j"`
	CreatedAt time.Time `json:"at"`
}

type TokenType int
