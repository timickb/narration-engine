package utils

func SliceToMap[T any, K comparable](slice []T, fn func(T) (K, T)) map[K]T {
	result := make(map[K]T)
	for _, item := range slice {
		key, value := fn(item)
		result[key] = value
	}
	return result
}

func MapSlice[T any, V any](slice []T, fn func(T) V) []V {
	result := make([]V, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}
