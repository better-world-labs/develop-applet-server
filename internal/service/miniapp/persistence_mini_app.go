package miniapp

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/gone/gone-lib/collection"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/page"
	"time"
)

type pMiniApp struct {
	gone.Goner
	xorm.Engine `gone:"gone-xorm"`
}

//go:gone
func NewPMiniApp() gone.Goner {
	return &pMiniApp{}
}

func (p pMiniApp) listApp(category int64) ([]*entity.MiniAppBaseInfo, error) {
	var arr []*entity.MiniAppBaseInfo
	session := p.Table(entity.MiniApp{}).
		Alias("m").
		Join("LEFT", []any{entity.StatisticMiniApp{}, "s"}, "m.uuid = s.app_id").
		Desc("top", "degree_of_heat", "m.id")
	if category > 0 {
		session = session.Where("category = ?", category)
	}

	return arr, session.Find(&arr)
}

func (p pMiniApp) listAppByUuids(uuids []string) ([]*entity.MiniAppBaseInfo, error) {
	var arr []*entity.MiniAppBaseInfo
	return arr, p.Table(entity.MiniApp{}).In("uuid", uuids).Find(&arr)
}

func (p pMiniApp) checkAppExists(uuid string) (bool, error) {
	return p.Table(entity.MiniApp{}).Where("uuid = ?", uuid).Exist()
}

func (p pMiniApp) checkOutputExists(outputId string) (bool, error) {
	return p.Table(entity.MiniAppOutput{}).Where("output_id = ?", outputId).Exist()
}

func (p pMiniApp) getOutputById(outputId string) (*entity.MiniAppOutput, bool, error) {
	var res entity.MiniAppOutput
	has, err := p.Where("output_id = ?", outputId).Get(&res)
	return &res, has, err
}

func (p pMiniApp) getLastNOutputByAppIds(appIds []string, last int) (map[string][]*entity.MiniAppOutput, error) {
	if len(appIds) == 0 {
		return map[string][]*entity.MiniAppOutput{}, nil
	}

	var results []*entity.MiniAppOutput
	err := p.Sqlx("select * from (select *, id % ? tag from mini_app_output where app_id in (?) order by id desc) a group by app_id, tag", last, appIds).Find(&results)
	if err != nil {
		return nil, err
	}

	return collection.GroupingBy(results, func(o *entity.MiniAppOutput) string {
		return o.AppId
	}, func(o *entity.MiniAppOutput) *entity.MiniAppOutput {
		return o
	}), nil
}

func (p pMiniApp) listAppByUser(userId, category int64) ([]*entity.MiniAppBaseInfo, error) {
	var arr []*entity.MiniAppBaseInfo
	session := p.Table(entity.MiniApp{}).Where("created_by = ?", userId)
	if category > 0 {
		session = session.And("category = ?", category)
	}

	return arr, session.Find(&arr)
}

func (p pMiniApp) listOpenedOutputs(appId string) ([]*entity.MiniAppOutput, error) {
	var arr []*entity.MiniAppOutput
	return arr, p.Where("app_id = ? and open = 1", appId).Desc("id").Find(&arr)
}

func (p pMiniApp) checkAppExistsByUuid(uuid string) (bool, error) {
	return p.Table(entity.MiniApp{}).Where("uuid = ?", uuid).Exist()
}

func (p pMiniApp) getAppByUuid(uuid string) (*entity.MiniApp, bool, error) {
	var app entity.MiniApp
	has, err := p.Table(entity.MiniApp{}).Where("uuid = ?", uuid).Get(&app)
	return &app, has, err
}

func (p pMiniApp) getAppById(id int64) (*entity.MiniApp, bool, error) {
	var app entity.MiniApp
	has, err := p.Table(entity.MiniApp{}).Where("id = ?", id).Get(&app)
	return &app, has, err
}

func (p pMiniApp) createApp(app *entity.MiniApp) error {
	return p.Transaction(func(session xorm.Interface) error {
		app.CreatedAt = time.Now()
		app.UpdatedAt = app.CreatedAt
		_, err := session.MustCols("status").Insert(app)
		return err
	})
}

func (p pMiniApp) updateApp(app *entity.MiniApp) error {
	app.UpdatedAt = time.Now()
	_, err := p.Where("uuid = ?", app.Uuid).Omit("duplicate_from").MustCols("status").Update(app)
	return err
}

func (p pMiniApp) createOutput(output *entity.MiniAppOutput) error {
	_, err := p.Insert(output)
	return err
}
func (p pMiniApp) deleteApp(appId string) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := session.Table(entity.MiniApp{}).Where("uuid = ?", appId).Delete()
		return err
	})
}

func (p pMiniApp) deleteOutputsByAppId(appId string) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := session.Table(entity.MiniAppOutput{}).Where("app_id = ?", appId).Delete()
		return err
	})

}

func (p pMiniApp) checkOutputExistsByAppIdAndUser(appId string, userId int64) (bool, error) {
	return p.Table(entity.MiniAppOutput{}).Where("app_id = ? and created_by = ?", appId, userId).Exist()
}

func (p pMiniApp) countAppByUserId(userId int64) (int64, error) {
	return p.Table(entity.MiniApp{}).Where("created_by = ?", userId).Count()
}

func (p pMiniApp) countOutputsByAppId(appId string) (int64, error) {
	return p.Table(entity.MiniAppOutput{}).Where("app_id = ?", appId).Count()
}

func (p pMiniApp) countUserAppRuntimes(userId int64) (int64, error) {
	return p.Where("app_id in (select uuid from mini_app where created_by = ?)", userId).Count(&entity.MiniAppOutput{})
}

func (p pMiniApp) countOutputsAppIdByUserId(userId int64) (int64, error) {
	return p.Table(entity.MiniAppOutput{}).
		Where("created_by = ?", userId).Distinct("app_id").Count()
}

func (p pMiniApp) pageAppsByUserId(query page.StreamQuery, userId int64) (*page.StreamResult[*entity.MiniAppBaseInfo], error) {
	var arr []*entity.MiniAppBaseInfo

	session := p.Table(entity.MiniApp{}).Where("created_by = ?", userId)
	if query.CursorIndicator() > 0 {
		session.Where("id < ?", query.CursorIndicator())
	}

	if err := session.Desc("id").Limit(query.Size(), 0).Find(&arr); err != nil {
		return nil, err
	}

	return page.NewStreamResult(arr), nil
}

func (p pMiniApp) pageOpenedOutputsByAppId(query page.StreamQuery, uuid string) (*page.StreamResult[*entity.MiniAppOutput], error) {
	var arr []*entity.MiniAppOutput

	session := p.Where("app_id = ? and open = 1", uuid)
	if query.CursorIndicator() > 0 {
		session.Where("id < ?", query.CursorIndicator())
	}

	if err := session.Desc("id").Limit(query.Size(), 0).Find(&arr); err != nil {
		return nil, err
	}

	return page.NewStreamResult(arr), nil
}
