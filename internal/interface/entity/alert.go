package entity

import "time"

type (
	Alert struct {
		Id         int64     `json:"id"`
		Title      string    `json:"title"`
		HeadImage  string    `json:"headImage"`
		Content    string    `json:"content"`
		CreatedAt  time.Time `json:"-"`
		TargetTime time.Time `json:"-"`
	}
)

func (a Alert) Available() bool {
	return time.Now().After(a.TargetTime)
}
