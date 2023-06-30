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
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Like      bool      `json:"like"`
}
