package utils

import "gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"

func List(list any, err error) (any, error) {
	return entity.ListWrap{List: list}, err
}
