package entity

import "time"

// AppRunDoneEvent APP 运行完成
type AppRunDoneEvent struct {
	AppId  string            `json:"appId"`
	User   int64             `json:"user"`
	Param  MiniAppRunParam   `json:"param"`
	Output MiniAppOutputCore `json:"output"`
	Time   time.Time         `json:"time"`
}

// AppOutputCreatedEvent APP 运行完成
type AppOutputCreatedEvent struct {
	OutputId string
	AppId    string
	UserId   int64
}

// AppCreatedEvent APP 创建时触发
type AppCreatedEvent struct {
	AppId         string    `json:"appId"`
	CreatedBy     int64     `json:"createdBy"`
	DuplicateFrom string    `json:"duplicateFrom"`
	Time          time.Time `json:"time"`
}

// AppViewEvent APP 被浏览时触发
type AppViewEvent struct {
	AppId string    `json:"appId"`
	User  int64     `json:"user"`
	Time  time.Time `json:"time"`
}
