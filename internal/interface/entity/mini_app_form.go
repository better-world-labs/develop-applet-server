package entity

import (
	"errors"
	"strings"
)

type MiniAppFormType string

const (
	MiniAppFormTypeText     MiniAppFormType = "text"
	MiniAppFormTypeSelect   MiniAppFormType = "select"
	MiniAppFormTypeCheckBox MiniAppFormType = "checkbox"
)

type MiniAppFormFields []MiniAppFormData

func (m MiniAppFormFields) MarshalJSON() ([]byte, error) {
	return EncodeMiniAppForms(m)
}

func (m *MiniAppFormFields) UnmarshalJSON(b []byte) error {
	forms, err := DecodeMiniAppForms(b)
	if err != nil {
		return err
	}

	*m = forms
	return nil
}

type MiniAppFormData interface {
	GetLabel() string
	GetType() MiniAppFormType
	GetId() string
	SetValue(value string) error
	GetValue() string
}

// MiniAppFormTextData 文本
type MiniAppFormTextData struct {
	id    string
	label string
	value string

	Placeholder string `json:"placeholder"`
}

func NewMiniAppFormTextData(id, label, placeholder string) *MiniAppFormTextData {
	return &MiniAppFormTextData{
		id:          id,
		label:       label,
		Placeholder: placeholder,
	}
}

func (d MiniAppFormTextData) GetLabel() string {
	return d.label
}

func (d *MiniAppFormTextData) GetId() string {
	return d.id
}

func (d MiniAppFormTextData) GetType() MiniAppFormType {
	return MiniAppFormTypeText
}

func (d *MiniAppFormTextData) GetValue() string {
	return d.value
}

func (d *MiniAppFormTextData) SetValue(value string) error {
	d.value = value
	return nil
}

// MiniAppFormSelectData 文本
type MiniAppFormSelectData struct {
	id    string
	label string
	value string

	Placeholder string `json:"placeholder"`
	Values      string `json:"values"`
}

func NewMiniAppFormSelectData(id, label, placeholder string, values string) *MiniAppFormSelectData {
	return &MiniAppFormSelectData{
		id:          id,
		label:       label,
		Values:      values,
		Placeholder: placeholder,
	}
}

func (d *MiniAppFormSelectData) GetLabel() string {
	return d.label
}

func (d *MiniAppFormSelectData) GetId() string {
	return d.id
}

func (d *MiniAppFormSelectData) GetType() MiniAppFormType {
	return MiniAppFormTypeSelect
}

func (d *MiniAppFormSelectData) GetValue() string {
	return d.value
}

func (d *MiniAppFormSelectData) SetValue(value string) error {
	for _, v := range strings.Split(d.Values, "\n") {
		if v == value {
			d.value = value
			return nil
		}
	}

	return errors.New("invalid value")
}
