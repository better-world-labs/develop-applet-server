package gpt

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner"
	"testing"
)

func Priests(cemetery gone.Cemetery) error {
	_ = goner.BasePriest(cemetery)
	cemetery.Bury(NewCompletionGpt())
	return nil
}
func TestGpt(t *testing.T) {
	gone.Test(func(c *completionGpt) {
		stream, err := c.CreateChatCompletionStream(ChatCompletionRequest{
			Model:            "gpt-3.5-turbo",
			Temperature:      0.9,
			MaxTokens:        2048,
			TopP:             0.9,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
			Messages: []ChatCompletionMessage{
				{Role: "user", Content: "写一段 golang hello world"},
			},
		})

		if err != nil {
			t.Fatal(err)
			return
		}

		for {
			read, err := stream.Read()
			if err != nil {
				break
			}

			t.Logf("msg: %v\n", read)
		}
	}, Priests)
}
