package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type IRecommendMiniApp interface {
	RecommendApp(entity.MiniAppRecommend) error
	DoRecommendApp(entity.MiniAppRecommend) error
	GetAppRecommend(appId string, userId int64) (entity.MiniAppRecommend, error)
	GetAppRecommendCountMap(appIds []string) (map[string]int64, error)
	IsAppsRecommended(appId []string, userId int64) (map[string]bool, error)
}
