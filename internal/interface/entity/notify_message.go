package entity

import (
	"github.com/gone-io/gone/goner/gin"
	"time"
)

const (
	NotifyMessageTypeReward             NotifyMessageType = "notify-reward"             //空投奖励
	NotifyMessageTypeInvited            NotifyMessageType = "notify-invited"            //好友邀请
	NotifyMessageTypePointsRecharge     NotifyMessageType = "notify-points-recharge"    //积分充值
	NotifyMessageDuplicateApp           NotifyMessageType = "notify-duplicate-app"      //一键同款
	NotifyMessageAppUsed                NotifyMessageType = "notify-app-be-used"        //程序被使用
	NotifyMessageCreateAppPointsLimited NotifyMessageType = "create-app-points-limited" //创建程序积分限制
	NotifyMessageLoginReward            NotifyMessageType = "notify-login-reward"       //登录奖励
	NotifyMessageCreateApp              NotifyMessageType = "notify-create-app"         //创建程序
	NotifyMessageUseApp                 NotifyMessageType = "notify-use-app"            //使用程序
)

type (
	NotifyMessageType string

	NotifyMessage struct {
		NotifyMessageInfo `xorm:"extends"`

		Id        int64     `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
		IsRead    bool      `json:"read"`
	}

	NotifyMessageInfo struct {
		Type        NotifyMessageType `json:"type"`
		OperationId string            `json:"operationId"`
		Title       string            `json:"title"`
		Content     string            `json:"content"`
		UserId      int64             `json:"userId"`
	}

	NotifyMessageListFilter struct {
		IsRead *bool `form:"isRead"`
	}
)

func (m NotifyMessage) Cursor() int64 {
	return m.Id
}

func (m NotifyMessage) Validate() error {
	if m.UserId == 0 {
		return gin.NewParameterError("invalid userId")
	}

	if m.Title == "" {
		return gin.NewParameterError("invalid title")
	}

	return nil
}
