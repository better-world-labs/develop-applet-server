package message

import (
	"encoding/json"
	"errors"
)

func EncodeContent(m Content) ([]byte, error) {
	return json.Marshal(m)
}

func DecodeContent(b []byte) (Content, error) {
	var header ContentHeader
	err := json.Unmarshal(b, &header)
	if err != nil {
		return nil, err
	}

	if provider, ok := messageContentFactory[header.Type]; ok {
		content := provider()

		return content, json.Unmarshal(b, &content)
	}

	return nil, errors.New("invalid message type")
}

func EncodeMessage(m *Message) ([]byte, error) {
	return json.Marshal(m)
}

func DecodeMessage(b []byte) (*Message, error) {
	var s struct {
		Header

		Content *json.RawMessage `json:"content"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	content, err := DecodeContent(*s.Content)
	if err != nil {
		return nil, err
	}

	return &Message{
		Header:  s.Header,
		Content: content,
	}, nil
}
