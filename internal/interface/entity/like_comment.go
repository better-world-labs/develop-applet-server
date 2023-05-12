package entity

import "time"

type MiniAppLike struct {
	AppId     string `json:"appId"`
	Like      bool   `json:"like"`
	CreatedBy int64  `json:"createdBy"`
	UpdatedAt int64  `json:"updatedAt"`
}

type LikeType int8

const (
	LikeValueLike   LikeType = 1
	LikeValueHate   LikeType = -1
	LikeValueNormal LikeType = 0
)

type MiniAppOutputLike struct {
	UserOutputLikeState `xorm:"extends"`

	CreatedBy int64 `json:"createdBy"`
	UpdatedAt int64 `json:"updatedAt"`
}

type UserLikeState struct {
	AppId string `json:"appId"`
	Like  bool   `json:"like"`
}

type UserOutputLikeState struct {
	OutputId string   `json:"outputId"`
	Like     LikeType `json:"like" binding:"min=-1,max=1"`
}

type MiniAppCount struct {
	AppId string `json:"appId"`
	Count int64  `json:"count"`
}

type MiniAppComment struct {
	Id        int64     `json:"id"`
	AppId     string    `json:"appId"`
	Content   string    `json:"content" binding:"required"`
	CreatedBy int64     `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type MiniAppCommentDto struct {
	MiniAppComment

	CreatedBy UserSimple `json:"createdBy"`
}
