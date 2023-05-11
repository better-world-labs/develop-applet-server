package gpt

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/imroc/req"
	"net/http"
)

type completionGpt struct {
	gone.Goner

	host string `gone:"config,server.ai"`
}

//go:gone
func NewCompletionGpt() gone.Goner {
	return &completionGpt{}
}

func (c *completionGpt) CreateChatCompletionStream(request ChatCompletionRequest) (*ChatCompletionStreamReader, error) {
	resp, err := req.Post(fmt.Sprintf("%s/api/gpt/chat-completions", c.host), req.BodyJSON(request))
	if err != nil {
		return nil, err
	}

	if resp.Response().StatusCode != http.StatusOK {
		return nil, errors.New(resp.String())
	}

	return &ChatCompletionStreamReader{
		reader:   bufio.NewReader(resp.Response().Body),
		response: resp.Response(),
	}, nil
}
