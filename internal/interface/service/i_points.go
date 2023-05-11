package service

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IPoints interface {
	// Add 如果 points 为负数，则是扣减。若余额不足，返回 businesserrors.ErrorPointsNotEnough
	Add(userId, points int64, _type string, description string) error
	GetUserPoints(userId int64) (int64, error)
	GetTodayIncomeByType(_type string, userId int64) (int64, error)
	PagePointFlow(query page.Query, userId int64) (page.Result[*entity.Points], error)

	RankingPointsToday(top int) ([]*entity.PointsRanking, error)
	RankingPointsTotal(top int) ([]*entity.PointsRanking, error)

	CheckExistsTodayByType(userId int64, _type string) (bool, error)
}
