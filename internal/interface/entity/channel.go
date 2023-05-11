package entity

import (
	"encoding/json"
	"time"
)

const (
	ChannelTypeNormal  ChannelType = 1
	ChannelTypePrivate ChannelType = 2
)

const (
	ChannelStatusValid   ChannelStatus = 1
	ChannelStatusInvalid ChannelStatus = 2
)

const (
	ChannelMemberStateNotJoin  ChannelMemberState = 0
	ChannelMemberStateApplying ChannelMemberState = 1
	ChannelMemberStateJoined   ChannelMemberState = 2
	ChannelMemberStateRemoved  ChannelMemberState = 3
)

type (
	MessageOffset struct {
		UserId            int64
		ChannelId         int64
		LastReadMessageId int64
	}

	ChannelType        uint
	ChannelStatus      uint
	ChannelMemberState uint

	Channel struct {
		Id        int64       `json:"id"`
		Name      string      `json:"name"`
		Type      ChannelType `json:"type"`
		Icon      string      `json:"icon"`
		GroupId   int64       `json:"groupId"`
		PlanetId  int64       `json:"planetId"`
		CreatedBy int64       `json:"createdBy"`
		CreatedAt time.Time   `json:"createdAt"`
		ExpiresAt *time.Time  `json:"expiresAt"`
		Sort      int64       `json:"sort"`
		State     int64       `json:"state"`
		Mute      bool        `json:"mute"`
		Notice    string      `json:"notice"`
	}

	ChannelGroup struct {
		Id        int64     `json:"id"`
		Name      string    `json:"name"`
		Icon      string    `json:"icon"`
		PlanetId  int64     `json:"planetId"`
		CreatedBy int64     `json:"createdBy"`
		CreatedAt time.Time `json:"createdAt"`
		Sort      int64     `json:"sort"`
	}

	ChannelMember struct {
		Id                int64              `json:"id"`
		ChannelId         int64              `json:"channelId"`
		UserId            int64              `json:"userId"`
		ApplyId           int64              `json:"applyId"`
		State             ChannelMemberState `json:"state"`
		LastReadMessageId int64              `json:"lastReadMessageId"`
		CreatedAt         time.Time          `json:"createdAt"`
	}
)

func (c *Channel) Status() ChannelStatus {
	if c.Type == ChannelTypeNormal {
		return ChannelStatusValid
	}

	if c.ExpiresAt == nil {
		return ChannelStatusInvalid
	}

	if expired := c.ExpiresAt.Before(time.Now()); expired {
		return ChannelStatusInvalid
	}

	return ChannelStatusValid
}

func (c *Channel) MarshalJSON() ([]byte, error) {
	var o = struct {
		Channel

		Status ChannelStatus `json:"status"`
	}{Channel: *c, Status: c.Status()}

	return json.Marshal(o)
}
