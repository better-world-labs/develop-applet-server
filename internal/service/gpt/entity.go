package gpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`

	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name,omitempty"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	N                int                     `json:"n,omitempty"`
	Stream           bool                    `json:"stream,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int          `json:"logit_bias,omitempty"`
	User             string                  `json:"user,omitempty"`
}

type ChatCompletionStreamChoiceDelta struct {
	Content string `json:"content"`
}

type ChatCompletionStreamChoice struct {
	Index        int                             `json:"index"`
	Delta        ChatCompletionStreamChoiceDelta `json:"delta"`
	FinishReason string                          `json:"finish_reason"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

type ChatCompletionStreamReader struct {
	response *http.Response
	reader   *bufio.Reader
}

func (s *ChatCompletionStreamReader) Read() (resp ChatCompletionResponse, err error) {
	eventGroup, err := s.reader.ReadBytes('\n')
	if err != nil {
		return
	}

	eventGroup = bytes.TrimSpace(eventGroup)
	donePrefix := []byte("done: ")
	dataPrefix := []byte("data: ")
	if bytes.HasPrefix(eventGroup, donePrefix) {
		eventGroup = bytes.TrimLeft(eventGroup, string(donePrefix))
		var r struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}

		err = json.Unmarshal(eventGroup, &r)
		if err != nil {
			return
		}

		if r.Code != 0 {
			err = errors.New(r.Msg)
			return
		}

		err = io.EOF
		return
	}

	if bytes.HasPrefix(eventGroup, dataPrefix) {
		eventGroup = bytes.TrimLeft(eventGroup, string(dataPrefix))
		err = json.Unmarshal(eventGroup, &resp)
		if err != nil {
			return
		}
	}

	return
}

func (s *ChatCompletionStreamReader) Close() {
	s.response.Body.Close()
}
