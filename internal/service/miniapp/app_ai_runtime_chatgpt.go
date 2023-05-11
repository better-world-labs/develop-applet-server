package miniapp

import (
	"github.com/gone-io/gone"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/service/gpt"
)

type AIRuntimeChatGpt struct {
	client gpt.ICompletionGpt `gone:"*"`

	gone.Goner
}

//go:gone
func NewChatGPTAppFlowExecutor() gone.Goner {
	return &AIRuntimeChatGpt{}
}

func (a AIRuntimeChatGpt) createCompletion(prompt string) (TextStreamChunkReader, error) {
	stream, err := a.client.CreateChatCompletionStream(gpt.ChatCompletionRequest{
		Model:            "gpt-3.5-turbo",
		Temperature:      0.9,
		MaxTokens:        2048,
		TopP:             0.9,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Messages: []gpt.ChatCompletionMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}

	return &gptStreamTrunkReader{
		upstream: stream,
	}, nil
}

func (a AIRuntimeChatGpt) model() string {
	return "chatgpt"
}
