package utils

func Map[T, F any](items []T, fn func(T) F) []F {
	result := make([]F, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}
