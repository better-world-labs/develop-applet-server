package entity

import "time"

type GPTRole string

const (
	GPTRoleAssistant = "assistant"
	GPTRoleUser      = "user"
)

type GptChatMessage struct {
	Id        int64     `json:"id"`
	Role      GPTRole   `json:"role"`
	UserId    int64     `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Like      bool      `xorm:"is_like" json:"like"`
}

func (m GptChatMessage) Cursor() int64 {
	return m.Id
}
