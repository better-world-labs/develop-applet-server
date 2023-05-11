package entity

import (
	"encoding/json"
	"time"
)

type MessageLike struct {
	Id        int64 `json:"id"`
	MessageId int64 `json:"messageId"`
	UserId    int64 `json:"userId"`
	IsLike    bool  `json:"isLike"`
}

type MessageLikeCount struct {
	Id   int64 `json:"id"`
	Like int   `json:"like"`
}

type MsgRecordSummary struct {
	UnreadCount int64 `json:"unreadCount"`
	LastReadId  int64 `json:"lastReadId"`
}

type MessageRecord struct {
	Id        int64           `json:"id"`
	SendId    string          `json:"sendId"`
	CreatedAt time.Time       `json:"createdAt"`
	SendAt    time.Time       `json:"sendAt"`
	UserId    int64           `json:"userId"`
	SeqId     int64           `json:"seqId"`
	ChannelId int64           `json:"channelId"`
	Content   json.RawMessage `json:"content"`
}
