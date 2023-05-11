package robot

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

type iConfigPersistence interface {
	List() ([]*entity.RobotConfig, error)
	ListByTriggerType(triggerType entity.TriggerType) ([]*entity.RobotConfig, error)
	CreateConfig(config *entity.RobotConfig) error
	GetByRobotId(appId string) (*entity.RobotConfig, bool, error)
	ListByUserIds(userIds []int64) ([]*entity.RobotConfig, error)
}
