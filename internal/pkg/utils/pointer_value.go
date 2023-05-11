package utils

func PointerValue[T any](i *T) T {
	var v T
	if i != nil {
		v = *i
	}

	return v
}
