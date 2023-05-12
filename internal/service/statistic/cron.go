package statistic

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/schedule"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

const (
	JobNameStatisticRevise     = "Job_Statistic_Revise"
	JobNameComputeDegreeOfHeat = "Job_ComputeDegreeOfHeat"
)

// cron 统计定时纠正
type cron struct {
	gone.Goner
	logrus.Logger `gone:"gone-logger"`

	miniApp     service.IStatisticMiniApp   `gone:"*"`
	miniAppSvc  service.IMiniApp            `gone:"*"`
	likeComment service.ILikeCommentMiniApp `gone:"*"`
	recommend   service.IRecommendMiniApp   `gone:"*"`

	jobSpecStatisticRevise     string `gone:"config,cron.statistic.revise"`
	jobSpecComputeDegreeOfHeat string `gone:"config,cron.statistic.compute-degree-of-heat"`
}

//go:gone
func NewCron() gone.Goner {
	return &cron{}
}

func (c cron) Cron(run schedule.RunFuncOnceAt) {
	run(c.jobSpecStatisticRevise, JobNameStatisticRevise, c.statisticRevise)
	run(c.jobSpecComputeDegreeOfHeat, JobNameComputeDegreeOfHeat, c.computeDegreeOfHeat)
}

func (c cron) statisticRevise() {
	c.Info("start statisticRevise:")
	appIds, err := c.miniApp.ListAppIds()
	if err != nil {
		c.Errorf("statisticRevise error: %v\n", err)
		return
	}

	for _, appId := range appIds {
		err := c.appCommentStatisticRevise(appId)
		if err != nil {
			c.Warnf("appCommentStatisticRevise error: appId=%s, %v\n", appId, err)
		}

		err = c.appRecommendStatisticRevise(appId)
		if err != nil {
			c.Warnf("appRecommendStatisticRevise error: appId=%s, %v\n", appId, err)
		}
	}
}

func (c cron) appRecommendStatisticRevise(appId string) error {
	countMap, err := c.recommend.GetAppRecommendCountMap([]string{appId})
	if err != nil {
		return err
	}

	if count, ok := countMap[appId]; ok {
		return c.miniApp.OverrideAppRecommendTimes(appId, count)
	}

	return nil
}

func (c cron) appCommentStatisticRevise(appId string) error {
	countMap, err := c.likeComment.GetAppCommentCountMap([]string{appId})
	if err != nil {
		return err
	}

	if count, ok := countMap[appId]; ok {
		return c.miniApp.OverrideAppCommentTimes(appId, count)
	}

	return nil
}

func (c cron) computeDegreeOfHeat() {
	appIds, err := c.miniApp.ListAppIds()
	if err != nil {
		c.Errorf("computeDegreeOfHeat error: %v\n", err)
		return
	}

	for _, appId := range appIds {
		err := c.computeAppDegreeOfHeat(appId)
		if err != nil {
			c.Errorf("computeAppDegreeOfHeat error: %v\n", err)
		}
	}
}

func (c cron) computeAppDegreeOfHeat(appId string) error {
	app, has, err := c.miniAppSvc.GetAppDetailByUuid(appId)
	if err != nil {
		return nil
	}
	if !has {
		return nil
	}

	outputs, err := c.miniAppSvc.ListOpenedAppOutputsByAppId(appId)
	if err != nil {
		return err
	}

	c.mergeStatistic(app, outputs)
	return c.miniApp.OverrideAppDegreeOfHeat(appId, float32(app.ViewTimes)*0.6+float32(app.RunTimes)*0.2)
}

func (c cron) mergeStatistic(app *entity.MiniAppDetailDto, outputs []*entity.MiniAppOutputDto) {
	for _, output := range outputs {
		app.CommentTimes += output.CommentTimes
		app.LikeTimes += output.LikeTimes
	}
}
