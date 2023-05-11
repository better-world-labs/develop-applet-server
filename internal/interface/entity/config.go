package entity

import "encoding/json"

type (
	ConfigValue struct {
		Value json.RawMessage `json:"value"`
	}
)
