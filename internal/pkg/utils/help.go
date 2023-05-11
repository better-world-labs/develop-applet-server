package utils

func ToInterfaceSlice[T any](t ...T) (list []interface{}) {
	for _, it := range t {
		list = append(list, it)
	}
	return
}
