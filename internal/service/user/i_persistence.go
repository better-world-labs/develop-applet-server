package user

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"time"
)

type (
	iUserPersistence interface {
		getUserByOpenId(openId string) (*entity.User, error)

		createUser(user *userDo) error

		getUserById(userId int64) (*entity.User, error)

		checkExists(userId int64) (bool, error)

		updateUserInfo(user *entity.User) error

		getUserInBatch(userIds []int64) ([]*entity.User, error)

		updateOnline(userId int64, online bool, t time.Time) error

		getOnlineUsers() (users []*entity.User, err error)

		updateConnectTime(userId int64) error

		updateTotalAccessDuration(userId, duration int64) error

		updateTotalBrowseDuration(userId, duration int64) error

		getUserOrderDesc(order string) ([]int64, error)

		getTotalBrowseDurationRankingList(top int) ([]*entity.User, error)

		updateAccessTimeByConnectTime(userIds []int64) error
	}

	iUserSettingsPersistence interface {
		getUserSettingsByUserId(userId int64) (*entity.UserSettings, error)

		getSimpleUserSettingsByUserId(userId int64, componentNames []entity.ComponentNameEnum) (*entity.UserSettings, error)

		getAllUserSettingsOrderByOffTime() ([]string, error)

		updateWorkOffTime(userId int64, offTime string) error

		updateBossKey(userId int64, bossKey string) error

		createUserSettings(settings entity.UserSettings) error

		updateUserSettings(settings entity.UserSettings, updates []entity.ComponentNameEnum) error
	}

	iUserAccessPersistence interface {
		updateEndTimeNull(userId int64) error

		insertWithoutEnd(userId int64) error
	}

	iUserBrowsePersistence interface {
		getLastRecordByUserId(userId int64, recordDate string) (*entity.UserBrowserRecord, error)

		update(record entity.UserBrowserRecord) error

		insert(record entity.UserBrowserRecord) error
	}
)

type userDo struct {
	*entity.User `xorm:"extends"`
}
