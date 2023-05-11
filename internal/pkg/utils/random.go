package utils

import "math/rand"

func RandomInt(min, max int) int {
	if min > max || min == 0 || max == 0 {
		return max
	}

	return rand.Intn(max-min+1) + min
}
