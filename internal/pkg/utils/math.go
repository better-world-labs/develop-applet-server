package utils

type Number interface {
	int | int64 | int8 | int16 | int32 |
		float32 | float64
}

func Abs[T Number](n T) T {
	if n < 0 {
		n = -n
	}
	return n
}
