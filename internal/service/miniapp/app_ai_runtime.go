package miniapp

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/gin"
)

type IAIRuntime interface {
	createCompletion(prompt string) (TextStreamChunkReader, error)

	model() string
}

type AIRuntime struct {
	gone.Goner

	ai    []IAIRuntime `gone:"*"`
	aiMap map[string]IAIRuntime
}

//go:gone
func NewAppFlowExecutor() gone.Goner {
	return &AIRuntime{
		aiMap: make(map[string]IAIRuntime),
	}
}

func (a AIRuntime) Start(cemetery gone.Cemetery) error {
	for _, client := range a.ai {
		a.aiMap[client.model()] = client
	}

	return nil
}

func (a AIRuntime) Stop(cemetery gone.Cemetery) error {
	return nil
}

func (a AIRuntime) createCompletion(aiModel, prompt string) (TextStreamChunkReader, error) {
	if client, ok := a.aiMap[aiModel]; ok {
		return client.createCompletion(prompt)
	}

	return nil, gin.NewParameterError("model not support yet")
}
