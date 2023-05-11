package limit

import (
	"fmt"
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/redis"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/utils"
	"time"
)

const (
	KeyAppCreateNotify = "limit-app-create-notify"
)

type marker struct {
	gone.Goner

	cache redis.Cache `gone:"*"`
}

//go:gone
func NewMarker() gone.Goner {
	return &marker{}
}

func (m marker) IsAppCreateNotifySent(userId int64) (bool, error) {
	var is int
	err := m.cache.Get(m.key(KeyAppCreateNotify, userId), &is)
	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return is == 1, err
}

func (m marker) MarkAppCreatePointsLimitNotifySent(userId int64) error {
	ttl := utils.TodayRemainder(time.Now())
	return m.cache.Put(m.key(KeyAppCreateNotify, userId), 1, ttl)
}

func (m marker) key(key string, userId int64) string {
	return fmt.Sprintf("%s-%d", key, userId)
}
