package service

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"time"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IUser interface {
	//GenClientId 生成clientId
	GenClientId() string

	//SetCookie 设置cookie
	SetCookie(ctx *gin.Context, key, content string)

	//ParseJwtInfo 解析jwt
	ParseJwtInfo(ctx *gin.Context) (userId int64, err gone.Error)

	ParseJwt(token, secret string) (int64, error)

	//Login 用户登录
	Login(ctx *gin.Context, user *entity.UserSimple) (rst *entity.LoginRst, err error)

	Logout(ctx *gin.Context, userId int64, clientId string) error

	GenQrToken() (*entity.QrTokenWarp, error)

	SetUserWithLoginToken(loginToken, nickName string, avatar entity.Url, invitedBy *int64, fromApp string, source string) (rst *entity.UserSimple, err error)

	GetQrTokenParams(token string, redirectUrl string) (tokenExpired time.Time, authUrl string, err error)

	AuthByCode(qrToken, code string) (*entity.UserSimple, error)

	//KickUserById 踢用户下线
	KickUserById(userId int64, msg ...string)

	GetUserByOpenId(openId string) (*entity.UserSimple, error)

	GetUserSettings(userId int64) (*domain.UserSettings, error)

	CheckUserExists(id int64) (bool, error)

	GetUserById(id int64) (*entity.User, error)

	GetUserInBatch(ids []int64) ([]*entity.User, error)

	GetUserSimpleInBatch(ids []int64) ([]*entity.UserSimple, error)

	UpdateUserInfo(id int64, nickname, avatar string) error

	UpdateWorkOffTime(userId int64, offTime string) error

	UpdateBossKey(userId int64, bossKey string) error

	UpdateOnline(userId int64, online bool, t time.Time) error

	CalculateEarlierThan(offTime string) (*domain.EarlierThan, error)

	UpdateUserSettings(userSetting entity.UserSettings) error

	GetSimpleUserSettings(userId int64, componentNames []entity.ComponentNameEnum) ([]*entity.ComponentElement, error)

	ListOnlineUsers() ([]*entity.User, error)

	UpdateConnectTime(userId int64) error

	AccumulateTotalAccessDuration(userId int64) error

	AccumulateTotalBrowseDuration(userId int64, timeQuantum []entity.TimeQuantum) error

	GenMoyuDetail(userId int64) (*domain.MoyuDetail, error)

	UpdateWorkInfo(userId int64, info domain.WorkSettings) error

	GetBrowseRankingList(top int, period entity.StatPeriod) ([]*domain.RankInfo, error)

	PostEvent(event *entity.UserEvent) error
}
