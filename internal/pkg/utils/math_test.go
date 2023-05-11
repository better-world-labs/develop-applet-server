package utils

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/pkg/test"
	"testing"
)

func TestAbsInt(t *testing.T) {
	var useCases = []test.UseCase[int, int]{
		{
			Input:        -9,
			ExceptOutput: 9,
		},
		{
			Input:        9,
			ExceptOutput: 9,
		},
	}

	for _, useCase := range useCases {
		if useCase.ExceptOutput != Abs(useCase.Input) {
			t.Fatalf("output mismatch: excepted=%d, actual=%d\n", useCase.ExceptOutput, useCase.Input)
			return
		}
	}
}

func TestAbsFloat(t *testing.T) {
	var useCases = []test.UseCase[float32, float32]{
		{
			Input:        -9.3,
			ExceptOutput: 9.3,
		},
		{
			Input:        9.9,
			ExceptOutput: 9.9,
		},
	}

	for _, useCase := range useCases {
		if useCase.ExceptOutput != Abs(useCase.Input) {
			t.Fatalf("output mismatch: excepted=%f, actual=%f\n", useCase.ExceptOutput, useCase.Input)
			return
		}
	}
}
