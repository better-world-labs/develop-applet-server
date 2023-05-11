package entity

import (
	"bytes"
	"encoding/json"
)

type ListWrap struct {
	List any `json:"list"`
}

func (r ListWrap) MarshalJSON() ([]byte, error) {
	listJson, err := json.Marshal(&r.List)
	if err != nil {
		return nil, err
	}

	if string(listJson) == "null" {
		listJson = []byte("[]")
	}

	b := bytes.Buffer{}
	b.WriteString(`{"list":`)
	b.Write(listJson)
	b.WriteString(`}`)

	return b.Bytes(), nil
}
