package utils

import (
	"strconv"
	"strings"
)

func SplitInt64(s string, sep string) ([]int64, error) {
	var res []int64
	split := strings.Split(s, sep)

	for _, i := range split {
		ite, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			return nil, err
		}

		res = append(res, ite)
	}

	return res, nil
}
