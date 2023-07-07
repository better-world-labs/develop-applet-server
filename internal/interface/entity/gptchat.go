package entity

import "time"

type GPTRole string

const (
	GPTRoleAssistant = "assistant"
	GPTRoleUser      = "user"
)

type LikeState uint8

const (
	LikeStateLike   = 1
	LikeStateNormal = 0
	LikeStateHate   = -1
)

type GptChatMessage struct {
	Id          int64     `json:"id"`
	MessageId   string    `json:"messageId"`
	Role        GPTRole   `json:"role"`
	UserId      int64     `json:"userId"`
	Content     string    `json:"content"`
	IsGptIgnore bool      `json:"-"`
	CreatedAt   time.Time `json:"createdAt"`
	Like        LikeState `xorm:"is_like" json:"like"`
}

func (m GptChatMessage) Cursor() int64 {
	return m.Id
}
