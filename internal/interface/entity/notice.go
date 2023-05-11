package entity

import "time"

const (
	NoticeTypeMention   NoticeType = "mention"
	NoticeTypeReference NoticeType = "reference"
	NoticeTypeApproval  NoticeType = "approval"
)

type (
	NoticeType string

	Notice struct {
		Id         int64      `json:"id"`
		UserId     int64      `json:"userId"`
		BusinessId int64      `json:"businessId"`
		Type       NoticeType `json:"type"`
		CreatedAt  time.Time  `json:"createdAt"`
		Read       bool       `xorm:"is_read" json:"read"`
	}
)
