package domain

type (
	User struct {
		Id       int64  `json:"id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Online   bool   `json:"online"`
	}

	Member struct {
		User   User  `json:"user"`
		Status int64 `json:"status"`
		Role   int64 `json:"role"`
	}
)
