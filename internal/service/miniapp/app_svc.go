package miniapp

import (
	"bytes"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	businesserrors "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/error"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"time"
)

type appSvc struct {
	gone.Goner
	Sender        emitter.Sender `gone:"gone-emitter"`
	logrus.Logger `gone:"gone-logger"`
	xorm.Engine   `gone:"gone-xorm"`

	pCategory    iCategoryPersistence          `gone:"*"`
	pAIModel     iAIModel                      `gone:"*"`
	pMiniApp     iAppPersistence               `gone:"*"`
	pCollection  iMiniAppCollectionPersistence `gone:"*"`
	userSvc      service.IUser                 `gone:"*"`
	statistic    service.IStatisticMiniApp     `gone:"*"`
	likeComments service.ILikeCommentMiniApp   `gone:"*"`

	aiRuntime *AIRuntime `gone:"*"`

	points service.IPointStrategy `gone:"*"`
}

//go:gone
func NewAppSvc() gone.Goner {
	return &appSvc{}
}

func (a appSvc) ListAppsCategories() ([]*entity.MiniAppCategory, error) {
	return a.pCategory.list()
}

func (a appSvc) IsAppsCollected(appIds []string, userId int64) (map[string]bool, error) {
	collectedAppIds, err := a.pCollection.getAppIds(userId)
	if err != nil {
		return nil, err
	}

	collectedSet := collection.NewSet[string]()
	collectedSet.AddAll(collectedAppIds)

	res := make(map[string]bool)
	for _, appId := range appIds {
		res[appId] = collectedSet.Contains(appId)
	}

	return res, nil
}

func (a appSvc) DoCollectApp(appIds string, userId int64) error {
	return a.pCollection.createIfNotExists(appIds, userId)
}

func (a appSvc) DoUnCollectApp(appIds string, userId int64) error {
	return a.pCollection.deleteIfExists(appIds, userId)
}

func (a appSvc) ListCollectedApps(userId int64) ([]*entity.MiniAppListDto, error) {
	appIds, err := a.pCollection.getAppIds(userId)
	if err != nil {
		return nil, err
	}

	app, err := a.pMiniApp.listAppByUuids(appIds)
	if err != nil {
		return nil, err
	}

	return a.transMiniAppToListDto(app)
}

func (a appSvc) ListApps(category int64) ([]*entity.MiniAppListDto, error) {
	app, err := a.pMiniApp.listApp(category)
	if err != nil {
		return nil, err
	}

	return a.transMiniAppToListDto(app)
}

func (a appSvc) ListAppsByUuids(uuids []string) ([]*entity.MiniAppListDto, error) {
	app, err := a.pMiniApp.listAppByUuids(uuids)
	if err != nil {
		return nil, err
	}

	return a.transMiniAppToListDto(app)
}

func (a appSvc) ListAppsByUser(userId, category int64) ([]*entity.MiniAppListDto, error) {
	app, err := a.pMiniApp.listAppByUser(userId, category)
	if err != nil {
		return nil, err
	}

	return a.transMiniAppToListDto(app)
}

func (a appSvc) GetAppById(id int64) (*entity.MiniApp, bool, error) {
	return a.pMiniApp.getAppById(id)
}

func (a appSvc) GetAppByUuid(uuid string) (*entity.MiniApp, bool, error) {
	return a.pMiniApp.getAppByUuid(uuid)
}

