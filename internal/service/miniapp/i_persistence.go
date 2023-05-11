package miniapp

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type iAIModel interface {
	list() ([]*entity.MiniAppAiModel, error)
	listCategory() ([]*entity.MiniAppAiModelCategory, error)
	checkExistsByName(name string) (bool, error)
}

type iCategoryPersistence interface {
	list() ([]*entity.MiniAppCategory, error)
	checkExists(id int64) (bool, error)
}

type iAppPersistence interface {
	listApp(category int64) ([]*entity.MiniAppBaseInfo, error)
	listAppByUuids(uuids []string) ([]*entity.MiniAppBaseInfo, error)
	listAppByUser(userId, category int64) ([]*entity.MiniAppBaseInfo, error)
	listOpenedOutputs(uuid string) ([]*entity.MiniAppOutput, error)
	getOutputById(outputId string) (*entity.MiniAppOutput, bool, error)
	getAppByUuid(uuid string) (*entity.MiniApp, bool, error)
	getLastNOutputByAppIds(appIds []string, last int) (map[string][]*entity.MiniAppOutput, error)
	checkAppExistsByUuid(uuid string) (bool, error)
	getAppById(id int64) (*entity.MiniApp, bool, error)
	updateApp(app *entity.MiniApp) error
	createApp(app *entity.MiniApp) error
	createOutput(output *entity.MiniAppOutput) error
	deleteApp(appId string) error
	deleteOutputsByAppId(appId string) error
	checkOutputExists(outputId string) (bool, error)
	checkAppExists(uuid string) (bool, error)
	checkOutputExistsByAppIdAndUser(appId string, userId int64) (bool, error)
	countAppByUserId(userId int64) (int64, error)
	countOutputsAppIdByUserId(userId int64) (int64, error)
}
