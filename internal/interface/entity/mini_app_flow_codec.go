package entity

import (
	"encoding/json"
	"errors"
)

func DecodeMiniAppFlowPrompts(b []byte) ([]AppFlowPrompt, error) {
	var promptRaw []*json.RawMessage
	prompts := make([]AppFlowPrompt, 0, len(promptRaw))

	err := json.Unmarshal(b, &promptRaw)
	if err != nil {
		return nil, err
	}

	for _, raw := range promptRaw {
		form, err := DecodeMiniAppFlowPrompt(*raw)
		if err != nil {
			return nil, err
		}

		prompts = append(prompts, form)
	}

	return prompts, nil
}

func DecodeMiniAppFlowPrompt(b []byte) (AppFlowPrompt, error) {
	var t struct {
		Type AppFlowPromptType `json:"type"`

		Properties *json.RawMessage `json:"properties"`
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}

	switch t.Type {
	case AppFlowPromptTypeText:
		var o = AppFlowPromptText{}
		return &o, json.Unmarshal(*t.Properties, &o)

	case AppFlowPromptTypeTag:
		var o = AppFlowPromptTag{}
		return &o, json.Unmarshal(*t.Properties, &o)

	default:
		return nil, errors.New("decode form error: invalid type")
	}
}

func EncodeMiniAppFlowPrompts(prompts []AppFlowPrompt) ([]byte, error) {
	arr := make([]json.RawMessage, 0, len(prompts))

	for _, f := range prompts {
		form, err := EncodeMiniAppFlowPrompt(f)
		if err != nil {
			return nil, err
		}

		arr = append(arr, form)
	}

	return json.Marshal(arr)
}

func EncodeMiniAppFlowPrompt(prompt AppFlowPrompt) ([]byte, error) {
	var o = struct {
		Type       AppFlowPromptType `json:"type"`
		Properties AppFlowPrompt     `json:"properties"`
	}{
		Type:       prompt.GetType(),
		Properties: prompt,
	}

	return json.Marshal(o)
}
