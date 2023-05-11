package entity

import (
	"encoding/json"
	"errors"
)

func DecodeMiniAppForms(b []byte) ([]MiniAppFormData, error) {
	var formsRaw []*json.RawMessage
	forms := make([]MiniAppFormData, 0, len(formsRaw))

	err := json.Unmarshal(b, &formsRaw)
	if err != nil {
		return nil, err
	}

	for _, raw := range formsRaw {
		form, err := DecodeMiniAppForm(*raw)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	return forms, nil
}

func DecodeMiniAppForm(b []byte) (MiniAppFormData, error) {
	var t struct {
		Id    string          `json:"id"`
		Label string          `json:"label"`
		Type  MiniAppFormType `json:"type"`

		Properties *json.RawMessage `json:"properties"`
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}

	switch t.Type {
	case MiniAppFormTypeText:
		var o = MiniAppFormTextData{
			id:    t.Id,
			label: t.Label,
		}
		return &o, json.Unmarshal(*t.Properties, &o)

	case MiniAppFormTypeSelect:
		var o = MiniAppFormSelectData{
			id:    t.Id,
			label: t.Label,
		}
		return &o, json.Unmarshal(*t.Properties, &o)

	default:
		return nil, errors.New("decode form error: invalid type")
	}
}

func EncodeMiniAppForms(forms []MiniAppFormData) ([]byte, error) {
	arr := make([]json.RawMessage, 0, len(forms))

	for _, f := range forms {
		form, err := EncodeMiniAppForm(f)
		if err != nil {
			return nil, err
		}

		arr = append(arr, form)
	}

	return json.Marshal(arr)
}

func EncodeMiniAppForm(form MiniAppFormData) ([]byte, error) {
	var o = struct {
		Id    string          `json:"id"`
		Label string          `json:"label"`
		Type  MiniAppFormType `json:"type"`

		Properties MiniAppFormData `json:"properties"`
	}{
		Id:         form.GetId(),
		Label:      form.GetLabel(),
		Type:       form.GetType(),
		Properties: form,
	}

	return json.Marshal(o)
}
