package jssdk

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/redis"
	"time"
)

const (
	Key              = "GLOBAL_JS_TICKET"
	RefreshAfterLast = time.Second * 60
)

type TicketCache struct {
	gone.Goner

	Redis redis.Cache `gone:"*"`
}

//go:gone
func NewTicketCache() gone.Goner {
	return &TicketCache{}
}

func (t *TicketCache) put(ticket JSTicket) error {
	return t.Redis.Put(Key, ticket, time.Duration(ticket.ExpiresIn)*time.Second-RefreshAfterLast)
}

func (t *TicketCache) get() (string, bool, error) {
	var ticket JSTicket
	if err := t.Redis.Get(Key, &ticket); err != nil {
		if err == redis.ErrNil {
			return "", false, nil
		}

		return "", false, err
	}

	return ticket.Ticket, true, nil
}
