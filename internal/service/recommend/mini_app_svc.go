package recommend

import (
	"fmt"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type miniAppSvc struct {
	gone.Goner

	xorm.Engine `gone:"gone-xorm"`
	Sender      emitter.Sender   `gone:"gone-emitter"`
	User        service.IUser    `gone:"*"`
	MiniApp     service.IMiniApp `gone:"*"`
}

const (
	TableNameMiniAppRecommend = "mini_app_recommend"
)

//go:gone
func NewSvc() gone.Goner {
	return &miniAppSvc{}
}

func (s miniAppSvc) RecommendApp(recommend entity.MiniAppRecommend) error {
	has, err := s.MiniApp.CheckAppExists(recommend.AppId)
	if err != nil {
		return err
	}

	if !has {
		return gin.NewParameterError("app not found")
	}

	return s.Sender.Send(&entity.MiniAppRecommendEvent{MiniAppRecommend: recommend})
}

func (s miniAppSvc) compareAndSetRecommend(excepted, actual entity.MiniAppRecommend) (bool, error) {
	row, err := s.Exec(fmt.Sprintf("update %s set `recommend` = ?, updated_at = ? where app_id = ? and `recommend` = ? and updated_at = ?", TableNameMiniAppRecommend),
		actual.Recommend, actual.UpdatedAt, actual.AppId, excepted.Recommend, excepted.UpdatedAt)
	if err != nil {
		return false, err
	}

	affected, err := row.RowsAffected()
	return affected > 0, err
}

func (s miniAppSvc) create(recommend entity.MiniAppRecommend) error {
	_, err := s.Insert(recommend)
	return err
}

func (s miniAppSvc) DoRecommendApp(recommend entity.MiniAppRecommend) error {
	return s.Transaction(func(session xorm.Interface) error {
		rows, err := session.Exec(fmt.Sprintf("insert %s (app_id, `recommend`, created_by, updated_at) values (?, ?, ?, ?) on duplicate key update"+
			" `recommend` = IF(`recommend` != ? and ? > updated_at, ?, `recommend`), updated_at = IF(`recommend` != ? and ? > updated_at, ?, updated_at)", TableNameMiniAppRecommend),
			recommend.AppId, recommend.Recommend, recommend.CreatedBy, recommend.Recommend, recommend.UpdatedAt, recommend.UpdatedAt, recommend.Recommend, recommend.Recommend, recommend.UpdatedAt, recommend.UpdatedAt)

		affected, _ := rows.RowsAffected()
		if affected > 0 {
			return s.Sender.Send(&entity.MiniAppRecommendChangedEvent{
				MiniAppRecommend: recommend,
			})
		}

		return err
	})
}

func (s miniAppSvc) GetAppRecommend(appId string, userId int64) (entity.MiniAppRecommend, error) {
	var res entity.MiniAppRecommend
	_, err := s.Where("app_id = ? and user_id = ?", appId, userId).Get(&res)
	return res, err
}

func (s miniAppSvc) GetAppRecommendCountMap(appIds []string) (map[string]int64, error) {
	var res []*entity.MiniAppCount

	err := s.Table(entity.MiniAppRecommend{}).Select("app_id, count(1) count").In("app_id", appIds).And("`recommend` = 1").GroupBy("app_id").Find(&res)
	if err != nil {
		return nil, err
	}

	return collection.ToMap(res, func(c *entity.MiniAppCount) (string, int64) {
		return c.AppId, c.Count
	}), nil
}

func (s miniAppSvc) listAppComments(appId string) ([]*entity.MiniAppComment, error) {
	var res []*entity.MiniAppComment
	return res, s.Where("app_id = ?", appId).Desc("id").Find(&res)
}

func (s miniAppSvc) IsAppsRecommended(appIds []string, userId int64) (map[string]bool, error) {
	states, err := s.mapRecommendStates(appIds, userId)
	if err != nil {
		return nil, err
	}

	for _, appId := range appIds {
		if _, ok := states[appId]; !ok {
			states[appId] = false
		}
	}

	return states, nil
}

func (s miniAppSvc) mapRecommendStates(appIds []string, userId int64) (map[string]bool, error) {
	var res []*entity.UserRecommendState
	if err := s.Table(entity.MiniAppRecommend{}).In("app_id", appIds).And("created_by = ?", userId).Find(&res); err != nil {
		return nil, err
	}

	return collection.ToMap(res, func(state *entity.UserRecommendState) (string, bool) {
		return state.AppId, state.Recommend
	}), nil
}
