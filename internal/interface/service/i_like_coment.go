package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type ILikeCommentMiniApp interface {
	LikeApp(entity.MiniAppLike) error
	CountUserAppsLikes(userId int64) (int64, error)
	LikeAppOutput(output entity.MiniAppOutputLike) error
	DoLikeApp(entity.MiniAppLike) error
	DoLikeAppOutput(entity.MiniAppOutputLike) error
	GetAppLike(appId string, userId int64) (entity.MiniAppLike, error)
	ListUserOutputLikeState(outputIds []string, userId int64) ([]*entity.UserOutputLikeState, error)
	IsAppsLiked(appId []string, userId int64) (map[string]bool, error)
	GetAppLikeCountMap(appIds []string) (map[string]int64, error)
	GetAppCommentCountMap(appIds []string) (map[string]int64, error)
	CreateAppComment(comment entity.MiniAppComment) error
	ListAppComments(appId string) ([]*entity.MiniAppCommentDto, error)
}
