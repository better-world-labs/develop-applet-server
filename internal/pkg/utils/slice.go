package utils

func CutTail[T any](source []T, length int) []T {
	if len(source) < length {
		return source
	}

	return source[len(source)-length:]
}
