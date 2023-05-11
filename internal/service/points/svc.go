package points

import (
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/xorm"
	"github.com/google/uuid"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	businesserrors "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/error"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"time"
)

type svc struct {
	gone.Goner

	xorm.Engine    `gone:"gone-xorm"`
	p              iPersistence  `gone:"*"`
	user           service.IUser `gone:"*"`
	emitter.Sender `gone:"gone-emitter"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svc{}
}

func (s svc) Add(userId, points int64, _type string, description string) error {
	if points == 0 {
		return nil
	}

	exists, err := s.user.CheckUserExists(userId)
	if err != nil {
		return err
	}

	if !exists {
		return gin.NewParameterError("user not found")
	}

	return s.Transaction(func(session xorm.Interface) error {
		remaining, err := s.p.sumPointsByUserIdForUpdate(userId)
		if err != nil {
			return nil
		}

		if remaining+points < 0 {
			return businesserrors.ErrorPointsNotEnough
		}

		points := entity.Points{
			Points:      points,
			Description: description,
			Type:        _type,
			UserId:      userId,
			CreatedAt:   time.Now(),
		}
		if err := s.p.create(&points); err != nil {
			return err
		}

		return s.Send(&entity.PointsChangeEvent{
			Points:      points,
			OperationId: uuid.NewString(),
		})
	})
}

func (s svc) CheckExistsTodayByType(userId int64, _type string) (bool, error) {
	return s.p.checkExistsTodayByType(userId, _type)
}

func (s svc) PagePointFlow(query page.Query, userId int64) (page.Result[*entity.Points], error) {
	return s.p.pageByUserId(query, userId)
}

func (s svc) GetUserPoints(userId int64) (int64, error) {
	return s.p.sumPointsByUserId(userId)
}

func (s svc) GetTodayIncomeByType(_type string, userId int64) (int64, error) {
	return s.p.sumTodayIncomesPointsByType(_type, userId)
}

func (s svc) RankingPointsToday(top int) ([]*entity.PointsRanking, error) {
	ranking, err := s.p.rankingPointsToday(top)
	if err != nil {
		return nil, err
	}

	return ranking, s.rankingPointsTranslate(ranking)
}

func (s svc) RankingPointsTotal(top int) ([]*entity.PointsRanking, error) {
	ranking, err := s.p.rankingPoints(top)
	if err != nil {
		return nil, err
	}

	return ranking, s.rankingPointsTranslate(ranking)
}

func (s svc) rankingPointsTranslate(data []*entity.PointsRanking) error {
	userIds := collection.Map(data, func(r *entity.PointsRanking) int64 {
		return r.UserId
	})

	users, err := s.user.GetUserSimpleInBatch(userIds)
	userMap := collection.ToMap(users, func(r *entity.UserSimple) (int64, *entity.UserSimple) {
		return r.Id, r
	})

	if err != nil {
		return err
	}

	for _, r := range data {
		if user, ok := userMap[r.UserId]; ok {
			r.User = *user
		}
	}

	return nil
}
