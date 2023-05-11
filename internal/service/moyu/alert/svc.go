package alert

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"time"
)

const TableName = "alert"

type svs struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewSvc() gone.Goner {
	return &svs{}
}

func (s *svs) GetAvailableAlert() (*entity.Alert, error) {
	var a entity.Alert
	exists, err := s.Table(TableName).Where("target_time < ?", time.Now()).Desc("id").Limit(1).Get(&a)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return &a, nil
}
