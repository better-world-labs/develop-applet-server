package entity

// UpdateUnreadEvent used by moyu-server
type UpdateUnreadEvent struct {
	UpdateUnreadRange UpdateUnreadRange `json:"updateUnreadRange"`

	UserIds []int64 `json:"userIds"`
}

type UpdateUnreadRange string

const (
	All     UpdateUnreadRange = "ALL"
	Partial UpdateUnreadRange = "PARTIAL"
)

const (
	EventTypeViewApp = "app-viewed"
)

type UserEvent struct {
	UserId int64  `json:"userId"`
	Type   string `json:"type"`
	Args   []any  `json:"args"`
}
