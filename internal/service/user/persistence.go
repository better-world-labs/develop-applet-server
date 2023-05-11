package user

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gone-io/emitter"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"strconv"
	"strings"
	"time"
)

const (
	TableName             = "user"
	SettingTableName      = "user_setting"
	AccessRecordTableName = "user_access_record"
	BrowseRecordTableName = "user_Browse_record"
)

//go:gone
func NewUserPersistence() gone.Goner {
	return &userPersistence{}
}

//go:gone
func NewUserSettingsPersistence() gone.Goner {
	return &userSettingsPersistence{}
}

//go:gone
func NewUserBrowseRecordPersistence() gone.Goner {
	return &userBrowseRecordPersistence{}
}

//go:gone
func NewUserAccessRecordPersistence() gone.Goner {
	return &userAccessRecordPersistence{}
}

type userPersistence struct {
	gone.Flag
	xorm.Engine    `gone:"gone-xorm"`
	emitter.Sender `gone:"gone-emitter"`
}

type userSettingsPersistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

type userBrowseRecordPersistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

type userAccessRecordPersistence struct {
	gone.Flag
	xorm.Engine `gone:"gone-xorm"`
}

func (p *userPersistence) getUserByOpenId(openId string) (*entity.User, error) {
	var user entity.User
	existed, err := p.Table(TableName).Where("wx_open_id = ?", openId).Get(&user)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if existed {
		return &user, nil
	}
	return nil, nil
}

func (p *userPersistence) createUser(user *userDo) error {
	return p.Transaction(func(session xorm.Interface) error {
		_, err := p.Table(TableName).Insert(user)
		if err != nil {
			return gin.ToError(err)
		}

		return err
	})
}

func (p *userPersistence) getUserById(userId int64) (*entity.User, error) {
	var user entity.User
	existed, err := p.Table(TableName).ID(userId).Get(&user)
	if err != nil {
		return nil, gin.ToError(err)
	}
	if existed {
		return &user, nil
	}
	return nil, nil
}

func (p *userPersistence) checkExists(userId int64) (bool, error) {
	exist, err := p.Table(TableName).Where("id = ?", userId).Exist()
	return exist, err
}

func (p *userPersistence) updateUserInfo(user *entity.User) error {
	_, err := p.Table(TableName).ID(user.Id).
		Cols("mobile", "nickname", "avatar", "sex", "invited_by", "from_app", "source", "login_at", "last_login_at").Update(user)
	if err != nil {
		return nil
	}

	return p.Send(&entity.UserInfoUpdated{UserId: user.Id})
}

func (p *userPersistence) getUsersByIds(ids []int64) (list []*entity.User, err error) {
	err = p.Table(TableName).In("id", ids).Find(&list)
	return
}

func (p *userPersistence) getUserInBatch(userIds []int64) ([]*entity.User, error) {
	res := make([]*entity.User, 0)
	var ids []int64
	for _, idPtr := range userIds {
		ids = append(ids, idPtr)
	}
	err := p.Table(TableName).Select("id, nickname, avatar, wx_open_id, wx_union_id, status, online, created_at, connect_time").In("id", ids).Desc("online").Find(&res)

	return res, err
}

func (p *userPersistence) updateOnline(userId int64, online bool, t time.Time) error {
	_, err := p.Exec("update user set online = ?, online_changed_at = ? where id = ? and online_changed_at < ?", online, t, userId, t)
	return err
}

func (p *userPersistence) getOnlineUsers() (users []*entity.User, err error) {
	err = p.Table(TableName).Where("online = 1").Find(&users)
	return
}

func (p *userPersistence) updateConnectTime(userId int64) error {
	currentTime := time.Now()
	user := entity.User{
		ConnectTime: &currentTime,
	}
	_, err := p.Table(TableName).Where("id = ?", userId).Cols("connect_time").Update(&user)
	return err
}

func (p *userPersistence) updateTotalAccessDuration(userId, duration int64) error {
	user := entity.User{
		TotalAccessDuration: duration,
	}
	_, err := p.Table(TableName).Where("id = ?", userId).Cols("total_access_duration").Update(&user)
	return err
}

func (p *userPersistence) updateTotalBrowseDuration(userId, duration int64) error {
	user := entity.User{
		TotalBrowseDuration: duration,
	}
	_, err := p.Table(TableName).Where("id = ?", userId).Cols("total_browse_duration").Update(&user)
	return err
}

func (p *userPersistence) getUserOrderDesc(order string) ([]int64, error) {
	var res []int64
	err := p.Table(TableName).Select(order).OrderBy(fmt.Sprintf("%s desc", order)).Find(&res)
	return res, err
}

