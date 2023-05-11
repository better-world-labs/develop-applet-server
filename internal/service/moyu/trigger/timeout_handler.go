package trigger

import (
	"sync"
	"time"
)

type TimeoutHandler struct {
	t    *time.Timer
	d    time.Duration
	once sync.Once
}

func NewTimeoutHandler(d time.Duration) *TimeoutHandler {
	return &TimeoutHandler{
		d:    d,
		once: sync.Once{},
	}
}

func (t *TimeoutHandler) Handle(handler func()) {
	t.once.Do(func() {
		t.t = time.NewTimer(t.d)
	})

	t.t.Reset(t.d)

	go func() {
		<-t.t.C
		handler()
	}()
}
