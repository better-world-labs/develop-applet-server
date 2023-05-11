package entity

import wsevent "gitlab.openviewtech.com/moyu-chat/ws-server/event"

type TriggerType int

const (
	TriggerTypeMention = 1
	TriggerTypeImage   = 2
)

type RobotConfig struct {
	Id                int64
	RobotId           string
	UserId            int64
	MessageReceiveUrl string
	Greeting          string
	CronExpression    string
	ContextLength     int
	Trigger           TriggerType
	Env               string // 只处理来自指定环境的消息
}

type RobotContext struct {
	ChannelId int64          `json:"channelId"`
	Robot     int64          `json:"robot"`
	Context   []*wsevent.Msg `json:"context"`
}
