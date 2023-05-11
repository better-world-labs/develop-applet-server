package points

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

type iPersistence interface {
	create(points *entity.Points) error
	sumPointsByUserId(userId int64) (int64, error)
	rankingPoints(top int) ([]*entity.PointsRanking, error)
	rankingPointsToday(top int) ([]*entity.PointsRanking, error)
	sumTodayIncomesPointsByType(_type string, userId int64) (int64, error)
	sumPointsByUserIdForUpdate(userId int64) (point int64, err error)
	pageByUserId(query page.Query, userId int64) (page.Result[*entity.Points], error)
	checkExistsTodayByType(userId int64, _type string) (bool, error)
}
