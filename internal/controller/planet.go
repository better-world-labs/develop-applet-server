package controller

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"strconv"
)

//go:gone
func NewPlanetController() gone.Goner {
	return &planetController{}
}

type planetController struct {
	*Base `gone:"*"`

	logrus.Logger  `gone:"gone-logger"`
	AuthRouter     gin.IRouter     `gone:"router-auth"`
	UserService    service.IUser   `gone:"*"`
	PlanetsService service.IPlanet `gone:"*"`
}

func (ctr *planetController) Mount() gin.MountError {
	ctr.
		AuthRouter.
		Group("/planets").
		GET("/:planetId/members-count", ctr.countPlanetMembers).
		GET("/:planetId/members", ctr.pagePlanetMembers).
		GET("/:planetId/members/me", ctr.getMyMemberInfo).
		PUT("/:planetId/msg", ctr.WithRole(entity.PlanetRoleAdmin, ctr.updatePlanetMessage)).
		PUT("/:planetId/members/status", ctr.WithRole(entity.PlanetRoleAdmin, ctr.updateMemberStatus)).
		PUT("/:planetId/members/role", ctr.WithRole(entity.PlanetRoleRoot, ctr.updateMemberRole)).
		GET("/:planetId/msg", ctr.getPlanet)
	return nil
}

func (ctr *planetController) updatePlanetMessage(ctx *gin.Context) (interface{}, error) {

	type Req struct {
		Icon       entity.Url `json:"icon"`
		FrontCover entity.Url `json:"frontCover"`
		Name       string     `json:"name"`
	}

	var req Req
	if err := ctx.BindJSON(&req); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	planetId, err := strconv.ParseInt(ctx.Param("planetId"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(planetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if member.Role < entity.PlanetRoleAdmin {
		return nil, gin.NewParameterError(PermissionError)
	}

	return nil, ctr.PlanetsService.UpdatePlanetMessage(req.Icon, req.FrontCover, req.Name, planetId)
}

func (ctr *planetController) countPlanetMembers(ctx *gin.Context) (interface{}, error) {
	planetId, err := utils.CtxPathParamInt64(ctx, "planetId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	count, err := ctr.PlanetsService.CountPlanetMembers(planetId)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"count": count,
	}, nil
}

func (ctr *planetController) pagePlanetMembers(ctx *gin.Context) (interface{}, error) {
	query, err := page.ParseQuery(ctx)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	planetId, err := utils.CtxPathParamInt64(ctx, "planetId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userIdCondition := ctx.Query("userId")
	return ctr.PlanetsService.PagePlanetMembers(planetId, query, userIdCondition)
}

func (ctr *planetController) updateMemberRole(ctx *gin.Context) (interface{}, error) {
	var param struct {
		UserIds []int64            `json:"userIds" binding:"required"`
		Role    *entity.PlanetRole `json:"role" binding:"required,min=0,max=2"`
	}

	planetId, err := utils.CtxPathParamInt64(ctx, "planetId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	err = ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	for _, v := range param.UserIds {
		if v == userId {
			return nil, gin.NewParameterError(fmt.Sprintf("change your role by yourself is not allowed, user id [%d]. ", userId))
		}
	}

	return nil, ctr.PlanetsService.UpdateMembersRole(planetId, param.UserIds, *param.Role)
}

func (ctr *planetController) updateMemberStatus(ctx *gin.Context) (interface{}, error) {
	var param struct {
		UserIds []int64                    `json:"userIds" binding:"required"`
		Status  *entity.PlanetMemberStatus `json:"status" binding:"required,min=0,max=1"`
	}

	planetId, err := utils.CtxPathParamInt64(ctx, "planetId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	err = ctx.ShouldBindJSON(&param)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	for _, v := range param.UserIds {
		if v == userId {
			return nil, gin.NewParameterError(fmt.Sprintf("change your status by yourself is not allowed, user id [%d]. ", userId))
		}
	}

	return nil, ctr.PlanetsService.UpdateMembersStatus(planetId, param.UserIds, *param.Status)
}

func (ctr *planetController) getPlanet(ctx *gin.Context) (interface{}, error) {
	planetId, err := strconv.ParseInt(ctx.Param("planetId"), 10, 64)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.PlanetsService.GetPlanet(planetId)
}

func (ctr *planetController) getMyMemberInfo(ctx *gin.Context) (interface{}, error) {
	planetId, err := utils.CtxPathParamInt64(ctx, "planetId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	userId := utils.CtxMustGetUserId(ctx)
	member, err := ctr.PlanetsService.GetPlanetMemberByUserId(planetId, userId)
	if err != nil {
		return nil, gin.ToError(err)
	}

	if member == nil {
		return nil, gin.NewParameterError("member not found")
	}

	return map[string]any{
		"userId": userId,
		"role":   member.Role,
		"status": member.Status,
	}, nil
}
