package miniapp

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"time"
)

type CollectGroupRes struct {
	AppId string
	Count int64
}
type AppId struct {
	AppId string
}

type miniAppCollectionP struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

const (
	TableNameCollection = "mini_app_collection"
)

//go:gone
func NewPersistence() gone.Goner {
	return &miniAppCollectionP{}
}

func (p miniAppCollectionP) createIfNotExists(appId string, userId int64) error {
	_, err := p.Exec(fmt.Sprintf("insert %s (app_id, user_id, created_at) values(?, ?, ?)"+
		"on duplicate key update created_at = ?", TableNameCollection), appId, userId, time.Now(), time.Now())
	return err
}

func (p miniAppCollectionP) deleteIfExists(appId string, userId int64) error {
	_, err := p.Exec(fmt.Sprintf("delete from %s where user_id = ? and app_id = ?", TableNameCollection), userId, appId)
	return err
}

func (p miniAppCollectionP) getAppIds(userId int64) ([]string, error) {
	var appIds []AppId

	err := p.Table(TableNameCollection).Where("user_id  = ?", userId).Desc("id").Find(&appIds)
	if err != nil {
		return nil, err
	}

	return collection.Map(appIds, func(appId AppId) string {
		return appId.AppId
	}), nil
}

func (p miniAppCollectionP) countByAppIds(appIds []string) (map[string]int64, error) {
	var res []CollectGroupRes
	err := p.Table(TableNameCollection).
		Select("app_id, count(1) count").
		In("app_id", appIds).GroupBy("app_id").
		Find(&res)
	if err != nil {
		return nil, err
	}

	return collection.ToMap(res, func(r CollectGroupRes) (string, int64) {
		return r.AppId, r.Count
	}), nil
}
