package entity

type (

	// NoticeMessageCreatedEvent 创建一个通知时触发
	NoticeMessageCreatedEvent struct {
		Notice
	}

	// NoticeMessageRead 通知被标记已读时触发
	NoticeMessageRead struct {
		UserId int64   `json:"userId"`
		Ids    []int64 `json:"ids"`
	}
)
