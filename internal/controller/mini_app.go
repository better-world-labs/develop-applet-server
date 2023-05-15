package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"io"
	"net/http"
	"strings"
	"time"
)

//go:gone
func NewMiniAppController() gone.Goner {
	return &miniAppController{}
}

type miniAppController struct {
	*Base `gone:"*"`

	logrus.Logger `gone:"gone-logger"`
	AuthRouter    gin.IRouter                 `gone:"router-auth"`
	PubRouter     gin.IRouter                 `gone:"router-pub"`
	svc           service.IMiniApp            `gone:"*"`
	likeComment   service.ILikeCommentMiniApp `gone:"*"`
	recommend     service.IRecommendMiniApp   `gone:"*"`
	system        service.ISystemConfig       `gone:"*"`
}

func (con *miniAppController) Mount() gin.MountError {
	con.PubRouter.
		GET("/apps", con.listApps).
		GET("/apps/:uuid", con.getAppById).
		GET("/app-categories", con.getAppCategories).
		GET("/app-tabs", con.getAppTabs)

	con.AuthRouter.
		GET("/users/:userId/apps", con.listUsersApps).
		GET("/apps/mine", con.listAppsByUser).
		GET("/apps/collected", con.listCollectedApps).
		POST("/apps/:uuid/collect", con.collectApp).
		POST("/apps/is-collected", con.checkAppsCollected).
		PUT("/apps/:uuid", con.saveApp).
		POST("/apps/:uuid/run", con.runApp).
		POST("/apps/:uuid/comments", con.addComment).
		GET("/apps/:uuid/comments", con.listComments).
		POST("/apps/:uuid/like", con.likeApp).
		POST("/apps/is-liked", con.isAppLiked).
		POST("/apps/:uuid/recommend", con.recommendApp).
		POST("/apps/is-recommended", con.isAppRecommended).
		POST("/outputs/:id/like", con.likeOutput).
		GET("/outputs/likes", con.listOutputsLikes).
		GET("/apps/:uuid/like", con.getAppLike).
		DELETE("/apps/:uuid", con.deleteApp).
		GET("/apps/:uuid/outputs", con.getAppOutputs).
		GET("/ai-models", con.listAIModels).
		GET("", con.listAIModels)

	return nil
}

func (con *miniAppController) listApps(ctx *gin.Context) (any, error) {
	var category struct {
		Category int64 `form:"category"`
	}

	err := ctx.ShouldBindQuery(&category)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	apps, err := con.svc.ListApps(category.Category)
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{List: apps}, nil
}

func (con *miniAppController) listAppsByUser(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	var category struct {
		Category int64 `form:"category"`
	}

	err := ctx.ShouldBindQuery(&category)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	apps, err := con.svc.ListAppsByUser(userId, category.Category)
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{List: apps}, nil
}

func (con *miniAppController) getAppById(ctx *gin.Context) (any, error) {
	uuid := ctx.Param("uuid")
	app, has, err := con.svc.GetAppDetailByUuid(uuid)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return app, nil
}

func (con *miniAppController) saveApp(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	uuid := ctx.Param("uuid")

	var app entity.MiniApp
	app.Uuid = uuid
	app.CreatedBy = userId
	err := ctx.ShouldBindJSON(&app)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	saved, err := con.svc.SaveApp(&app)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (con *miniAppController) runApp(ctx *gin.Context) (any, error) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")

	var done struct {
		Code int    `json:"code"`
		Msg  string `json:"msg,omitempty"`
	}

	f, ok := ctx.Writer.(http.Flusher)
	if !ok {
		return nil, nil
	}

	uuid := ctx.Param("uuid")
	userId := utils.CtxMustGetUserId(ctx)
	var param entity.MiniAppRunParam
	err := ctx.ShouldBindJSON(&param)
	if err != nil {
		done.Code = -1
		done.Msg = err.Error()
		_, _ = io.WriteString(ctx.Writer, "event: done\n")
		doneJson, _ := json.Marshal(done)
		_, _ = io.WriteString(ctx.Writer, fmt.Sprintf("data: %s\n\n", doneJson))
		f.Flush()
		return nil, nil
	}

	reader, err := con.svc.RunApp(userId, uuid, param)
	if err != nil {
		if goneErr, ok := err.(gone.Error); ok {
			done.Code = goneErr.Code()
			done.Msg = goneErr.Error()
		} else {
			done.Code = -1
			done.Msg = err.Error()
		}

		_, _ = io.WriteString(ctx.Writer, "event: done\n")
		doneJson, _ := json.Marshal(done)
		_, _ = io.WriteString(ctx.Writer, fmt.Sprintf("data: %s\n\n", doneJson))
		f.Flush()
		return nil, nil
	}

	for {
		read, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			done.Code = -1
			done.Msg = err.Error()
			break
		}

		json, err := json.Marshal(read)
		if err != nil {
			done.Code = -1
			done.Msg = err.Error()
			break
		}

		_, _ = io.WriteString(ctx.Writer, "event: data\n")
		_, _ = io.WriteString(ctx.Writer, fmt.Sprintf("data: %s\n\n", json))
		f.Flush()
	}

	_, _ = io.WriteString(ctx.Writer, "event: done\n")
	doneJson, err := json.Marshal(done)
	_, _ = io.WriteString(ctx.Writer, fmt.Sprintf("data: %s\n\n", doneJson))
	f.Flush()
	reader.Close()

	return nil, nil
}

func (con *miniAppController) getAppOutputs(ctx *gin.Context) (any, error) {
	uuid := ctx.Param("uuid")
	var query page.StreamQuery
	if err := query.BindQuery(ctx); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	outputs, total, err := con.svc.PageOpenedAppOutputsByAppId(query, uuid)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"total":      total,
		"list":       outputs.GetList(),
		"nextCursor": outputs.GetNextCursor(),
	}, nil
}

