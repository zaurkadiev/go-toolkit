package tslices

// Group partitions slice into buckets keyed by K and collects mapped values.
// Elements for which ok=false are skipped. Returns nil for empty input.
//
// Behavior details:
// - The mapper is called in order for each element.
// - Keys are created on demand; values are appended in the order encountered.
// - If no element passes the ok check, the result is an empty map.
//
// Example:
//
//	out := Group([]int{1, 2, 3, 4}, func(v int) (int, string, bool) {
//		return v % 2, strconv.Itoa(v), true
//	})
//	// out == map[int][]string{0: {"2", "4"}, 1: {"1", "3"}}
func Group[T any, K comparable, V any](slice []T, fn func(T) (K, V, bool)) map[K][]V {
	if len(slice) == 0 {
		return nil
	}

	r := make(map[K][]V, len(slice))
	for _, v := range slice {
		if key, value, ok := fn(v); ok {
			r[key] = append(r[key], value)
		}
	}

	return r
}

// GroupFuncAsIs is a helper that returns the same value as both key and value.
// Intended to be used with Group.
// It always returns ok=true and does not transform the input.
//
// Behavior details:
// - The same input is used as the group key and the grouped value.
// - Suitable when the element itself should serve as both key and value.
//
// Example:
//
//	out := Group([]string{"a", "b", "a"}, GroupFuncAsIs)
//	// out == map[string][]string{"a": {"a", "a"}, "b": {"b"}}
//
// See also: Group.
func GroupFuncAsIs[T any](t T) (T, T, bool) {
	return t, t, true
}

// GroupErr groups elements like Group but allows the mapper to return an error.
// It stops on the first error and returns it. Returns (nil, nil) for empty input.
//
// Behavior details:
// - The mapper is called in order for each element until an error occurs.
// - When err != nil, GroupErr stops and returns (nil, err).
// - When ok == false, the element is skipped.
// - Values are appended in encounter order per key.
//
// Example:
//
//	out, err := GroupErr([]string{"1", "x", "2"}, func(v string) (int, int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return 0, 0, false, err
//		}
//		return n % 2, n, true, nil
//	})
//	// err != nil (because of "x")
//
//	out, err = GroupErr([]string{"1", "2", "3"}, func(v string) (int, int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return 0, 0, false, err
//		}
//		return n % 2, n, true, nil
//	})
//	// out == map[int][]int{0: {2}, 1: {1, 3}}, err == nil
func GroupErr[T any, K comparable, V any](slice []T, fn func(T) (K, V, bool, error)) (map[K][]V, error) {
	if len(slice) == 0 {
		return nil, nil
	}

	r := make(map[K][]V, len(slice))
	for _, v := range slice {
		key, value, ok, err := fn(v)
		if err != nil {
			return nil, err
		}
		if ok {
			r[key] = append(r[key], value)
		}
	}

	return r, nil
}

// GroupErrFuncAsIs is the error-returning variant of GroupFuncAsIs.
// Intended to be used with GroupErr.
//
// It never returns an error and always returns ok=true.
//
// Example:
//
//	out, err := GroupErr([]string{"a", "b", "a"}, GroupErrFuncAsIs)
//	// out == map[string][]string{"a": {"a", "a"}, "b": {"b"}}, err == nil
//
// See also: GroupErr.
func GroupErrFuncAsIs[T any](t T) (T, T, bool, error) {
	return t, t, true, nil
}
