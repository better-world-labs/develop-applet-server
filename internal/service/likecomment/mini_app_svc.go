package likecomment

import (
	"fmt"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"time"
)

type miniAppSvc struct {
	gone.Goner

	xorm.Engine `gone:"gone-xorm"`
	Sender      emitter.Sender   `gone:"gone-emitter"`
	User        service.IUser    `gone:"*"`
	MiniApp     service.IMiniApp `gone:"*"`
}

const (
	TableNameMiniAppLike       = "mini_app_like"
	TableNameMiniAppOutputLike = "mini_app_output_like"
	TableNameMiniAppComment    = "mini_app_comment"
)

//go:gone
func NewSvc() gone.Goner {
	return &miniAppSvc{}
}

func (s miniAppSvc) LikeApp(like entity.MiniAppLike) error {
	has, err := s.MiniApp.CheckAppExists(like.AppId)
	if err != nil {
		return err
	}

	if !has {
		return gin.NewParameterError("app not found")
	}

	return s.Sender.Send(&entity.MiniAppLikeEvent{MiniAppLike: like})
}

func (s miniAppSvc) LikeAppOutput(like entity.MiniAppOutputLike) error {
	has, err := s.MiniApp.CheckOutputExists(like.OutputId)
	if err != nil {
		return err
	}

	if !has {
		return gin.NewParameterError("output not found")
	}

	return s.Sender.Send(&entity.MiniAppOutputLikeEvent{MiniAppOutputLike: like})
}

func (s miniAppSvc) GetAppLike(appId string, userId int64) (res entity.MiniAppLike, err error) {
	_, err = s.Where("created_by = ? and app_id = ?", userId, appId).Get(&res)
	return
}

func (s miniAppSvc) getByOutputId(outputId string) (entity.MiniAppOutputLike, bool, error) {
	var like entity.MiniAppOutputLike
	has, err := s.Where("output_id = ?", outputId).Get(&like)
	return like, has, err
}

func (s miniAppSvc) compareAndSetLike(excepted, actual entity.MiniAppOutputLike) (bool, error) {
	row, err := s.Exec(fmt.Sprintf("update %s set `like` = ?, updated_at = ? where output_id = ? and `like` = ? and updated_at = ?", TableNameMiniAppOutputLike),
		actual.Like, actual.UpdatedAt, actual.OutputId, excepted.Like, excepted.UpdatedAt)
	if err != nil {
		return false, err
	}

	affected, err := row.RowsAffected()
	return affected > 0, err
}

func (s miniAppSvc) create(like entity.MiniAppOutputLike) error {
	_, err := s.Insert(like)
	return err
}

func (s miniAppSvc) DoLikeAppOutput(like entity.MiniAppOutputLike) error {
	excepted, has, err := s.getByOutputId(like.OutputId)
	if err != nil {
		return err
	}

	if !has {
		err := s.create(like)
		if err != nil {
			return err
		}
	} else {
		if excepted.UpdatedAt >= like.UpdatedAt ||
			excepted.Like == like.Like {
			return nil
		}
		ok, err := s.compareAndSetLike(excepted, like)
		if err != nil {
			return err
		}

		if !ok {
			time.Sleep(200 * time.Millisecond)
			return s.DoLikeAppOutput(like)
		}
	}

	return s.Sender.Send(&entity.MiniAppOutputLikeChangedEvent{
		FromLikeState:     excepted.Like,
		MiniAppOutputLike: like,
	})
}

func (s miniAppSvc) DoLikeApp(like entity.MiniAppLike) error {
	return s.Transaction(func(session xorm.Interface) error {
		rows, err := session.Exec(fmt.Sprintf("insert %s (app_id, `like`, created_by, updated_at) values (?, ?, ?, ?) on duplicate key update"+
			" `like` = IF(`like` != ? and ? > updated_at, ?, `like`), updated_at = IF(`like` != ? and ? > updated_at, ?, updated_at)", TableNameMiniAppLike),
			like.AppId, like.Like, like.CreatedBy, like.Like, like.UpdatedAt, like.UpdatedAt, like.Like, like.Like, like.UpdatedAt, like.UpdatedAt)

		affected, _ := rows.RowsAffected()
		if affected > 0 {
			return s.Sender.Send(&entity.MiniAppLikeChangedEvent{
				MiniAppLike: like,
			})
		}

		return err
	})
}

func (s miniAppSvc) CreateAppComment(comment entity.MiniAppComment) error {
	return s.Transaction(func(session xorm.Interface) error {
		var t time.Time

		if comment.CreatedAt == t {
			comment.CreatedAt = time.Now()
		}

		_, err := s.Insert(comment)
		if err != nil {
			return err
		}

		return s.Sender.Send(&entity.MiniAppCommentedEvent{
			AppId: comment.AppId, UserId: comment.CreatedBy,
		})
	})
}

func (s miniAppSvc) ListAppComments(appId string) ([]*entity.MiniAppCommentDto, error) {
	comments, err := s.listAppComments(appId)
	if err != nil {
		return nil, err
	}

	users, err := s.User.GetUserSimpleInBatch(collection.Map(comments, func(c *entity.MiniAppComment) int64 {
		return c.CreatedBy
	}))
	if err != nil {
		return nil, err
	}

	userMap := collection.ToMap(users, func(u *entity.UserSimple) (int64, *entity.UserSimple) {
		return u.Id, u
	})

	return collection.Map(comments, func(c *entity.MiniAppComment) *entity.MiniAppCommentDto {
		dto := entity.MiniAppCommentDto{
			MiniAppComment: *c,
		}

		if u, ok := userMap[c.CreatedBy]; ok {
			dto.CreatedBy = *u
		}

		return &dto
	}), nil
}

func (s miniAppSvc) GetAppCommentCountMap(appIds []string) (map[string]int64, error) {
	var res []*entity.MiniAppCount

	err := s.Table(entity.MiniAppComment{}).Select("app_id, count(1) count").In("app_id", appIds).GroupBy("app_id").Find(&res)
	if err != nil {
		return nil, err
	}

	return collection.ToMap(res, func(c *entity.MiniAppCount) (string, int64) {
		return c.AppId, c.Count
	}), nil
}

func (s miniAppSvc) GetAppLikeCountMap(appIds []string) (map[string]int64, error) {
	var res []*entity.MiniAppCount

	err := s.Table(entity.MiniAppLike{}).Select("app_id, count(1) count").In("app_id", appIds).And("`like` = 1").GroupBy("app_id").Find(&res)
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

func (s miniAppSvc) ListUserOutputLikeState(outputIds []string, userId int64) ([]*entity.UserOutputLikeState, error) {
	states, err := s.mapOutputsLikeStates(outputIds, userId)
	return collection.Map(outputIds, func(id string) *entity.UserOutputLikeState {
		if state, ok := states[id]; ok {
			return state
		}

		return &entity.UserOutputLikeState{
			OutputId: id,
		}
	}), err
}

func (s miniAppSvc) listOutputsLikeStates(outputIds []string, userId int64) ([]*entity.UserOutputLikeState, error) {
	var res []*entity.UserOutputLikeState
	return res, s.Table(entity.MiniAppOutputLike{}).In("output_id", outputIds).And("created_by = ?", userId).Find(&res)
}

func (s miniAppSvc) mapOutputsLikeStates(outputIds []string, userId int64) (map[string]*entity.UserOutputLikeState, error) {
	states, err := s.listOutputsLikeStates(outputIds, userId)
	return collection.ToMap(states, func(in *entity.UserOutputLikeState) (string, *entity.UserOutputLikeState) {
		return in.OutputId, in
	}), err
}
