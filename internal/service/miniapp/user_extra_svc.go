package miniapp

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/xorm"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/service"
	"time"
)

type userExtra struct {
	gone.Goner
	logrus.Logger `gone:"gone-logger"`
	xorm.Engine   `gone:"gone-xorm"`

	user service.IUser `gone:"*"`
}

const TableNameUserExtra = "mini_app_user_extra"

//go:gone
func NewUserExtra() gone.Goner {
	return &userExtra{}
}

func (u userExtra) GetByUserId(userId int64) (entity.MiniAppUserExtra, bool, error) {
	var res entity.MiniAppUserExtra
	has, err := u.Where("user_id = ?", userId).Get(&res)
	return res, has, err
}

func (u userExtra) checkUserExists(userId int64) error {
	has, err := u.user.CheckUserExists(userId)
	if err != nil {
		return err
	}

	if !has {
		return gin.NewParameterError("user not found")
	}

	return nil
}

func (u userExtra) CompleteGuidance(userId int64) error {
	err := u.checkUserExists(userId)
	if err != nil {
		return err
	}

	_, err = u.Exec(fmt.Sprintf("insert %s (user_id, complete_guidance, created_at) "+
		"values (?, ?, ?) on duplicate key update complete_guidance = ?", TableNameUserExtra),
		userId, true, time.Now(), true,
	)
	return err
}