func (con *miniAppController) listAIModels(ctx *gin.Context) (any, error) {
	models, err := con.svc.ListAIModels()
	if err != nil {
		return nil, err
	}

	categories, err := con.svc.ListAIModelCategories()
	if err != nil {
		return nil, err
	}

	for _, m := range models {
		if category, has := getCategory(categories, m.Category); has {
			m.Category = 0
			category.Models = append(category.Models, m)
		}

	}

	return entity.ListWrap{
		List: categories,
	}, nil
}

func (con *miniAppController) getAppCategories(ctx *gin.Context) (any, error) {
	categories, err := con.svc.ListAppsCategories()
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{
		List: categories,
	}, nil
}

func (con *miniAppController) getAppTabs(ctx *gin.Context) (any, error) {
	value, err := con.system.Get("MINI_APP_HOME_TABS")
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (con *miniAppController) deleteApp(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	appId := ctx.Param("uuid")
	return con.svc.DeleteApp(userId, appId), nil
}

func (con *miniAppController) addComment(ctx *gin.Context) (any, error) {
	appId := ctx.Param("uuid")
	userId := utils.CtxMustGetUserId(ctx)
	comment := entity.MiniAppComment{
		AppId:     appId,
		CreatedBy: userId,
		CreatedAt: time.Now(),
	}

	err := ctx.ShouldBindJSON(&comment)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, con.likeComment.CreateAppComment(comment)
}

func (con *miniAppController) likeApp(ctx *gin.Context) (any, error) {
	appId := ctx.Param("uuid")
	userId := utils.CtxMustGetUserId(ctx)
	like := entity.MiniAppLike{
		AppId:     appId,
		CreatedBy: userId,
		UpdatedAt: time.Now().UnixMilli(),
	}

	err := ctx.ShouldBindJSON(&like)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, con.likeComment.LikeApp(like)
}

func (con *miniAppController) recommendApp(ctx *gin.Context) (any, error) {
	appId := ctx.Param("uuid")
	userId := utils.CtxMustGetUserId(ctx)
	recommend := entity.MiniAppRecommend{
		AppId:     appId,
		CreatedBy: userId,
		UpdatedAt: time.Now().UnixMilli(),
	}

	err := ctx.ShouldBindJSON(&recommend)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, con.recommend.RecommendApp(recommend)
}

func (con *miniAppController) listComments(ctx *gin.Context) (any, error) {
	appId := ctx.Param("uuid")
	comments, err := con.likeComment.ListAppComments(appId)
	if err != nil {
		return nil, err
	}

	return entity.ListWrap{
		List: comments,
	}, nil
}

func (con *miniAppController) getAppLike(ctx *gin.Context) (any, error) {
	appId := ctx.Param("uuid")
	userId := utils.CtxMustGetUserId(ctx)

	like, err := con.likeComment.GetAppLike(appId, userId)

	return struct {
		Like bool `json:"like"`
	}{Like: like.Like}, err
}

func (con *miniAppController) likeOutput(ctx *gin.Context) (any, error) {
	outputId := ctx.Param("id")

	userId := utils.CtxMustGetUserId(ctx)
	like := entity.MiniAppOutputLike{
		UserOutputLikeState: entity.UserOutputLikeState{
			OutputId: outputId,
		},
		CreatedBy: userId,
		UpdatedAt: time.Now().UnixMilli(),
	}

	err := ctx.ShouldBindJSON(&like)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return nil, con.likeComment.LikeAppOutput(like)
}

func (con *miniAppController) listOutputsLikes(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	ids, has := ctx.GetQuery("outputIds")
	if !has {
		return nil, gin.NewParameterError("invalid outputIds")
	}

	outputIds := strings.Split(ids, ",")
	states, err := con.likeComment.ListUserOutputLikeState(outputIds, userId)
	return entity.ListWrap{
		List: states,
	}, err
}

func (con *miniAppController) listCollectedApps(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	apps, err := con.svc.ListCollectedApps(userId)
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return entity.ListWrap{
		List: apps,
	}, nil
}

func (con *miniAppController) collectApp(ctx *gin.Context) (any, error) {
	uuid := ctx.Param("uuid")
	userId := utils.CtxMustGetUserId(ctx)

	var param struct {
		Collected bool `json:"collected"`
	}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	if param.Collected {
		return nil, con.svc.DoCollectApp(uuid, userId)
	} else {
		return nil, con.svc.DoUnCollectApp(uuid, userId)
	}
}

func (con *miniAppController) checkAppsCollected(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	var param struct {
		Uuids []string `json:"uuids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return con.svc.IsAppsCollected(param.Uuids, userId)
}

func (con *miniAppController) isAppLiked(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	var param struct {
		AppIds []string `json:"appIds" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return con.likeComment.IsAppsLiked(param.AppIds, userId)
}

func (con *miniAppController) isAppRecommended(ctx *gin.Context) (any, error) {
	userId := utils.CtxMustGetUserId(ctx)
	var param struct {
		AppIds []string `json:"appIds" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&param); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return con.recommend.IsAppsRecommended(param.AppIds, userId)
}

func (con *miniAppController) listUsersApps(ctx *gin.Context) (any, error) {
	userId, err := utils.CtxPathParamInt64(ctx, "userId")
	if err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	var query page.StreamQuery
	if err := query.BindQuery(ctx); err != nil {
		return nil, gin.NewParameterError(err.Error())
	}

	return con.svc.PageUsersApps(query, userId)
}
