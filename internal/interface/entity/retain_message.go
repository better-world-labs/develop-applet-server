package entity

const (
	RetainMessageTypeFriendsFirstLogin = "friends-first-login"
)

type RetainMessageFriendsFirstLogin struct {
	Points  int           `json:"points"`
	Friends []*UserSimple `json:"friends"`
}

type RetainMessage struct {
	Id      int64                          `json:"id"`
	Type    string                         `json:"type"`
	UserId  int64                          `json:"userId"`
	Payload RetainMessageFriendsFirstLogin `xorm:"json" json:"payload"`
}
