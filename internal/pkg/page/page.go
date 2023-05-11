package page

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone/goner/gin"
)

const (
	DefaultPage = 1
	DefaultSize = 100
)

type (
	Query struct {
		Page int `form:"page"`
		Size int `form:"size"`
	}

	Result[T any] struct {
		Total int64 `json:"total"`
		List  []T   `json:"list"`
	}
)

func (q Query) LimitStart() int {
	return (q.Page - 1) * q.Size
}

func (q Query) LimitOffset() int {
	return q.Size
}

func ParseQuery(ctx *gin.Context) (Query, error) {
	query := Query{Page: DefaultPage, Size: DefaultSize}
	return query, ctx.ShouldBindQuery(&query)
}

func NewResult[T any](total int64, list []T) *Result[T] {
	return &Result[T]{Total: total, List: list}
}

func (r Result[T]) MarshalJSON() ([]byte, error) {
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
	b.WriteString(fmt.Sprintf(",\"total\":%d}", r.Total))
	fmt.Println(b.String())

	return b.Bytes(), nil
}
