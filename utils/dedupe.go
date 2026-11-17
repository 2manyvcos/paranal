package utils

// DedupeStart deduplicates a slice by comparing item keys, keeping the first occurrence.
func DedupeStart[T any, K comparable](items []T, key func(T) K) []T {
	seen := make(map[K]struct{}, len(items))
	result := make([]T, 0, len(items))

	for _, item := range items {
		k := key(item)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, item)
	}

	return result
}

// DedupeEnd deduplicates a slice by comparing item keys, keeping the last occurrence.
func DedupeEnd[T any, K comparable](items []T, key func(T) K) []T {
	seen := make(map[K]struct{}, len(items))
	result := make([]T, 0, len(items))

	for i := len(items) - 1; i >= 0; i-- {
		k := key(items[i])
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, items[i])
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}
