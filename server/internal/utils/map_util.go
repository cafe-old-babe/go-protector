package utils

func sliceToStringMap[T any, V any](
	slice []T, key func(elem T) string, val func(elem T) V) (toMap map[string]V) {
	toMap = make(map[string]V)
	if len(slice) <= 0 {
		return
	}

	// 遍历切片，将元素添加到 map 中
	for _, elem := range slice {
		if k := key(elem); len(k) > 0 {
			toMap[k] = val(elem)
		}
	}
	return
}

func SliceToUint64Map[T any, V any](
	slice []T, key func(elem T) uint64, val func(elem T) V) (toMap map[uint64]V) {
	toMap = make(map[uint64]V)
	if len(slice) <= 0 {
		return
	}

	// 遍历切片，将元素添加到 map 中
	for _, elem := range slice {
		if k := key(elem); k >= 0 {
			toMap[k] = val(elem)
		}
	}
	return
}
