package utils

import (
	"fmt"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"log"
	"testing"
	"time"
)

func TestBatchConsume(t *testing.T) {
	end := make(End)
	consume := BatchConsume[*entity.NoticeMessageCreatedEvent](func(evt []*entity.NoticeMessageCreatedEvent) {
		log.Printf("handleBatchMessage: size=%d", len(evt))
		for _, e := range evt {
			fmt.Println(e)
		}

	}, 5, 10*time.Second, end)

	for i := 0; i < 1004; i++ {
		err := consume(&entity.NoticeMessageCreatedEvent{Notice: entity.Notice{Id: int64(i)}})
		if err != nil {
			log.Printf("error: %v", err)
		}
	}

	err := consume(&entity.NoticeMessageCreatedEvent{Notice: entity.Notice{Id: int64(333333)}})
	if err != nil {
		log.Printf("error: %v", err)
	}
	err = consume(&entity.NoticeMessageCreatedEvent{Notice: entity.Notice{Id: int64(33333)}})
	if err != nil {
		log.Printf("error: %v", err)
	}
	time.Sleep(20 * time.Second)
}
