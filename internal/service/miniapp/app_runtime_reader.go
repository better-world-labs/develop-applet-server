package miniapp

import (
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/service/gpt"
	"io"
)

type OutputReader interface {
	io.Closer
	Read() (*entity.MiniAppOutputCore, error)
}

type TextStreamChunkReader interface {
	io.Closer
	Read() (entity.TextStreamChunk, error)
}

type gptStreamTrunkReader struct {
	upstream *gpt.ChatCompletionStreamReader
}

func (g gptStreamTrunkReader) Close() error {
	g.upstream.Close()
	return nil
}

func (g gptStreamTrunkReader) Read() (chunk entity.TextStreamChunk, err error) {
	r, err := g.upstream.Read()
	if err != nil {
		return
	}

	chunk = entity.TextStreamChunk(r.Choices[0].Delta.Content)
	return
}
