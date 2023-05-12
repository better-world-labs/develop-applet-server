package statistic

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

const TableNameMiniApp = "statistic_mini_app"
const TableNameMiniAppOutput = "statistic_mini_app_output"

type miniAppSvc struct {
	gone.Goner

	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewMiniAppSvc() gone.Goner {
	return &miniAppSvc{}
}

func (m miniAppSvc) GetByAppId(appId string) (entity.StatisticMiniApp, error) {
	s, err := m.GetAppMapByAppIds([]string{appId})
	if err != nil {
		return entity.StatisticMiniApp{}, nil
	}

	if statistic, ok := s[appId]; ok {
		return *statistic, nil
	}

	return entity.StatisticMiniApp{}, nil
}

func (m miniAppSvc) GetAppOutputMapByOutputIds(outputIds []string) (map[string]*entity.StatisticMiniAppOutput, error) {
	var result []*entity.StatisticMiniAppOutput
	err := m.In("output_id", outputIds).Find(&result)
	if err != nil {
		return nil, err
	}

	return collection.ToMap(result, func(s *entity.StatisticMiniAppOutput) (string, *entity.StatisticMiniAppOutput) {
		return s.OutputId, s
	}), nil
}

func (m miniAppSvc) GetAppMapByAppIds(appIds []string) (map[string]*entity.StatisticMiniApp, error) {
	var result []*entity.StatisticMiniApp
	err := m.In("app_id", appIds).Find(&result)
	if err != nil {
		return nil, err
	}

	return collection.ToMap(result, func(s *entity.StatisticMiniApp) (string, *entity.StatisticMiniApp) {
		return s.AppId, s
	}), nil
}

func (m miniAppSvc) ListAppIds() ([]string, error) {
	var result []*entity.StatisticMiniApp
	err := m.Cols("app_id").Find(&result)
	if err != nil {
		return nil, err
	}

	return collection.Map(result, func(s *entity.StatisticMiniApp) string {
		return s.AppId
	}), nil
}

func (m miniAppSvc) IncrementAppRuntimes(appId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, run_times) values (?, 1)  on duplicate key update "+
		"run_times =  run_times + 1, run_times_updated_at =  current_timestamp() ", TableNameMiniApp), appId)
	return err
}

func (m miniAppSvc) IncrementAppUseTimes(appId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, use_times) values (?, 1)  on duplicate key update "+
		"use_times =  use_times + 1, use_times_updated_at =  current_timestamp()", TableNameMiniApp), appId)
	return err
}

func (m miniAppSvc) IncrementAppLikeTimes(appId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, like_times) values (?, 1)  on duplicate key update "+
		"like_times =  like_times + 1, like_times_updated_at = current_timestamp() ", TableNameMiniApp), appId)
	return err
}

func (m miniAppSvc) DecrementAppLikeTimes(appId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, like_times) values (?, 1)  on duplicate key update "+
		"like_times =  IF(like_times = 0, 0, like_times - 1), like_times_updated_at =  current_timestamp()", TableNameMiniApp), appId)
	return err
}

func (m miniAppSvc) IncrementAppRecommendTimes(appId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, recommend_times) values (?, 1)  on duplicate key update "+
		"recommend_times =  recommend_times + 1, recommend_times_updated_at = current_timestamp() ", TableNameMiniApp), appId)
	return err
}

func (m miniAppSvc) DecrementAppRecommendTimes(appId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, recommend_times) values (?, 1)  on duplicate key update "+
		"recommend_times =  IF(recommend_times = 0, 0, recommend_times - 1), recommend_times_updated_at =  current_timestamp()", TableNameMiniApp), appId)
	return err
}

func (m miniAppSvc) IncrementAppCommentTimes(appId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, comment_times) values (?, 1)  on duplicate key update "+
		"comment_times =  comment_times + 1, comment_times_updated_at =  current_timestamp()", TableNameMiniApp), appId)
	return err
}

func (m miniAppSvc) IncrementAppViewTimes(appId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, view_times) values (?, 1)  on duplicate key update "+
		"view_times = view_times + 1, view_times_updated_at = current_timestamp() ", TableNameMiniApp), appId)
	return err
}

