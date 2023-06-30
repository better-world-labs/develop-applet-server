package page

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gone-io/gone/goner/gin"
	"strconv"
)

const DefaultStreamSize = 10

type StreamCursor interface {
	Cursor() int64
}

type StreamQuery struct {
	nextCursor string
	size       int
}

func CreateStreamQuery(size int, cursor string) StreamQuery {
	return StreamQuery{
		nextCursor: cursor,
		size:       size,
	}
}

func (q StreamQuery) CursorIndicator() int64 {
	return unwrapCursor(q.nextCursor)
}

func (q StreamQuery) Size() int {
	return q.size
}

func (q *StreamQuery) BindQuery(ctx *gin.Context) error {
	var s struct {
		Cursor string `form:"cursor"`
		Size   int    `form:"size"`
	}

	err := ctx.ShouldBindQuery(&s)
	if err != nil {
		return nil
	}

	if s.Size == 0 {
		s.Size = DefaultStreamSize
	}
	*q = CreateStreamQuery(s.Size, s.Cursor)
	return nil
}

type StreamResult[T StreamCursor] struct {
	nextCursor string
	list       []T
}

func NewStreamResult[T StreamCursor](data []T) *StreamResult[T] {
	if data == nil {
		data = make([]T, 0, 0)
	}

	return &StreamResult[T]{
		list:       data,
		nextCursor: getCursor(DefaultStreamSize, data),
	}
}

func getCursor[T StreamCursor](size int, data []T) string {
	if len(data) < size {
		return ""
	}

	return wrapCursor(data[len(data)-1].Cursor())
}

func wrapCursor(cursorData int64) string {
	return base64.URLEncoding.EncodeToString([]byte(strconv.FormatInt(cursorData, 10)))
}

func unwrapCursor(cursor string) int64 {
	data, _ := base64.URLEncoding.DecodeString(cursor)
	cursorData, _ := strconv.ParseInt(string(data), 10, 64)
	return cursorData
}

func (q StreamResult[T]) GetList() []T {
	return q.list
}

func (q StreamResult[T]) GetNextCursor() string {
	return q.nextCursor
}

func (q StreamResult[T]) MarshalJSON() ([]byte, error) {
	s := struct {
		List       []T    `json:"list"`
		NextCursor string `json:"nextCursor"`
	}{
		List:       q.list,
		NextCursor: q.nextCursor,
	}

	return json.Marshal(s)
}