func (p *userPersistence) getTotalBrowseDurationRankingList(top int) ([]*entity.User, error) {
	var res []*entity.User
	err := p.Table(TableName).OrderBy("total_browse_duration desc").Limit(top).Find(&res)
	//err := p.Table(TableName).OrderBy("total_access_duration desc").Limit(top).Find(&res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (p *userPersistence) updateAccessTimeByConnectTime(userIds []int64) error {
	var ids []string
	linq.From(userIds).SelectT(func(id int64) string {
		return strconv.FormatInt(id, 10)
	}).ToSlice(&ids)

	_, err := p.Exec(fmt.Sprintf("UPDATE user SET total_access_duration = total_access_duration + (UNIX_TIMESTAMP(NOW()) - UNIX_TIMESTAMP(connect_time)), connect_time = now() where id in (%s) and connect_time is not null", strings.Join(ids, ",")))
	return err
}

func (p *userSettingsPersistence) getUserSettingsByUserId(userId int64) (*entity.UserSettings, error) {
	var res entity.UserSettings
	exist, err := p.Table(SettingTableName).Where("user_id = ?", userId).Get(&res)
	if err != nil {
		return nil, err
	}

	if exist {
		return &res, err
	}

	return nil, nil
}

func (p *userSettingsPersistence) getSimpleUserSettingsByUserId(userId int64, componentNames []entity.ComponentNameEnum) (*entity.UserSettings, error) {
	var res entity.UserSettings

	// generate query cols
	var selectStat []string
	for _, name := range componentNames {
		col := name.ToQueryCol()
		if len(col) > 0 {
			selectStat = append(selectStat, name.ToQueryCol())
		}
	}

	if len(selectStat) == 0 {
		return nil, nil
	}

	query := strings.Join(selectStat, ",")

	exist, err := p.Table(SettingTableName).Select(query).Where("user_id = ? ", userId).Get(&res)
	if err != nil {
		return nil, err
	}

	if exist {
		return &res, err
	}

	return nil, nil
}

func (p *userSettingsPersistence) getAllUserSettingsOrderByOffTime() ([]string, error) {
	res := make([]string, 0)
	err := p.Table(SettingTableName).Select("end_off_time").OrderBy("end_off_time").Find(&res)

	return res, err
}

func (p *userSettingsPersistence) updateWorkOffTime(userId int64, offTime string) error {
	setting := entity.UserSettings{
		EndOffTime: offTime,
	}
	_, err := p.Table(SettingTableName).Where("user_id = ?", userId).Cols("end_off_time").Update(setting)
	return err
}

func (p *userSettingsPersistence) updateBossKey(userId int64, bossKey string) error {
	setting := entity.UserSettings{
		BossKey: bossKey,
	}
	_, err := p.Table(SettingTableName).Where("user_id = ?", userId).Cols("boss_key").Update(setting)
	return err
}

func (p *userSettingsPersistence) createUserSettings(settings entity.UserSettings) error {
	_, err := p.Table(SettingTableName).Insert(settings)
	return err
}

func (p *userSettingsPersistence) updateUserSettings(settings entity.UserSettings, updates []entity.ComponentNameEnum) error {
	// generate update cols
	var updateStat []string
	for _, name := range updates {
		col := name.ToQueryCol()
		if len(col) > 0 {
			updateStat = append(updateStat, name.ToQueryCol())
		}
	}

	if len(updateStat) == 0 {
		return nil
	}

	_, err := p.Table(SettingTableName).Where("user_id = ? ", settings.UserId).Cols(updateStat...).Update(settings)

	return err
}

func (p *userBrowseRecordPersistence) getLastRecordByUserId(userId int64, recordDate string) (*entity.UserBrowserRecord, error) {
	var res entity.UserBrowserRecord
	exist, err := p.Table(BrowseRecordTableName).Where("user_id = ? and browse_date = ?", userId, recordDate).Get(&res)
	if err != nil {
		return nil, err
	}

	if exist {
		return &res, err
	}

	return nil, err
}

func (p *userBrowseRecordPersistence) update(record entity.UserBrowserRecord) error {
	_, err := p.Table(BrowseRecordTableName).Where("id = ? ", record.Id).Cols("browse_duration", "last_browse_time").Update(&record)

	return err
}

func (p *userBrowseRecordPersistence) insert(record entity.UserBrowserRecord) error {
	_, err := p.Table(BrowseRecordTableName).Insert(record)
	return err
}

func (p *userAccessRecordPersistence) updateEndTimeNull(userId int64) error {
	var accessRecord entity.UserAccessRecord
	currentTime := time.Now()
	accessRecord.EndTime = &currentTime
	_, err := p.Table(AccessRecordTableName).Where("user_id = ? and end_time is null", userId).Cols("end_time").Update(&accessRecord)

	return err
}

func (p *userAccessRecordPersistence) insertWithoutEnd(userId int64) error {
	currentTime := time.Now()
	userAccess := entity.UserAccessRecord{
		UserId:    userId,
		StartTime: currentTime,
	}
	_, err := p.Table(AccessRecordTableName).Insert(userAccess)

	return err
}
