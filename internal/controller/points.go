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
)

//go:gone
func NewPointsController() gone.Goner {
	return &pointsController{}
}

type pointsController struct {
	*Base         `gone:"*"`
	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter `gone:"router-auth"`

	pointsStrategy service.IPointStrategy `gone:"*"`
	points         service.IPoints        `gone:"*"`
}

func (ctr *pointsController) Mount() gin.MountError {
	ctr.
		AuthRouter.
		Group("/points").
		GET("", ctr.pagePointFlow).
		GET("/total", ctr.getTotalPoints).
		POST("", ctr.WithRole(entity.PlanetRoleAdmin, ctr.addPoints))

	ctr.AuthRouter.
		GET("/users/points-ranking", ctr.getPointsRankingTotal).
		GET("/users/points-ranking/daily", ctr.getPointsRankingDaily)

	return nil
}

func (ctr *pointsController) pagePointFlow(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	query, err := page.ParseQuery(ctx)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return ctr.points.PagePointFlow(query, userId)
}

func (ctr *pointsController) addPoints(ctx *gin.Context) (any, error) {
	//var param struct {
	//	UserId      int64  `json:"userId" binding:"required"`
	//	Points      int64  `json:"points" binding:"required"`
	//	Description string `json:"description" binding:"required"`
	//}
	//
	//err := ctx.ShouldBindJSON(&param)
	//if err != nil {
	//	return nil, gin.NewParameterError(err.Error())
	//}
	//
	//return nil, ctr.points.ApplyPoints(param.)
	return nil, nil
}

func (ctr *pointsController) getTotalPoints(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	points, err := ctr.points.GetUserPoints(userId)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"total":          points,
		"withdrawAmount": fmt.Sprintf("%.2f", float32(points)/5*0.0075),
	}, nil
}

func (ctr *pointsController) getPointsRankingTotal(ctx *gin.Context) (any, error) {
	var param struct {
		Top int `form:"top"`
	}

	if err := ctx.ShouldBindQuery(&param); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	ranking, err := ctr.points.RankingPointsTotal(param.Top)
	return entity.ListWrap{
		List: ranking,
	}, err
}

func (ctr *pointsController) getPointsRankingDaily(ctx *gin.Context) (any, error) {
	var param struct {
		Top int `form:"top"`
	}

	if err := ctx.ShouldBindQuery(&param); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	ranking, err := ctr.points.RankingPointsToday(param.Top)
	return entity.ListWrap{
		List: ranking,
	}, err
}
