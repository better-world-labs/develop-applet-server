package service

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type IRobot interface {
	SendMessage(env, robotId string, channelId int64, content []byte) error
	GetConfigByRobotId(robotId string) (*entity.RobotConfig, bool, error)
	CreateConfig(userId, channelId int64, messageReceiveUrl string) error
	ListConfigsByUserIds(userIds []int64) ([]*entity.RobotConfig, error)
	ListConfigsByTrigger(triggerType entity.TriggerType) ([]*entity.RobotConfig, error)
	ListConfigs() ([]*entity.RobotConfig, error)
}
