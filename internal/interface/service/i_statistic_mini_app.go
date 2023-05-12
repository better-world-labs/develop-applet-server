package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type IStatisticMiniApp interface {
	ListAppIds() ([]string, error)

	GetByAppId(appId string) (entity.StatisticMiniApp, error)

	GetAppMapByAppIds(appIds []string) (map[string]*entity.StatisticMiniApp, error)

	GetAppOutputMapByOutputIds(outputIds []string) (map[string]*entity.StatisticMiniAppOutput, error)

	IncrementAppViewTimes(appId string) error

	IncrementAppRuntimes(appId string) error

	IncrementAppUseTimes(appId string) error

	IncrementAppCommentTimes(appId string) error

	IncrementAppLikeTimes(appId string) error

	DecrementAppLikeTimes(appId string) error

	IncrementAppRecommendTimes(appId string) error

	DecrementAppRecommendTimes(appId string) error

	OverrideAppViewTimes(appId string, t int64) error

	OverrideAppRuntimes(appId string, t int64) error

	OverrideAppUseTimes(appId string, t int64) error

	OverrideAppCommentTimes(appId string, t int64) error

	OverrideAppLikeTimes(appId string, t int64) error

	OverrideAppRecommendTimes(appId string, t int64) error

	OverrideAppDegreeOfHeat(appId string, v float32) error

	IncrementAppOutputCommentTimes(outputId string) error

	IncrementAppOutputLikeTimes(outputId string) error

	DecrementAppOutputLikeTimes(outputId string) error

	IncrementAppOutputHateTimes(outputId string) error

	DecrementAppOutputHateTimes(outputId string) error

	OverrideAppOutputCommentTimes(outputId string, t int64) error

	OverrideAppOutputLikeTimes(outputId string, t int64) error
}
