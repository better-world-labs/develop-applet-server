package points

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
)

type p struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewP() gone.Goner {
	return &p{}
}

func (c p) sumPointsByUserIdForUpdate(userId int64) (point int64, err error) {
	var res struct {
		Points int64
	}
	_ = c.Transaction(func(session xorm.Interface) error {
		_, err = session.Table("points").Select("sum(points) points").Where("user_id = ?", userId).ForUpdate().Get(&res)
		return err
	})

	point = res.Points
	return
}

func (c p) create(points *entity.Points) error {
	return c.Transaction(func(session xorm.Interface) error {
		_, err := session.Insert(points)
		return err
	})
}

func (c p) pageByUserId(query page.Query, userId int64) (result page.Result[*entity.Points], err error) {
	result.Total, err = c.Where("user_id = ?", userId).
		Desc("id").
		Limit(query.LimitOffset(), query.LimitStart()).
		FindAndCount(&result.List)

	return
}

func (c p) sumPointsByUserId(userId int64) (point int64, err error) {
	point, err = c.SumInt(&entity.Points{
		UserId: userId,
	}, "points")

	return
}

func (c p) sumTodayIncomesPointsByType(_type string, userId int64) (int64, error) {
	point, err := c.Where("created_at > current_date").SumInt(&entity.Points{
		UserId: userId,
		Type:   _type,
	}, "points")

	return point, err
}

func (c p) checkExistsTodayByType(userId int64, _type string) (bool, error) {
	return c.Where("user_id = ? and type = ? and created_at > current_date", userId, _type).
		Exist(&entity.Points{})
}

func (c p) rankingPoints(top int) ([]*entity.PointsRanking, error) {
	var res []*entity.PointsRanking

	session := c.Table(entity.Points{}).Select("user_id, sum(points) points").
		Where("points > 0").GroupBy("user_id").Desc("points")

	if top > 0 {
		session = session.Limit(top, 0)
	}

	return res, session.Find(&res)
}

func (c p) rankingPointsToday(top int) ([]*entity.PointsRanking, error) {
	var res []*entity.PointsRanking

	session := c.Table(entity.Points{}).Select("user_id, sum(points) points").
		Where("points > 0 and created_at > current_date()").GroupBy("user_id").Desc("points")

	if top > 0 {
		session = session.Limit(top, 0)
	}

	return res, session.Find(&res)
}
