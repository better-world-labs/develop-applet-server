package robot

import (
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/schedule"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
)

type cron struct {
	gone.Goner
	logrus.Logger `gone:"gone-logger"`

	svc     service.IRobot   `gone:"*"`
	Channel service.IChannel `gone:"*"`
	Env     string           `gone:"config,server.env"`
}

//go:gone
func NewCron() gone.Goner {
	return &cron{}
}

func (c cron) Cron(run schedule.RunFuncOnceAt) {
	configs, err := c.svc.ListConfigs()
	if err != nil {
		panic(err)
	}

	for _, config := range configs {
		config := config
		jobName := fmt.Sprintf("robot_%s", config.RobotId)
		run(config.CronExpression, schedule.JobName(jobName), func() {
			c.Infof("RobotGreeting %v", config)

			msg := map[string]any{
				"type": "text",
				"text": config.Greeting,
			}

			b, err := json.Marshal(msg)
			if err != nil {
				c.Errorf("Marshal error: %v", err)
				return
			}

			channelIds, err := c.Channel.ListJoinedChannelIds(config.UserId)
			if err != nil {
				return
			}

			for _, channelId := range channelIds {
				err = c.svc.SendMessage(c.Env, config.RobotId, channelId, b)
				if err != nil {
					c.Errorf("SendMessage error: %v", err)
					return
				}
			}
		})
	}

}
