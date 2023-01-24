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
		result := make([]T, 0, len(arr)-1)
		result = append(result, arr[:index]...)
		result = append(result, arr[index+1:]...)
		return result
	}

	return arr
}

func ForEach[T any](arr []T, method func(item T)) {
	for _, item := range arr {
		method(item)
	}
}
