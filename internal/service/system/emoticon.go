package system

import (
	"github.com/gone-io/gone/goner/schedule"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

func (s *svc) GetEmoticonList(group int, sort bool) ([]*entity.Emoticon, error) {
	return s.P.listByGroupId(group, sort)
}

const jobName = "stat-use-emoticon"

func (s *svc) Cron(run schedule.RunFuncOnceAt) {
	run(s.CronEmoticonStat, jobName, func() {
		err := s.P.updateRefStat()
		if err != nil {
			s.Errorf("stat-use-emoticon s.P.updateRefStat() err:%v", err)
		}
	})
}
