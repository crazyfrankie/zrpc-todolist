package slice

func Transform[A, B any](src []A, fn func(A) B) []B {
	if src == nil {
		return nil
	}

	dst := make([]B, 0, len(src))
	for _, a := range src {
		dst = append(dst, fn(a))
	}

	return dst
}

func TransformWithErrorCheck[A, B any](src []A, fn func(A) (B, error)) ([]B, error) {
	if src == nil {
		return nil, nil
	}

	dst := make([]B, 0, len(src))
	for _, a := range src {
		item, err := fn(a)
		if err != nil {
			return nil, err
		}
		dst = append(dst, item)
	}

	return dst, nil
}

func Unique[T comparable](src []T) []T {
	if src == nil {
		return nil
	}
	dst := make([]T, 0, len(src))
	m := make(map[T]struct{}, len(src))
	for _, s := range src {
		if _, ok := m[s]; ok {
			continue
		}
		dst = append(dst, s)
		m[s] = struct{}{}
	}

	return dst
}

func ToMap[E any, K comparable, V any](src []E, fn func(e E) (K, V)) map[K]V {
	if src == nil {
		return nil
	}

	dst := make(map[K]V, len(src))
	for _, e := range src {
		k, v := fn(e)
		dst[k] = v
	}

	return dst
}

func Batch[T any, V any](fn func(T) V, ts []T) []V {
	if ts == nil {
		return nil
	}
	res := make([]V, 0, len(ts))
	for i := range ts {
		res = append(res, fn(ts[i]))
	}
	return res
}
