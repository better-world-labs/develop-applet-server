package entity

// MiniAppLikeEvent 点赞时触发
type MiniAppLikeEvent struct {
	MiniAppLike
}

// MiniAppLikeChangedEvent 点赞发生变化时
type MiniAppLikeChangedEvent struct {
	MiniAppLike
}

// MiniAppOutputLikeEvent 点赞时触发
type MiniAppOutputLikeEvent struct {
	MiniAppOutputLike
}

// MiniAppOutputLikeChangedEvent 点赞发生变化时
type MiniAppOutputLikeChangedEvent struct {
	MiniAppOutputLike

	FromLikeState LikeType `json:"fromLikeState"`
}

func (m MiniAppOutputLikeChangedEvent) BecomeLike() bool {
	return m.FromLikeState != LikeValueLike && m.Like == LikeValueLike
}

func (m MiniAppOutputLikeChangedEvent) BecomeNormal() bool {
	return m.FromLikeState != LikeValueNormal && m.Like == LikeValueNormal
}

func (m MiniAppOutputLikeChangedEvent) BecomeHate() bool {
	return m.FromLikeState != LikeValueHate && m.Like == LikeValueHate
}

func (m MiniAppOutputLikeChangedEvent) FromLike() bool {
	return m.FromLikeState == LikeValueLike && m.Like != LikeValueLike
}

func (m MiniAppOutputLikeChangedEvent) FromHate() bool {
	return m.FromLikeState == LikeValueHate && m.Like != LikeValueHate
}

// MiniAppCommentedEvent 评论保存时触发
type MiniAppCommentedEvent struct {
	AppId  string `json:"appId"`
	UserId int64  `json:"userId"`
}
