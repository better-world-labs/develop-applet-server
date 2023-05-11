package miniapp

type iMiniAppCollectionPersistence interface {
	createIfNotExists(appId string, userId int64) error
	deleteIfExists(appId string, userId int64) error
	getAppIds(userId int64) ([]string, error)
	countByAppIds(appIds []string) (map[string]int64, error)
}
