package controller

import (
	"errors"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	businesserrors "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/error"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/domain"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"strconv"
	"time"
)

const (
	PermissionError = "Current user has no permission to modify it. "
)

//go:gone
func NewChannelController() gone.Goner {
	return &channelController{}
}

type channelController struct {
	*Base `gone:"*"`

	logrus.Logger  `gone:"gone-logger"`
	AuthRouter     gin.IRouter      `gone:"router-auth"`
	UserService    service.IUser    `gone:"*"`
	PlanetsService service.IPlanet  `gone:"*"`
	ChannelService service.IChannel `gone:"*"`
}

func (ctr *channelController) Mount() gin.MountError {

	ctr.
		AuthRouter.
		Group("/channels").
		DELETE("/:channelId/users/:userId", ctr.WithRole(entity.PlanetRoleAdmin, ctr.adminRemoveMember)).
		GET("/groups", ctr.listPlanetGroups).
		POST("/groups", ctr.createPlanetGroup).
		POST("/query-many", ctr.listChannelsById).
		PUT("/groups/group-name", ctr.updateGroupName).
		DELETE("/groups/:id", ctr.deletePlanetGroup).
		PUT("/groups/sort", ctr.updatePlanetGroupChannelSort).
		GET("", ctr.listChannels).
		POST("", ctr.createChannel).
		GET("/:channelId", ctr.getChannel).
		POST("/:channelId/apply", ctr.applyChannel).
		PUT("/:channelId/notice", ctr.WithRole(entity.PlanetRoleAdmin, ctr.updateNotice)).
		GET("/:channelId/member-state", ctr.getMemberState).
		PUT("/channel-name", ctr.updateChannelName).
		DELETE("/:channelId", ctr.deleteChannel).
		GET("/unread-msg-num", ctr.getUnreadNumDetail).
		GET("/:channelId/members", ctr.listChannelMembers).
		GET("/:channelId/last-read", ctr.getUserLastReadMsgId)

	ctr.AuthRouter.PUT("/groups/sort", ctr.updateChannelGroupsSort)
	return nil
}