func (a appSvc) GetAppDetailByUuid(uuid string) (*entity.MiniAppDetailDto, bool, error) {
	app, has, err := a.pMiniApp.getAppByUuid(uuid)
	if err != nil {
		return nil, has, err
	}

	user, err := a.userSvc.GetUserById(app.CreatedBy)
	if err != nil {
		return nil, has, err
	}

	if !has {
		return nil, has, nil
	}

	statistic, err := a.statistic.GetByAppId(uuid)
	if err != nil {
		return nil, has, nil
	}

	likeMap, err := a.likeComments.GetAppLikeCountMap([]string{uuid})
	if err != nil {
		return nil, has, err
	}

	commentMap, err := a.likeComments.GetAppCommentCountMap([]string{uuid})
	if err != nil {
		return nil, has, err
	}

	collectMap, err := a.pCollection.countByAppIds([]string{uuid})
	if err != nil {
		return nil, has, err
	}

	if like, ok := likeMap[uuid]; ok {
		statistic.LikeTimes = int(like)
	}

	if comment, ok := commentMap[uuid]; ok {
		statistic.CommentTimes = int(comment)
	}

	if collect, ok := collectMap[uuid]; ok {
		statistic.CollectTimes = int(collect)
	}

	if err := a.fixAppPrice(app); err != nil {
		return nil, false, err
	}

	return &entity.MiniAppDetailDto{
		MiniApp: *app,
		CreatedBy: entity.UserSimple{
			Id:       user.Id,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
		MiniAppStatisticInfo: statistic.MiniAppStatisticInfo,
	}, true, err
}

func (a appSvc) SaveApp(app *entity.MiniApp) (*entity.MiniApp, error) {
	if app.Uuid == "" {
		return nil, gin.NewParameterError("uuid must set")
	}

	if has, _ := a.pCategory.checkExists(app.Category); !has {
		return nil, gin.NewParameterError("invalid category")
	}

	has, err := a.pMiniApp.checkAppExistsByUuid(app.Uuid)
	if err != nil {
		return nil, err
	}

	if !has {
		err := a.createApp(app)
		if err != nil {
			return nil, err
		}
	} else {
		err := a.pMiniApp.updateApp(app)
		if err != nil {
			return nil, err
		}
	}

	app, _, err = a.pMiniApp.getAppByUuid(app.Uuid)
	return app, err
}

func (a appSvc) DeleteApp(userId int64, appId string) error {
	app, has, err := a.pMiniApp.getAppByUuid(appId)
	if err != nil {
		return err
	}

	if !has {
		return nil
	}

	if app.CreatedBy != userId {
		return businesserrors.ErrorPermissionDenied
	}

	return a.Transaction(func(session xorm.Interface) error {
		err := a.pMiniApp.deleteApp(appId)
		if err != nil {
			return err
		}

		return a.pMiniApp.deleteOutputsByAppId(appId)
	})
}

func (a appSvc) RunApp(userId int64, uuid string, param entity.MiniAppRunParam) (reader *service.ChannelStreamTrunkReader[*entity.MiniAppOutputStreamChunk], err error) {
	app, has, err := a.pMiniApp.getAppByUuid(uuid)
	if err != nil {
		return
	}

	if !has {
		err = gin.NewParameterError("app not found")
		return
	}

	_ = a.Transaction(func(session xorm.Interface) error {
		_, err = a.points.ApplyPoints(userId, entity.StrategyArgUsingApp{
			Form: *app.Form,
		})
		if err != nil {
			return err
		}

		if userId != app.CreatedBy {
			_, err = a.points.ApplyPoints(app.CreatedBy, entity.StrategyArgAppUsed{
				App:       *app,
				RunUserId: userId,
			})
			if err != nil {
				return err
			}

		}
		reader, err = a.runApp(userId, app, param)
		return err
	})
	return
}

func (a appSvc) runApp(userId int64, app *entity.MiniApp, param entity.MiniAppRunParam) (*service.ChannelStreamTrunkReader[*entity.MiniAppOutputStreamChunk], error) {
	err := app.Input(param.Values)
	if err != nil {
		return nil, err
	}

	runtime := NewAppRuntime(app, a.aiRuntime)
	runtime.OnComplete(func(outputs map[string]entity.MiniAppOutputCore) {
		var b bytes.Buffer
		var output entity.MiniAppOutputCore
		for _, f := range app.Flow {
			if o, ok := outputs[f.Id]; ok {
				output.Type = entity.MiniAppOutputTypeText
				output.OutputId = o.OutputId
				b.WriteString(o.Content)
			}
		}

		output.Content = b.String()
		err := a.Sender.Send(&entity.AppRunDoneEvent{
			AppId:  app.Uuid,
			User:   userId,
			Param:  param,
			Output: output,
			Time:   time.Now(),
		})
		if err != nil {
			a.Errorf("Send AppRunDoneEvent error: %v\n", err)
		}
	})

	reader, err := runtime.Run()
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (a appSvc) CreateOutput(output *entity.MiniAppOutput) error {
	return a.pMiniApp.createOutput(output)
}

func (a appSvc) ListOpenedAppOutputsByAppId(uuid string) ([]*entity.MiniAppOutputDto, error) {
	outputs, err := a.pMiniApp.listOpenedOutputs(uuid)
	if err != nil {
		return nil, err
	}

	return a.transMiniAppOutputToDto(outputs)
}

func (a appSvc) ListAIModels() ([]*entity.MiniAppAiModel, error) {
	return a.pAIModel.list()
}

func (a appSvc) ListAIModelCategories() ([]*entity.MiniAppAiModelCategory, error) {
	return a.pAIModel.listCategory()
}

func (a appSvc) CheckAppExists(uuid string) (bool, error) {
	return a.pMiniApp.checkAppExists(uuid)
}

func (a appSvc) CheckOutputExists(outputId string) (bool, error) {
	return a.pMiniApp.checkOutputExists(outputId)
}

func (a appSvc) CheckAppRanByUser(appId string, userId int64) (bool, error) {
	return a.pMiniApp.checkOutputExistsByAppIdAndUser(appId, userId)
}

func (a appSvc) GetAppOutputById(outputId string) (*entity.MiniAppOutput, bool, error) {
	return a.pMiniApp.getOutputById(outputId)
}

func (a appSvc) transMiniAppOutputToDto(outputs []*entity.MiniAppOutput) ([]*entity.MiniAppOutputDto, error) {
	userIds := collection.Map(outputs, func(a *entity.MiniAppOutput) int64 {
		return a.CreatedBy
	})

	outputIds := collection.Map(outputs, func(a *entity.MiniAppOutput) string {
		return a.OutputId
	})

	users, err := a.userSvc.GetUserInBatch(userIds)
	if err != nil {
		return nil, err
	}

	likeMap, err := a.statistic.GetAppOutputMapByOutputIds(outputIds)
	if err != nil {
		return nil, err
	}

	userSimpleMap := collection.ToMap(users, func(u *entity.User) (int64, *entity.UserSimple) {
		return u.Id, &entity.UserSimple{
			Id:       u.Id,
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
		}
	})

	return collection.Map(outputs, func(a *entity.MiniAppOutput) *entity.MiniAppOutputDto {
		dto := &entity.MiniAppOutputDto{
			MiniAppOutput: *a,
		}

		if createdBy, ok := userSimpleMap[a.CreatedBy]; ok {
			dto.CreatedBy = *createdBy
		}

		if like, ok := likeMap[a.OutputId]; ok {
			dto.StatisticMiniAppOutput = *like
		}

		return dto
	}), nil
}

func (a appSvc) transMiniAppToListDto(apps []*entity.MiniAppBaseInfo) ([]*entity.MiniAppListDto, error) {
	appIds := collection.Map(apps, func(app *entity.MiniAppBaseInfo) string {
		return app.Uuid
	})
	outputs, err := a.pMiniApp.getLastNOutputByAppIds(appIds, 2)
	if err != nil {
		return nil, err
	}

	userIds := collection.Map(apps, func(a *entity.MiniAppBaseInfo) int64 {
		return a.CreatedBy
	})

	users, err := a.userSvc.GetUserInBatch(userIds)
	if err != nil {
		return nil, err
	}

	statisticMap, err := a.statistic.GetAppMapByAppIds(appIds)
	if err != nil {
		return nil, err
	}

	likeMap, err := a.likeComments.GetAppLikeCountMap(appIds)
	if err != nil {
		return nil, err
	}

	commentMap, err := a.likeComments.GetAppCommentCountMap(appIds)
	if err != nil {
		return nil, err
	}

	collectMap, err := a.pCollection.countByAppIds(appIds)
	if err != nil {
		return nil, err
	}

	userSimpleMap := collection.ToMap(users, func(u *entity.User) (int64, *entity.UserSimple) {
		return u.Id, &entity.UserSimple{
			Id:       u.Id,
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
		}
	})

	return collection.Map(apps, func(a *entity.MiniAppBaseInfo) *entity.MiniAppListDto {
		dto := &entity.MiniAppListDto{
			MiniAppBaseInfo: *a,
			Results:         make([]*entity.MiniAppOutput, 0, 0),
		}

		if outputs, ok := outputs[a.Uuid]; ok {
			dto.Results = outputs
		}

		if createdBy, ok := userSimpleMap[a.CreatedBy]; ok {
			dto.CreatedBy = *createdBy
		}

		if statistic, ok := statisticMap[a.Uuid]; ok {
			dto.MiniAppStatisticInfo = statistic.MiniAppStatisticInfo
			dto.SoldPoints = int64(statistic.RunTimes * a.Price)
		}

		if like, ok := likeMap[a.Uuid]; ok {
			dto.LikeTimes = int(like)
		}

		if comment, ok := commentMap[a.Uuid]; ok {
			dto.CommentTimes = int(comment)
		}

		if collect, ok := collectMap[a.Uuid]; ok {
			dto.CollectTimes = int(collect)
		}

		return dto
	}), nil
}

func (a appSvc) fixAppPrice(app *entity.MiniApp) error {
	points, err := a.points.GetStrategyPoints(entity.StrategyArgUsingApp{
		Form: *app.Form,
	})

	if err != nil {
		return err
	}

	app.Price = utils.Abs(points)
	return nil
}

func (a appSvc) createApp(app *entity.MiniApp) error {
	return a.Transaction(func(session xorm.Interface) error {
		if app.DuplicateFrom != "" {
			_, err := a.points.ApplyPoints(app.CreatedBy, entity.StrategyArgDuplicatingApp{})
			if err != nil {
				return err
			}
		}

		err := a.pMiniApp.createApp(app)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		return a.Sender.Send(&entity.AppCreatedEvent{
			AppId:         app.Uuid,
			CreatedBy:     app.CreatedBy,
			DuplicateFrom: app.DuplicateFrom,
			Time:          app.CreatedAt,
		})
	})
}

func (a appSvc) CountUserCreatedApps(userId int64) (int64, error) {
	return a.pMiniApp.countAppByUserId(userId)
}

func (a appSvc) CountUserRanApps(userId int64) (int64, error) {
	return a.pMiniApp.countOutputsAppIdByUserId(userId)
}
