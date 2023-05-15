package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	_interface "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"strconv"
	"time"
)

const (
	DefaultTimeLayout = "15:04:05"
)

//go:gone
func NewUserController() gone.Goner {
	return &userController{}
}

type userController struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter                 `gone:"router-auth"`
	PubRouter     gin.IRouter                 `gone:"router-pub"`
	UserService   service.IUser               `gone:"*"`
	LikeComment   service.ILikeCommentMiniApp `gone:"*"`
	MiniApp       service.IMiniApp            `gone:"*"`
	Points        service.IPoints             `gone:"*"`
}

func (ctr *userController) Mount() gin.MountError {
	ctr.AuthRouter.
		Group("/users").
		POST("/list", ctr.getUsers).
		GET("/me/user-settings", ctr.getUserSettings).
		POST("/me/simple-settings", ctr.getSimpleUserSettings).
		PUT("/me/user-settings", ctr.updateUserSettings).
		PUT("/me/info", ctr.updateUserInfo).
		GET("/me/info", ctr.getUserInfo).
		PUT("/me/work-off-time", ctr.updateWorkOffTime).
		PUT("/me/boss-key", ctr.updateBossKey).
		GET("/off-time-earlier", ctr.calculateEarlierThan).
		GET("/statistic", ctr.userStatistic).
		POST("/me/browse-duration", ctr.updateBrowseDuration).
		GET("/moyu-detail", ctr.getMoyuDetail).
		PUT("/me/work-settings", ctr.updateWorkSettings).
		GET("/moyu-time-ranking", ctr.getMoyuTimeRank)
	return nil
}

func (ctr *userController) getUserSettings(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	return ctr.UserService.GetUserSettings(userId)
}

func (ctr *userController) getSimpleUserSettings(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	var req struct {
		ComponentNames []entity.ComponentNameEnum `json:"componentNames"`
	}
	err := ctx.BindJSON(&req)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	componentSettings, err := ctr.UserService.GetSimpleUserSettings(userId, req.ComponentNames)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return entity.ListWrap{List: componentSettings}, nil
}

func (ctr *userController) updateWorkOffTime(ctx *gin.Context) (interface{}, error) {
	type req struct {
		Time string `json:"time"`
	}

	var reqBody req
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}
	userId := utils.CtxMustGetUserId(ctx)

	_, err = time.Parse(DefaultTimeLayout, reqBody.Time)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, ctr.UserService.UpdateWorkOffTime(userId, reqBody.Time)
}

func (ctr *userController) updateBossKey(ctx *gin.Context) (interface{}, error) {
	type req struct {
		BossKey string `json:"bossKey"`
	}

	var reqBody req
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}
	userId := utils.CtxMustGetUserId(ctx)

	return nil, ctr.UserService.UpdateBossKey(userId, reqBody.BossKey)
}

func (ctr *userController) getUserInfo(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)

	user, err := ctr.UserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return &entity.UserInfo{
		UserSimple: entity.UserSimple{
			Id:       user.Id,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
		LoginAt:     user.LoginAt,
		LastLoginAt: user.LastLoginAt,
		InvitedBy:   utils.PointerValue(user.InvitedBy),
		Points:      _interface.PointsStrategyFirstLogin,
	}, nil
}

func (ctr *userController) updateUserInfo(ctx *gin.Context) (interface{}, error) {
	var param struct {
		Nickname string `json:"nickname" binding:"required"`
		Avatar   string `json:"avatar" binding:"required"`
	}

	userId := utils.CtxMustGetUserId(ctx)
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	if len([]rune(param.Nickname)) > 15 {
		return nil, gin.NewParameterError("invalid nickname: too long")
	}

	if len([]rune(param.Avatar)) == 0 { // TODO required 了，可以省掉
		return nil, gin.NewParameterError("invalid avatar: avatar cannot be empty")
	}

	return nil, ctr.UserService.UpdateUserInfo(userId, param.Nickname, param.Avatar)
}

func (ctr *userController) getUsers(ctx *gin.Context) (interface{}, error) {
	var param struct {
		Ids []int64 `json:"ids" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	users, err := ctr.UserService.GetUserSimpleInBatch(param.Ids)
	return entity.ListWrap{List: users}, err
}

func (ctr *userController) calculateEarlierThan(ctx *gin.Context) (interface{}, error) {
	offTime := ctx.Query("offTime")
	_, err := time.Parse(DefaultTimeLayout, offTime)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.UserService.CalculateEarlierThan(offTime)
}

func (ctr *userController) updateUserSettings(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	var req struct {
		AppearanceTheme entity.AppearanceThemeEnum `json:"appearanceTheme"`
		SiteSettings    *entity.SiteSettings       `json:"siteSettings"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	var settings []byte
	if req.SiteSettings != nil {
		settings, err = json.Marshal(*req.SiteSettings)
		if err != nil {
			return nil, gin.NewParameterError(err.Error())
		}
	}

	userSetting := entity.UserSettings{UserId: userId, AppearanceTheme: req.AppearanceTheme, SiteSettings: string(settings)}

	return nil, ctr.UserService.UpdateUserSettings(userSetting)
}

func (ctr *userController) updateBrowseDuration(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)

	var req struct {
		TimeQuantum []entity.TimeQuantum `json:"timeQuantum"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	if len(req.TimeQuantum) == 0 {
		return nil, gin.NewParameterError("no time-quantum found in request body. ")
	}

	return nil, ctr.UserService.AccumulateTotalBrowseDuration(userId, req.TimeQuantum)
}

func (ctr *userController) getMoyuDetail(ctx *gin.Context) (interface{}, error) {
	userId, err := strconv.ParseInt(ctx.Query("userId"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.UserService.GenMoyuDetail(userId)
}

func (ctr *userController) updateWorkSettings(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)

	var params domain.WorkSettings
	err := ctx.BindJSON(&params)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	if params.MonthlyWorkingDays < 1 ||
		params.MonthlyWorkingDays > 31 ||
		params.MonthlySalary > 1000000 {
		return nil, gin.NewParameterError(fmt.Sprintf("invalid value of day [%d] or salary [%d]. ", params.MonthlyWorkingDays, params.MonthlySalary))
	}

	return nil, ctr.UserService.UpdateWorkInfo(userId, params)
}

func (ctr *userController) getMoyuTimeRank(ctx *gin.Context) (interface{}, error) {
	top, err := strconv.Atoi(ctx.Query("top"))
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	period, err := strconv.ParseInt(ctx.Query("period"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	res, err := ctr.UserService.GetBrowseRankingList(top, entity.StatPeriod(period))
	if err != nil {
		return nil, err
	}
	return entity.ListWrap{List: res}, nil
}

func (ctr *userController) postEvent(ctx *gin.Context) (any, error) {
	var event entity.UserEvent

	err := ctx.ShouldBindJSON(&event)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, ctr.UserService.PostEvent(&event)
}

func (ctr *userController) userStatistic(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	likes, err := ctr.LikeComment.CountUserAppsLikes(userId)
	if err != nil {
		return nil, err
	}

	userInfo, err := ctr.UserService.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	apps, err := ctr.MiniApp.CountUserCreatedApps(userId)
	if err != nil {
		return nil, err
	}

	runTimes, err := ctr.MiniApp.CountUsersAppsRuntimes(userId)
	if err != nil {
		return nil, err
	}

	points, err := ctr.Points.GetUserPoints(userId)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"apps":           apps,
		"registeredDays": userInfo.RegisteredDays(),
		"points":         points,
		"appUses":        runTimes,
		"appLikes":       likes,
	}, nil
}
