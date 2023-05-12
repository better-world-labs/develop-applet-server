package entity

// MiniAppRecommendEvent 推荐时触发
type MiniAppRecommendEvent struct {
	MiniAppRecommend
}

// MiniAppRecommendChangedEvent 推荐发生变化时
type MiniAppRecommendChangedEvent struct {
	MiniAppRecommend
}
