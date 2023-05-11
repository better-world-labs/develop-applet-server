package user

import (
	"encoding/json"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/dgrijalva/jwt-go"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/redis"
	"github.com/gone-io/gone/goner/tracer"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultPlanetToJoin       = 1
	DefaultOffTime            = "18:00:00"
	DefaultAppearanceTheme    = entity.AppearanceThemeBright
	DefaultBossKey            = "ALT + Z"
	DefaultSiteSettings       = "{\"type\": \"default\", \"customIcon\": \"\", \"customTitle\": \"\"}"
	DefaultMonthlySalary      = 10000
	DefaultMonthlyWorkingDays = 22

	TodayBrowseTimeRedisZSetKey = "moyu#today_browse_duration#zset"
	TodayAccessTimeRedisZSetKey = "moyu#today_access_duration#zset"
	ZMember                     = "mem#%d"
)

var CnLocation, _ = time.LoadLocation("Asia/Shanghai")

//go:gone
func NewUserService() gone.Goner {
	return &svc{}
}

type svc struct {
	gone.Flag
	T      xorm.Engine   `gone:"gone-xorm"`
	tracer tracer.Tracer `gone:"gone-tracer"`

	PUser iUserPersistence `gone:"*"`

	PUserSettings     iUserSettingsPersistence `gone:"*"`
	PUserBrowseRecord iUserBrowsePersistence   `gone:"*"`
	PUserAccessRecord iUserAccessPersistence   `gone:"*"`
	MessageRecord     service.IMessageRecord   `gone:"*"`
	RedisService      service.IRedisService    `gone:"*"`

	Points   service.IPointStrategy     `gone:"*"`
	Identity service.IAnonymousIdentity `gone:"*"`

	Planet service.IPlanet `gone:"*"`

	JwtSecret    string        `gone:"config,jwt.secret"`
	JwtExpiresIn time.Duration `gone:"config,jwt.expiresIn"`
	JwtRenewal   time.Duration `json:"jwt.renewal"`

	CookieJwtKey      string        `gone:"config,cookie.jwt-key"`
	CookieClientIdKey string        `gone:"config,cookie.client-id-key"`
	CookieExpiresIn   time.Duration `gone:"config,cookie.expiresIn"`
	CookieSecure      bool          `gone:"config,cookie.secure"`
	TestClientId      string        `gone:"config,cookie.client-id.for-test"`

	QrTokenExpiresIn    time.Duration `gone:"config,user.login.qrToken.expiresIn"`
	LoginTokenExpiresIn time.Duration `gone:"config,user.login.token.expiresIn"`
	LoginTokenSecret    string        `gone:"config,user.login.token.secret"`
	WxAppId             string        `gone:"config,user.login.wechat.appid"`
	WxSecret            string        `gone:"config,user.login.wechat.secret"`

	redis.Cache    `gone:"gone-redis-cache"`
	logrus.Logger  `gone:"gone-logger"`
	emitter.Sender `gone:"gone-emitter"`

	mainDomain string `gone:"config,server.domain"`
}

