package slice

func SubSlice[T comparable](a, b []T) []T {
	if len(b) == 0 {
		return a
	}

	set := make(map[T]struct{})
	for _, item := range b {
		set[item] = struct{}{}
	}

	seen := make(map[T]struct{})
	result := make([]T, 0, len(a))

	for _, item := range a {
		if _, inB := set[item]; inB {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		result = append(result, item)
		seen[item] = struct{}{}
	}

	return result
}
