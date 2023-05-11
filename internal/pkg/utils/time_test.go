package utils

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/test"
	"testing"
	"time"
)

func TestTodayRemainder(t *testing.T) {
	var useCases = []test.UseCase[string, time.Duration]{
		{
			Input:        "2023-05-08 08:28:00",
			ExceptOutput: 16*time.Hour - 28*time.Minute,
		},
		{
			Input:        "2023-05-08 23:59:59",
			ExceptOutput: 1 * time.Second,
		},
	}

	for _, useCase := range useCases {
		uc, err := time.Parse("2006-01-02 15:04:05", useCase.Input)
		if err != nil {
			t.Fatalf("%v\n", err)
			return
		}

		output := TodayRemainder(uc)
		if output != useCase.ExceptOutput {
			t.Fatalf("output mismatch: excepted=%d, actual=%d\n", useCase.ExceptOutput, output)
			return
		}
	}
}