func (s *svc) UpdateUser(userId int64, req *entity.UpdateUserReq) (*entity.User, error) {
	user, err := s.PUser.getUserById(userId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	err = s.PUser.updateUserInfo(user)
	if err != nil {
		return nil, gin.ToError(err)
	}
	return user, nil
}

func (s *svc) GetUser(userId int64) (*entity.User, error) {
	return s.PUser.getUserById(userId)
}

func (s *svc) GetUserByOpenId(openId string) (*entity.UserSimple, error) {
	user, err := s.PUser.getUserByOpenId(openId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if user == nil {
		return nil, gin.NewParameterError("user does not exist")
	}

	return &entity.UserSimple{
		Id:       user.Id,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId    int64  `json:"u"`
	CsrfToken string `json:"x"`
}

func (s *svc) findOrCreate(wx *wxUserInfo) (*entity.UserSimple, error) {
	user, err := s.PUser.getUserByOpenId(wx.Openid)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if user == nil {
		u := userDo{
			User: &entity.User{
				WxOpenId:  wx.Openid,
				Nickname:  wx.Nickname,
				Avatar:    entity.Url(wx.HeadImageUrl),
				CreatedAt: time.Now(),
			},
		}

		err = s.T.Transaction(func(session xorm.Interface) error {
			err = s.PUser.createUser(&u)
			if err != nil {
				return gin.ToError(err)
			}
			user = u.User

			err = s.PUserSettings.createUserSettings(entity.UserSettings{
				UserId:             u.Id,
				EndOffTime:         DefaultOffTime,
				BossKey:            DefaultBossKey,
				AppearanceTheme:    DefaultAppearanceTheme,
				SiteSettings:       DefaultSiteSettings,
				MonthlySalary:      DefaultMonthlySalary,
				MonthlyWorkingDays: DefaultMonthlyWorkingDays,
			})
			if err != nil {
				return gin.ToError(err)
			}

			err = s.Planet.MemberJoin(u.Id, DefaultPlanetToJoin)
			if err != nil {
				return gin.ToError(err)
			}

			_, err := s.Points.ApplyPoints(u.Id, entity.StrategyArgNewRegister{})
			return err
		})

		if err != nil {
			return nil, err
		}
	}
	return &entity.UserSimple{Id: user.Id, Nickname: user.Nickname, Avatar: user.Avatar}, nil
}

func (s *svc) PostEvent(event *entity.UserEvent) error {
	switch event.Type {
	case entity.EventTypeViewApp:
		return s.sendAppViewEvent(event)

	default:
		return gin.NewParameterError("invalid event type")
	}
}

func (s *svc) sendAppViewEvent(event *entity.UserEvent) error {
	if len(event.Args) != 1 {
		return gin.NewParameterError("invalid param")
	}

	appId, ok := event.Args[0].(string)
	if !ok {
		return gin.NewParameterError("invalid appId")
	}

	return s.Send(&entity.AppViewEvent{
		AppId: appId,
		User:  event.UserId,
	})
}
func (s *svc) GetUserSettings(userId int64) (*domain.UserSettings, error) {
	settings, err := s.PUserSettings.getUserSettingsByUserId(userId)
	if err != nil {
		return nil, err
	}

	defaultSiteSettings, _ := json.Marshal(entity.SiteSettings{
		Type: entity.SiteSettingsDefault,
	})

	if settings == nil || settings.Id == 0 {
		settings = &entity.UserSettings{
			UserId:             userId,
			EndOffTime:         DefaultOffTime,
			AppearanceTheme:    DefaultAppearanceTheme,
			SiteSettings:       string(defaultSiteSettings),
			MonthlySalary:      DefaultMonthlySalary,
			MonthlyWorkingDays: DefaultMonthlyWorkingDays,
		}

		err = s.PUserSettings.createUserSettings(*settings)
		if err != nil {
			return nil, err
		}
	}

	if len(settings.SiteSettings) == 0 {
		settings.SiteSettings = string(defaultSiteSettings)
	}

	var domainSetting entity.SiteSettings
	err = json.Unmarshal([]byte(settings.SiteSettings), &domainSetting)
	if err != nil {
		return nil, gone.NewInnerError(500, fmt.Sprintf("translate settings to domain setting failed, err: %s", err.Error()))
	}

	if len(domainSetting.Type) == 0 {
		domainSetting.Type = entity.SiteSettingsDefault
	}

	res := &domain.UserSettings{
		BossKey:            settings.BossKey,
		EndOffTime:         settings.EndOffTime,
		AppearanceTheme:    settings.AppearanceTheme,
		SiteSettings:       domainSetting,
		MonthlySalary:      settings.MonthlySalary,
		MonthlyWorkingDays: settings.MonthlyWorkingDays,
	}

	return res, nil
}

func (s *svc) CheckUserExists(id int64) (bool, error) {
	return s.PUser.checkExists(id)
}

func (s *svc) GetUserById(id int64) (*entity.User, error) {
	return s.PUser.getUserById(id)
}

func (s *svc) GetUserSimpleInBatch(ids []int64) ([]*entity.UserSimple, error) {
	var res []*entity.UserSimple

	users, err := s.GetUserInBatch(ids)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		res = append(res, &entity.UserSimple{
			Id:       u.Id,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
		})
	}

	return res, nil
}

func (s *svc) GetUserInBatch(ids []int64) ([]*entity.User, error) {

	batchSize := 200
	batch := len(ids) / batchSize
	extra := len(ids)%batchSize != 0
	if extra {
		batch += 1
	}
	res := make([]*entity.User, 0)
	for i := 0; i < batch; i++ {
		var subIds []int64
		if i == batch-1 {
			subIds = ids[i*batchSize:]
		} else {
			subIds = ids[i*batchSize : (i+1)*batchSize]
		}

		subUsers, err := s.PUser.getUserInBatch(subIds)
		if err != nil {
			return nil, err
		}
		res = append(res, subUsers...)
	}

	return res, nil
}

func (s *svc) UpdateUserInfo(id int64, nickname, avatar string) error {
	u, err := s.PUser.getUserById(id)
	if err != nil {
		return err
	}

	u.Nickname = nickname
	u.Avatar = entity.Url(avatar)

	return s.PUser.updateUserInfo(u)
}

func (s *svc) UpdateWorkOffTime(userId int64, offTime string) error {
	settings, err := s.PUserSettings.getUserSettingsByUserId(userId)
	if err != nil {
		return err
	}

	if settings == nil || settings.Id == 0 {
		return s.PUserSettings.createUserSettings(entity.UserSettings{
			UserId:             userId,
			AppearanceTheme:    DefaultAppearanceTheme,
			EndOffTime:         offTime,
			MonthlySalary:      DefaultMonthlySalary,
			MonthlyWorkingDays: DefaultMonthlyWorkingDays,
		})
	}

	return s.PUserSettings.updateWorkOffTime(userId, offTime)
}

func (s *svc) UpdateBossKey(userId int64, bossKey string) error {
	settings, err := s.PUserSettings.getUserSettingsByUserId(userId)
	if err != nil {
		return err
	}

	if settings == nil || settings.Id == 0 {
		return s.PUserSettings.createUserSettings(entity.UserSettings{
			UserId:             userId,
			EndOffTime:         DefaultOffTime,
			AppearanceTheme:    DefaultAppearanceTheme,
			BossKey:            bossKey,
			MonthlySalary:      DefaultMonthlySalary,
			MonthlyWorkingDays: DefaultMonthlyWorkingDays,
		})
	}
	return s.PUserSettings.updateBossKey(userId, bossKey)
}

func (s *svc) UpdateOnline(userId int64, online bool, t time.Time) error {
	return s.PUser.updateOnline(userId, online, t)
}

func (s *svc) UpdateUserSettings(userSetting entity.UserSettings) error {
	var updates []entity.ComponentNameEnum
	if len(userSetting.AppearanceTheme) != 0 {
		updates = append(updates, entity.AppearanceThemeComponent)
	}
	if len(userSetting.SiteSettings) != 0 {
		updates = append(updates, entity.SiteSettingsComponent)
	}
	return s.PUserSettings.updateUserSettings(userSetting, updates)
}

func (s *svc) CalculateEarlierThan(offTime string) (*domain.EarlierThan, error) {
	offTimes, err := s.PUserSettings.getAllUserSettingsOrderByOffTime()
	if err != nil {
		return nil, err
	}

	pos := findStrFirstPos(offTime, offTimes)
	length := int64(len(offTimes))

	var percentFloat float64
	if pos == 0 { // first
		percentFloat = 100.00
	} else if pos == int64(len(offTimes)-1) { // last
		percentFloat = 0.00
	} else {
		percent := float64(length-pos) * 100 / float64(length)
		percentFloat, err = strconv.ParseFloat(fmt.Sprintf("%.2f", percent), 64)
		if err != nil {
			return nil, err
		}
	}

	res := &domain.EarlierThan{
		EarlierThan: percentFloat,
	}

	return res, nil
}

func (s *svc) GetSimpleUserSettings(userId int64, componentNames []entity.ComponentNameEnum) ([]*entity.ComponentElement, error) {
	var res []*entity.ComponentElement

	var validComponentNames []entity.ComponentNameEnum
	linq.From(componentNames).WhereT(func(component entity.ComponentNameEnum) bool {
		valid := component.ValidComponent()
		if !valid {
			s.Logger.Infof("It's a invalid component name of [%s] ", component)
		}
		return valid
	}).ToSlice(&validComponentNames)

	settings, err := s.PUserSettings.getSimpleUserSettingsByUserId(userId, validComponentNames)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return nil, nil
	}

	linq.From(validComponentNames).ForEachT(func(component entity.ComponentNameEnum) {
		var element entity.ComponentElement
		switch component {
		case entity.AppearanceThemeComponent:
			element.ComponentName = entity.AppearanceThemeComponent
			element.ComponentSettings = string(settings.AppearanceTheme)
		case entity.SiteSettingsComponent:
			var setting entity.SiteSettings
			err = json.Unmarshal([]byte(settings.SiteSettings), &setting)
			if err != nil {
				s.Logger.Errorf("translate to site-setting failed, err: %s", err.Error())
			}
			element.ComponentName = entity.SiteSettingsComponent
			element.ComponentSettings = setting
		}

		res = append(res, &element)
	})

	return res, nil
}

func (s *svc) ListOnlineUsers() ([]*entity.User, error) {
	return s.PUser.getOnlineUsers()
}

func (s *svc) UpdateConnectTime(userId int64) error {
	err := s.PUser.updateConnectTime(userId)
	if err != nil {
		return err
	}

	return s.PUserAccessRecord.insertWithoutEnd(userId)
}

func (s *svc) AccumulateTotalAccessDuration(userId int64) error {
	user, err := s.PUser.getUserById(userId)
	if err != nil {
		s.Logger.Errorf("get user info failed, user id: [%d], err: %s", userId, err.Error())
		return err
	}

	if user == nil {
		return nil
	}

	currentTime := time.Now()
	lastConnectTime := user.ConnectTime
	if lastConnectTime == nil {
		s.Logger.Infof("last connect time is nil, use current time. ")
		lastConnectTime = &currentTime
	}

	if currentTime.Before(*lastConnectTime) {
		s.Logger.Errorf("current time is before user's connect time. cannot accumulate it.")
		return nil
	}

	duration := currentTime.Unix() - (*lastConnectTime).Unix()
	err = s.PUser.updateTotalAccessDuration(userId, duration+user.TotalAccessDuration)
	if err != nil {
		return err
	}

	// incr to redis
	todayDuration := duration
	beginningDate := genBeginningDate(currentTime, CnLocation)
	if lastConnectTime.Before(beginningDate) {
		s.Logger.Infof("last connect time of user [%d] is yesterday, get time from 00:00. ", userId)
		todayDuration = currentTime.Unix() - beginningDate.Unix()
	}

	memKey := fmt.Sprintf(ZMember, userId)
	expireTime := genTtlSec(currentTime, CnLocation)
	s.Logger.Infof("increase access duration for user [%d], today duration is [%d]. ", userId, todayDuration)
	err = s.RedisService.ZIncrbyWitEx(TodayAccessTimeRedisZSetKey, memKey, todayDuration, expireTime)
	if err != nil && err != redis.ErrNil {
		return err
	}

	return s.PUserAccessRecord.updateEndTimeNull(userId)
}

func (s *svc) AccumulateTotalBrowseDuration(userId int64, timeQuantum []entity.TimeQuantum) error {
	if len(timeQuantum) == 0 {
		return nil
	}

	browseDate := genBrowseDate(time.Now())
	lastBrowseRecord, err := s.PUserBrowseRecord.getLastRecordByUserId(userId, browseDate)
	if err != nil {
		return err
	}

	if lastBrowseRecord == nil {
		err = s.PUserBrowseRecord.insert(entity.UserBrowserRecord{
			UserId:         userId,
			BrowseDate:     browseDate,
			BrowseDuration: 0,
			LastBrowseTime: time.Now(),
		})
		if err != nil {
			return err
		}
		lastBrowseRecord, err = s.PUserBrowseRecord.getLastRecordByUserId(userId, browseDate)
	}

	curDuration := int64(len(timeQuantum)) * 10

	// 1. accumulate browse duration to redis
	memKey := fmt.Sprintf(ZMember, userId)
	expireTime := genTtlSec(time.Now(), CnLocation)

	err = s.RedisService.ZIncrbyWitEx(TodayBrowseTimeRedisZSetKey, memKey, curDuration, expireTime)
	if err != nil && err != redis.ErrNil {
		return err
	}

	// 2. accumulate to mysql
	user, err := s.PUser.getUserById(userId)
	if err != nil {
		return err
	}
	err = s.PUser.updateTotalBrowseDuration(userId, user.TotalBrowseDuration+curDuration)
	if err != nil {
		return err
	}

	// 3. update daily record
	s.Logger.Infof("accumulate duration [%d] for user [%d]. ", curDuration, userId)
	lastBrowseRecord.BrowseDuration += curDuration
	lastBrowseRecord.LastBrowseTime = timeQuantum[len(timeQuantum)-1].EndTime

	return s.PUserBrowseRecord.update(*lastBrowseRecord)
}

func (s *svc) GenMoyuDetail(userId int64) (*domain.MoyuDetail, error) {
	user, err := s.PUser.getUserById(userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, gone.NewInnerError(400, fmt.Sprintf("user [%s] not found.", userId))
	}

	// 1. get browse duration from redis
	memKey := fmt.Sprintf(ZMember, userId)
	todayBrowseDuration, _ := s.RedisService.ZScore(TodayBrowseTimeRedisZSetKey, memKey)
	//todayBrowseDuration, _ := s.RedisService.ZScore(TodayAccessTimeRedisZSetKey, memKey)

	// 2. get more than percent
	//durations, err := s.PUser.getUserOrderDesc("total_access_duration")
	durations, err := s.PUser.getUserOrderDesc("total_browse_duration")
	if err != nil {
		s.Logger.Errorf("Get all users' browse duration from db failed, err: %s", err.Error())
		return nil, err
	}
	linq.From(durations).OrderByDescendingT(func(aInt int64) int64 {
		return aInt
	}).ToSlice(&durations)
	length := int64(len(durations))
	pos := findIntFirstPos(user.TotalBrowseDuration, durations)
	//pos := findIntFirstPos(user.TotalAccessDuration, durations)

	var percentFloat float64
	if pos == 0 { // first
		percentFloat = 100.00
	} else if pos == int64(length-1) { // last
		percentFloat = 0.00
	} else {
		percent := float64(length-pos) * 100 / float64(length)
		percentFloat, err = strconv.ParseFloat(fmt.Sprintf("%.1f", percent), 64)
		if err != nil {
			return nil, err
		}
	}

	// 3. find last report time
	browseDate := genBrowseDate(time.Now())
	browserRecord, err := s.PUserBrowseRecord.getLastRecordByUserId(userId, browseDate)
	if err != nil {
		s.Logger.Errorf("Get last report time from db failed, err: %s", err.Error())
		return nil, err
	}

	var lastTime time.Time
	if browserRecord != nil {
		lastTime = browserRecord.LastBrowseTime
	}

	// 4. count msg num
	cntByUserId, err := s.MessageRecord.GetMessageCntByUserId(userId)
	if err != nil {
		s.Logger.Errorf("Get msg count from db failed, err: %s", err.Error())
		return nil, err
	}

	// 5. calculate second salary
	settings, err := s.PUserSettings.getUserSettingsByUserId(userId)
	if err != nil {
		s.Logger.Errorf("Get user [%d] settings from db failed, err: %s", userId, err.Error())
		return nil, err
	}

	if settings == nil {
		settings = &entity.UserSettings{
			MonthlySalary:      10000,
			MonthlyWorkingDays: 22,
		}
	}
	secondSalary := float64(settings.MonthlySalary) / float64(settings.MonthlyWorkingDays*9*60*60)

	// 6. update access data
	var exDuration int64
	var todayDuration int64
	if user.ConnectTime != nil && user.Online {
		s.Logger.Infof("user [%d] last connect at [%v]. ", userId, user.ConnectTime)
		exDuration = time.Now().Unix() - user.ConnectTime.Unix()
		connectTime := user.ConnectTime
		beginningTime := genBeginningDate(time.Now(), CnLocation)
		if user.ConnectTime.Before(beginningTime) {
			connectTime = &beginningTime
		}
		todayDuration = time.Now().Unix() - connectTime.Unix()
	}

	s.Logger.Infof("increase access duration for user [%d], total duration is [%d], today exDuration is [%d]. ",
		userId, exDuration, todayDuration)

	err = s.PUser.updateConnectTime(userId)
	err = s.PUser.updateTotalAccessDuration(userId, user.TotalAccessDuration+exDuration)

	expireTime := genTtlSec(time.Now(), CnLocation)
	err = s.RedisService.ZIncrbyWitEx(TodayAccessTimeRedisZSetKey, memKey, todayDuration, expireTime)

	res := domain.MoyuDetail{
		JoinDate:             user.CreatedAt,
		TodayBrowseDuration:  todayBrowseDuration + exDuration,
		TotalBrowseDuration:  user.TotalBrowseDuration,
		MoreThan:             percentFloat,
		SecondSalary:         secondSalary,
		LastReportBrowseTime: lastTime,
		AccumulateMsgCnt:     cntByUserId,
		User: domain.UserBaseInfo{
			Id:       user.Id,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
	}

	return &res, nil
}

func (s *svc) UpdateWorkInfo(userId int64, info domain.WorkSettings) error {
	var updates []entity.ComponentNameEnum
	if len(info.OffWorkTime) != 0 {
		updates = append(updates, entity.OffWorkTime)
	}
	if info.MonthlySalary != 0 {
		updates = append(updates, entity.MonthlySalary)
	}
	if info.MonthlyWorkingDays != 0 {
		updates = append(updates, entity.MonthlyWorkingDays)
	}

	userSetting := entity.UserSettings{
		UserId:             userId,
		EndOffTime:         info.OffWorkTime,
		MonthlySalary:      info.MonthlySalary,
		MonthlyWorkingDays: info.MonthlyWorkingDays,
	}

	return s.PUserSettings.updateUserSettings(userSetting, updates)
}

func (s *svc) GetBrowseRankingList(top int, period entity.StatPeriod) ([]*domain.RankInfo, error) {
	users, err := s.ListOnlineUsers()
	if err != nil {
		return nil, err
	}

	err = s.updateUserAccessTime(TodayAccessTimeRedisZSetKey, users)
	if err != nil {
		return nil, err
	}

	switch period {
	case entity.Daily:
		return s.getTodayBrowseRankingList(int64(top))
	default:
		return s.getTotalBrowseRankingList(top)
	}
}

func (s *svc) getTodayBrowseRankingList(top int64) ([]*domain.RankInfo, error) {
	var res []*domain.RankInfo
	memIds, err := s.RedisService.ZRevrangeByScore(TodayBrowseTimeRedisZSetKey, 0, top, true)
	//memIds, err := s.RedisService.ZRevrangeByScore(TodayAccessTimeRedisZSetKey, 0, top, true)
	if err != nil {
		s.Logger.Errorf("Get today ranking list from redis failed, err: %s", err.Error())
		return nil, err
	}

	if len(memIds)%2 != 0 {
		return nil, gone.NewInnerError(500, "result from zrange is not one pair.")
	}

	zPair := make(map[int64]int64)
	var userIds []int64
	for i := 0; i < len(memIds); i += 2 {
		splits := strings.Split(memIds[i], "#")
		if len(splits) != 2 {
			continue
		}

		userId, err := strconv.ParseInt(splits[1], 10, 64)
		if err != nil {
			s.Logger.Warnf("cannot from [%s] parse userId, err: %s", memIds[i], err.Error())
			continue
		}

		duration, err := strconv.ParseInt(memIds[i+1], 10, 64)
		if err != nil {
			s.Logger.Warnf("cannot from [%s] parse duration, err: %s", memIds[i+1], err.Error())
			continue
		}

		if duration == 0 {
			s.Logger.Warnf("user [%d] duration is 0, skip it.", userId)
			continue
		}

		zPair[userId] = duration
		userIds = append(userIds, userId)
	}

	users, err := s.PUser.getUserInBatch(userIds)
	if err != nil {
		s.Logger.Errorf("Get users failed, user ids [%v], err: %s", userIds, err.Error())
		return nil, err
	}

	linq.From(users).ForEachT(func(oneUser *entity.User) {
		res = append(res, &domain.RankInfo{User: domain.UserBaseInfo{
			Id:       oneUser.Id,
			Nickname: oneUser.Nickname,
			Avatar:   oneUser.Avatar,
		},
			BrowseDuration: zPair[oneUser.Id],
		})
	})

	linq.From(res).WhereT(func(oneUser *domain.RankInfo) bool {
		return oneUser.BrowseDuration > 0
	}).OrderByDescendingT(func(oneEle *domain.RankInfo) int64 {
		return oneEle.BrowseDuration
	}).ToSlice(&res)

	var rank int64
	linq.From(res).ForEachT(func(oneEle *domain.RankInfo) {
		rank += 1
		oneEle.Rank = rank
	})

	return res, nil
}

func (s *svc) getTotalBrowseRankingList(top int) ([]*domain.RankInfo, error) {
	var res []*domain.RankInfo
	rankingList, err := s.PUser.getTotalBrowseDurationRankingList(top)
	if err != nil {
		s.Logger.Errorf("Get total ranking list from db failed, err: %s", err.Error())
		return nil, err
	}

	var rank int64
	linq.From(rankingList).WhereT(func(user *entity.User) bool {
		return user.TotalBrowseDuration > 0
	}).OrderByDescendingT(func(user *entity.User) int64 {
		return user.TotalBrowseDuration
	}).ForEachT(func(user *entity.User) {
		rank += 1
		res = append(res, &domain.RankInfo{
			Rank: rank,
			User: domain.UserBaseInfo{
				Id:       user.Id,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
			},
			BrowseDuration: user.TotalBrowseDuration,
		})
	})

	return res, nil
}

func (s *svc) updateUserAccessTime(mSetKey string, users []*entity.User) error {
	var err error
	currentTime := time.Now()
	beginTime := genBeginningDate(currentTime, CnLocation)
	ttl := genTtlSec(currentTime, CnLocation)
	for _, user := range users {
		lastConnectTime := user.ConnectTime
		if lastConnectTime == nil {
			lastConnectTime = &currentTime
		}

		// 只统计今日的
		if lastConnectTime.Before(beginTime) {
			s.Logger.Infof("last connect time of user [%d] is yesterday, get time from 00:00. ", user.Id)
			lastConnectTime = &beginTime
		}

		exDuration := currentTime.Unix() - lastConnectTime.Unix()
		err = s.RedisService.ZIncrbyWitEx(mSetKey, fmt.Sprintf(ZMember, user.Id), exDuration, ttl)
		if err != nil {
			s.Logger.Errorf("increase user [%d] today access duration failed, err： %s", user.Id, err.Error())
		}
	}

	var userIds []int64
	linq.From(users).SelectT(func(user *entity.User) int64 {
		return user.Id
	}).ToSlice(&userIds)

	go func() {
		inErr := s.PUser.updateAccessTimeByConnectTime(userIds)
		if inErr != nil {
			s.Logger.Errorf("update access time by connect time failed. err: %s", err.Error())
		}
	}()

	return err
}

func findStrFirstPos(src string, strArr []string) int64 {
	var l, r int64
	l = 0
	r = int64(len(strArr) - 1)
	for l+1 < r {
		mid := l + (r-l)/2
		if strArr[mid] == src {
			return mid
		} else if strings.Compare(src, strArr[mid]) == -1 {
			r = mid
		} else {
			l = mid
		}
	}

	// 未匹配上，且在所有元素之前
	if src < strArr[l] {
		return 0
	}

	// 为匹配上，且在所有元素之后
	if src > strArr[r] {
		return int64(len(strArr) - 1)
	}

	// 在l和r之间
	return l + 1
}

func findIntFirstPos(src int64, intArr []int64) int64 {
	var l, r int64
	l = 0
	r = int64(len(intArr) - 1)
	for l+1 < r {
		mid := l + (r-l)/2
		if intArr[mid] == src {
			return mid
		} else if src > intArr[mid] {
			r = mid
		} else {
			l = mid
		}
	}

	// 未匹配上，且在所有元素之前
	if src > intArr[l] {
		return 0
	}

	// 为匹配上，且在所有元素之后
	if src < intArr[r] {
		return int64(len(intArr) - 1)
	}

	// 在l和r之间
	return l + 1
}

// times: time segments, threshold: the valid time diff (s)
func timeMerge(times []entity.TimeQuantum, threshold int64) []entity.TimeQuantum {
	if len(times) == 0 {
		return times
	}

	var res []entity.TimeQuantum
	linq.From(times).OrderByT(func(quantum entity.TimeQuantum) int64 {
		return quantum.StartTime.Unix()
	}).ToSlice(&times)

	curStart := times[0].StartTime
	curEnd := times[0].EndTime

	for i := 1; i < len(times); i++ {
		if times[i].StartTime.Unix()-curEnd.Unix() > threshold { // 超过0s认为应该分为两段时间记录
			// 存储当前这段时间
			res = append(res, entity.TimeQuantum{
				StartTime: curStart,
				EndTime:   curEnd,
			})

			curStart = times[i].StartTime
			curEnd = times[i].EndTime

			continue
		}

		curEnd = times[i].EndTime
	}

	res = append(res, entity.TimeQuantum{
		StartTime: curStart,
		EndTime:   curEnd,
	})

	return res
}

func genBrowseDate(browseTime time.Time) string {
	dateLayout := "2006-01-02"
	return browseTime.Format(dateLayout)
}

func genBeginningDate(browseTime time.Time, location *time.Location) time.Time {
	dateLayout := "2006-01-02"
	dateTimeLayout := "2006-01-02"
	date := browseTime.Format(dateLayout)

	res, _ := time.ParseInLocation(dateTimeLayout, date, location)

	return res
}

func genTtlSec(currentTime time.Time, location *time.Location) int64 {
	nextDayStr := currentTime.AddDate(0, 0, 1).Format("2006-01-02")
	nextZeroTime, _ := time.ParseInLocation("2006-01-02", nextDayStr, location)

	expireTime := nextZeroTime.Unix() - currentTime.Unix()
	return expireTime
}
