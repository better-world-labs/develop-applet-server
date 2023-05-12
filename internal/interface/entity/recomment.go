package entity

type MiniAppRecommend struct {
	AppId     string `json:"appId"`
	Recommend bool   `json:"recommend"`
	CreatedBy int64  `json:"createdBy"`
	UpdatedAt int64  `json:"updatedAt"`
}

type UserRecommendState struct {
	AppId     string
	UserId    int64
	Recommend bool
}
