package robot

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"github.com/google/uuid"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type configPersistence struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewConfigPersistence() gone.Goner {
	return &configPersistence{}
}

func (c configPersistence) ListByTriggerType(triggerType entity.TriggerType) ([]*entity.RobotConfig, error) {
	var res []*entity.RobotConfig
	return res, c.Where("`trigger` = ?", triggerType).Find(&res)
}

func (c configPersistence) CreateConfig(config *entity.RobotConfig) error {
	if config.RobotId == "" {
		config.RobotId = uuid.NewString()
	}

	_, err := c.Insert(config)
	return err
}

func (c configPersistence) GetByRobotId(appId string) (*entity.RobotConfig, bool, error) {
	var res entity.RobotConfig
	exists, err := c.Where("robot_id = ?", appId).Get(&res)
	return &res, exists, err
}

func (c configPersistence) List() ([]*entity.RobotConfig, error) {
	var res []*entity.RobotConfig
	return res, c.Find(&res)
}

func (c configPersistence) ListByUserIds(userIds []int64) ([]*entity.RobotConfig, error) {
	var res []*entity.RobotConfig
	return res, c.In("user_id", userIds).Find(&res)
}
