package config

import (
	"encoding/json"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
)

type system struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewSystem() gone.Goner {
	return &system{}
}

func (s system) Put(key string, value any) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.Exec("insert `system_config`(`key`, `value`) values (?, ?) on duplicate key update `value` = ?", key, v, v)
	return err
}

func (s system) Get(Key string) (any, error) {
	var value entity.ConfigValue
	_, err := s.Table("system_config").
		Select("value").Where("`key` = ?", Key).
		Get(&value)

	return value.Value, err
}
