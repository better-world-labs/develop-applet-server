package utils

import (
	"errors"
	"github.com/gone-io/emitter"
	"log"
	"time"
)

type End <-chan struct{}

func ListenBatch[T emitter.DomainEvent](
	on emitter.OnEvent,
	bufferSize int,
	maxInterval time.Duration,
	end End,
	handler func([]T),
) {
	consume := BatchConsume[T](handler, bufferSize, maxInterval, end)
	on(func(t T) error {
		return consume(t)
	})
}

func BatchConsume[T any](deal func([]T), maxLength int, maxDuration time.Duration, end End) (wrapFn func(T) error) {
	var ch = make(chan T)
	var buf = make([]T, 0, maxLength)

	tick := time.NewTicker(maxDuration)

	go func() {
		var closed bool

		for {
			select {
			case t, ok := <-ch:
				if !ok {
					tick.Stop()
					if len(buf) > 0 {
						deal(buf)
					}

					return
				}

				buf = append(buf, t)
				if len(buf) >= maxLength {
					tick.Reset(maxDuration)
					deal(buf)
					buf = make([]T, 0, maxLength)
				}

			case <-tick.C:
				if len(buf) > 0 {
					deal(buf)
					buf = make([]T, 0, maxLength)
				}

			case <-end:
				if !closed {
					close(ch)
					closed = true
					log.Printf("consume close, bufferSize=%d\n", len(buf))
				}
			}
		}
	}()

	return func(t T) error {
		select {
		case <-end:
			return errors.New("consume closed")

		default:
			ch <- t
		}

		return nil
	}
}