func (m miniAppSvc) OverrideAppRuntimes(appId string, t int64) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, run_times) values (?, ?)  on duplicate key update "+
		"run_times = ?, run_times_updated_at = current_timestamp() ", TableNameMiniApp), appId, t, t)
	return err
}

func (m miniAppSvc) OverrideAppUseTimes(appId string, t int64) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, use_times) values (?, ?)  on duplicate key update "+
		"use_times = ?, use_times_updated_at = current_timestamp() ", TableNameMiniApp), appId, t, t)
	return err
}

func (m miniAppSvc) OverrideAppLikeTimes(appId string, t int64) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, like_times) values (?, ?)  on duplicate key update "+
		"like_times = ?, like_times_updated_at = current_timestamp() ", TableNameMiniApp), appId, t, t)
	return err
}

func (m miniAppSvc) OverrideAppCommentTimes(appId string, t int64) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, comment_times) values (?, ?)  on duplicate key update "+
		"comment_times = ?, comment_times_updated_at = current_timestamp() ", TableNameMiniApp), appId, t, t)
	return err
}

func (m miniAppSvc) OverrideAppRecommendTimes(appId string, t int64) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, recommend_times) values (?, ?)  on duplicate key update "+
		"recommend_times = ?, recommend_times_updated_at = current_timestamp() ", TableNameMiniApp), appId, t, t)
	return err
}

func (m miniAppSvc) OverrideAppViewTimes(appId string, t int64) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (app_id, view_times) values (?, ?)  on duplicate key update "+
		"view_times = ?, view_times_updated_at = current_timestamp() ", TableNameMiniApp), appId, t, t)
	return err
}

func (m miniAppSvc) OverrideAppDegreeOfHeat(appId string, v float32) error {
	_, err := m.Exec(fmt.Sprintf("update %s set "+
		"degree_of_heat = ?, degree_of_heat_updated_at = current_timestamp() where app_id = ? ", TableNameMiniApp), v, appId)
	return err
}

func (m miniAppSvc) IncrementAppOutputCommentTimes(outputId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (output_id, comment_times) values (?, 1)  on duplicate key update "+
		"comment_times =  comment_times + 1, comment_times_updated_at =  current_timestamp()", TableNameMiniAppOutput), outputId)
	return err
}

func (m miniAppSvc) IncrementAppOutputLikeTimes(outputId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (output_id, like_times) values (?, 1)  on duplicate key update "+
		"like_times =  like_times + 1, like_times_updated_at = current_timestamp() ", TableNameMiniAppOutput), outputId)
	return err
}

func (m miniAppSvc) DecrementAppOutputLikeTimes(outputId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (output_id, like_times) values (?, 1)  on duplicate key update "+
		"like_times =  IF(like_times = 0, 0, like_times - 1), like_times_updated_at =  current_timestamp()", TableNameMiniAppOutput), outputId)
	return err
}

func (m miniAppSvc) IncrementAppOutputHateTimes(outputId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (output_id, hate_times) values (?, 1)  on duplicate key update "+
		"hate_times =  hate_times + 1, hate_times_updated_at = current_timestamp() ", TableNameMiniAppOutput), outputId)
	return err
}

func (m miniAppSvc) DecrementAppOutputHateTimes(outputId string) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (output_id, hate_times) values (?, 1)  on duplicate key update "+
		"hate_times =  IF(hate_times = 0, 0, hate_times - 1), hate_times_updated_at =  current_timestamp()", TableNameMiniAppOutput), outputId)
	return err
}

func (m miniAppSvc) OverrideAppOutputCommentTimes(outputId string, t int64) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (output_id, comment_times) values (?, ?)  on duplicate key update "+
		"comment_times = ?, comment_times_updated_at = current_timestamp() ", TableNameMiniAppOutput), outputId, t, t)
	return err
}

func (m miniAppSvc) OverrideAppOutputLikeTimes(outputId string, t int64) error {
	_, err := m.Exec(fmt.Sprintf("insert %s (output_id, like_times) values (?, ?)  on duplicate key update "+
		"like_times = ?, like_times_updated_at = current_timestamp() ", TableNameMiniAppOutput), outputId, t, t)
	return err
}
