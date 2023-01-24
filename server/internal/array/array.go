package array

func Except[T any](arr []T, predicate func(item T) bool) []T {
	var index = -1
	for i, item := range arr {
		if predicate(item) {
			index = i
			break
		}
	}

	if index != -1 {
		if index+1 < len(arr) {
			return append(arr[:index], arr[index+1:]...)
		} else {
			return arr[:index]
		}
	}

	return arr
}

func ForEach[T any](arr []T, method func(item T)) {
	for _, item := range arr {
		method(item)
	}
}