func (ctr *channelController) deletePlanetGroup(ctx *gin.Context) (interface{}, error) {
	groupId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	group, err := ctr.ChannelService.GetChannelGroupByGroupId(groupId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if group == nil || group.Id == 0 {
		return nil, gin.NewParameterError(fmt.Sprintf("group with id [%d] not exist. ", groupId))
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(group.PlanetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	return nil, ctr.ChannelService.DeleteChannelGroup(groupId)
}

func (ctr *channelController) listPlanetGroups(context *gin.Context) (interface{}, error) {

	planetId, err := strconv.ParseInt(context.Query("planetId"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.ChannelService.ListChannelGroups(planetId)
}

func (ctr *channelController) createPlanetGroup(ctx *gin.Context) (interface{}, error) {

	type req struct {
		PlanetId int64  `json:"planetId"`
		Name     string `json:"name"`
		Icon     string `json:"icon"`
	}

	var reqBody req
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(reqBody.PlanetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	return ctr.ChannelService.CreateChannelGroup(userId, reqBody.PlanetId, reqBody.Name, reqBody.Icon)
}

func (ctr *channelController) updateGroupName(ctx *gin.Context) (interface{}, error) {

	type req struct {
		PlanetId       int64  `json:"planetId"`
		ChannelGroupId int64  `json:"channelGroupId"`
		Name           string `json:"name"`
	}

	var reqBody req
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(reqBody.PlanetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	return nil, ctr.ChannelService.UpdateChannelGroupName(reqBody.ChannelGroupId, reqBody.Name)
}

func (ctr *channelController) updatePlanetGroupChannelSort(ctx *gin.Context) (interface{}, error) {
	type req struct {
		GroupId int64   `json:"groupId"`
		Sorted  []int64 `json:"sortedChannelIds"`
	}

	var reqBody req
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	if reqBody.Sorted == nil || len(reqBody.Sorted) == 0 {
		return nil, gin.NewParameterError(err.Error())
	}

	group, err := ctr.ChannelService.GetChannelGroupByGroupId(reqBody.GroupId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(group.PlanetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	return nil, ctr.ChannelService.UpdateChannelSort(reqBody.GroupId, reqBody.Sorted)
}

func (ctr *channelController) updateChannelGroupsSort(ctx *gin.Context) (interface{}, error) {
	var param struct {
		PlanetId int64   `json:"planetId" binding:"required"`
		Sorted   []int64 `json:"sortedGroupIds" binding:"required"`
	}

	err := ctx.BindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(param.PlanetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if member == nil {
		return nil, gin.ToError(errors.New("planet not found"))
	}

	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	return nil, ctr.ChannelService.UpdateChannelGroupSort(param.PlanetId, param.Sorted)
}

func (ctr *channelController) listChannels(ctx *gin.Context) (interface{}, error) {
	planetId, err := strconv.ParseInt(ctx.Query("planetId"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.ChannelService.ListChannels(planetId)
}

func (ctr *channelController) applyChannel(ctx *gin.Context) (interface{}, error) {
	var param struct {
		Reason string `json:"reason" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	channelId, err := utils.CtxPathParamInt64(ctx, "channelId")
	if err != nil {
		return nil, gin.NewParameterError("invalid channelId")
	}

	return nil, ctr.ChannelService.ApplyPrivateChannel(userId, channelId, param.Reason)
}

func (ctr *channelController) createChannel(ctx *gin.Context) (interface{}, error) {
	type req struct {
		PlanetId  int64              `json:"planetId"`
		Name      string             `json:"name"`
		Icon      string             `json:"icon"`
		Type      entity.ChannelType `json:"type" binding:"min=1,max=2"`
		GroupId   int64              `json:"groupId"`
		Mute      bool               `json:"mute"`
		ExpiresIn time.Duration      `json:"expiresIn"`
	}

	var reqBody req
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(reqBody.PlanetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	var channelId int64
	if reqBody.Type == entity.ChannelTypeNormal {
		channelId, err = ctr.ChannelService.CreateNormalChannel(reqBody.Name, reqBody.Icon, reqBody.GroupId, reqBody.PlanetId, userId, reqBody.Mute)
		if err != nil {
			return nil, err
		}
	}

	if reqBody.Type == entity.ChannelTypePrivate {
		channelId, err = ctr.ChannelService.CreatePrivateChannel(reqBody.Name, reqBody.Icon, reqBody.GroupId, reqBody.PlanetId, userId, reqBody.Mute, reqBody.ExpiresIn*time.Second)
		if err != nil {
			return nil, err
		}
	}

	return struct {
		ChannelId int64 `json:"channelId"`
	}{ChannelId: channelId}, nil
}

func (ctr *channelController) updateChannelName(ctx *gin.Context) (interface{}, error) {

	type req struct {
		PlanetId  int64  `json:"planetId"`
		ChannelId int64  `json:"channelId"`
		Name      string `json:"name"`
	}

	var reqBody req
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(reqBody.PlanetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	return nil, ctr.ChannelService.UpdateChannelName(reqBody.ChannelId, reqBody.Name)
}

func (ctr *channelController) deleteChannel(ctx *gin.Context) (interface{}, error) {
	channelId, err := strconv.ParseInt(ctx.Param("channelId"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	channel, exists, err := ctr.ChannelService.GetChannelByChannelId(channelId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if !exists {
		return nil, gin.NewParameterError(fmt.Sprintf("channel with id [%d] not exist. ", channelId))
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(channel.PlanetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	return nil, ctr.ChannelService.DeleteChannel(channelId)
}

func (ctr *channelController) listChannelMembers(ctx *gin.Context) (interface{}, error) {
	channelId, err := strconv.ParseInt(ctx.Param("channelId"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.ChannelService.ListChannelMembers(channelId)
}

func (ctr *channelController) getUserLastReadMsgId(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	channelId, err := strconv.ParseInt(ctx.Param("channelId"), 10, 64)
	if err != nil {
		return nil, err
	}

	channelMember, err := ctr.ChannelService.GetLastReadMsgId(userId, channelId)
	if err != nil {
		return nil, err
	}

	if channelMember == nil {
		return nil, gin.NewParameterError(fmt.Sprintf("this user is not the member of chennel[%d]", channelId))
	}

	return domain.ChannelMemberLastReadMessageId{LastReadMessageId: channelMember.LastReadMessageId}, nil
}

func (ctr *channelController) getUnreadNumDetail(ctx *gin.Context) (interface{}, error) {
	userId := utils.CtxMustGetUserId(ctx)
	planetId, err := strconv.ParseInt(ctx.Query("planetId"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	cnts, err := ctr.ChannelService.GetAllUnReadCountByUser(userId, planetId)
	if err != nil {
		return nil, err
	}
	return entity.ListWrap{List: cnts}, nil
}

func (ctr *channelController) getMemberState(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)

	channelId, err := utils.CtxPathParamInt64(ctx, "channelId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	member, err := ctr.ChannelService.GetChannelMember(channelId, userId)
	if err != nil {
		return nil, err
	}

	var state struct {
		State uint `json:"state"`
	}

	if member != nil {
		state.State = uint(member.State)
	}

	return state, nil
}

func (ctr *channelController) adminRemoveMember(ctx *gin.Context) (any, error) {
	channelId, err := utils.CtxPathParamInt64(ctx, "channelId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	memberId, err := utils.CtxPathParamInt64(ctx, "userId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, ctr.ChannelService.RemoveChannelMember(channelId, memberId)
}

func (ctr *channelController) getChannel(ctx *gin.Context) (any, error) {
	channelId, err := utils.CtxPathParamInt64(ctx, "channelId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	channel, exists, err := ctr.ChannelService.GetChannelByChannelId(channelId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, businesserrors.ErrorChannelNotFound
	}

	return channel, nil
}

func (ctr *channelController) listChannelsById(ctx *gin.Context) (any, error) {
	param := struct {
		Ids []int64 `json:"ids" binding:"required"`
	}{}

	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	channels, err := ctr.ChannelService.ListChannelsByIds(param.Ids)
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{List: channels}, nil
}

func (ctr *channelController) updateNotice(ctx *gin.Context) (any, error) {
	channelId, err := utils.CtxPathParamInt64(ctx, "channelId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	var param struct {
		Notice string `json:"notice"`
	}

	err = ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, ctr.ChannelService.UpdateNotice(userId, channelId, param.Notice)
}
