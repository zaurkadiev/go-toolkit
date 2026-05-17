package tmaps

// Map transforms map keys and values in one pass.
// Pairs for which fn returns ok=false are skipped.
//
// Behavior details:
// - For a nil map, it returns nil.
// - For an empty but non-nil map, it returns an empty map.
// - The mapper receives both key and value and returns a new key and value.
// - If fn returns ok=false, that pair is omitted from the result.
// - If multiple inputs map to the same output key, the last one wins.
// - The output map has no guaranteed ordering.
//
// Example:
//
//	out := Map(map[int]string{
//		1: "a",
//		2: "b",
//	}, func(k int, v string) (string, string, bool) {
//		return fmt.Sprintf("k-%d", k), v + v, k%2 == 0
//	})
//	// out == map[string]string{"k-2": "bb"}
func Map[K, K1 comparable, V, V1 any](m map[K]V, fn func(K, V) (K1, V1, bool)) map[K1]V1 {
	if m == nil {
		return nil
	}
	if len(m) == 0 {
		return map[K1]V1{}
	}

	r := make(map[K1]V1, len(m))
	for k, v := range m {
		if key, value, ok := fn(k, v); ok {
			r[key] = value
		}
	}

	return r
}

// MapErr transforms map keys and values with an error-returning mapper.
// Pairs for which fn returns ok=false are skipped; the first error stops processing.
//
// Behavior details:
// - For a nil map, it returns (nil, nil).
// - For an empty but non-nil map, it returns (empty map, nil).
// - The mapper receives key and value and returns a new key, value, ok flag, and error.
// - If fn returns ok=false, that pair is omitted.
// - If fn returns a non-nil error, MapErr returns (nil, err) immediately.
// - If multiple inputs map to the same output key, the last one wins.
//
// Example:
//
//	out, err := MapErr(map[int]string{
//		1: "10",
//		2: "x",
//	}, func(k int, v string) (string, int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return "", 0, false, err
//		}
//		return fmt.Sprintf("k-%d", k), n, true, nil
//	})
//	// err != nil because of "x"
func MapErr[K, K1 comparable, V, V1 any](m map[K]V, fn func(K, V) (K1, V1, bool, error)) (map[K1]V1, error) {
	if m == nil {
		return nil, nil
	}
	if len(m) == 0 {
		return map[K1]V1{}, nil
	}

	r := make(map[K1]V1, len(m))
	for k, v := range m {
		key, value, ok, err := fn(k, v)
		if err != nil {
			return nil, err
		}
		if ok {
			r[key] = value
		}
	}

	return r, nil
}

// ToSlice converts a map into a slice using a mapper function.
// Entries for which fn returns ok=false are skipped.
//
// Behavior details:
//   - For a nil map, it returns nil.
//   - For an empty but non-nil map, it returns an empty slice.
//   - The mapper receives both key and value and returns the output element and ok flag.
//   - If fn returns ok=false, that entry is omitted.
//   - The output slice order is not deterministic and depends on map iteration order.
//   - The returned slice may have a length greater than the number of accepted entries
//     if the mapper filters items; do not assume a tight length without checking.
//
// Example:
//
//	out := ToSlice(map[int]string{
//		1: "a",
//		2: "b",
//	}, func(k int, v string) (string, bool) {
//		if k%2 == 0 {
//			return v + v, true
//		}
//		return "", false
//	})
//	// out contains "bb" in an unspecified position.
func ToSlice[K comparable, V, V1 any](m map[K]V, fn func(K, V) (V1, bool)) []V1 {
	if m == nil {
		return nil
	}
	if len(m) == 0 {
		return []V1{}
	}

	r := make([]V1, 0, len(m))
	for k, v := range m {
		if v1, ok := fn(k, v); ok {
			r = append(r, v1)
		}
	}

	return r
}

// ToSliceErr converts a map into a slice using a mapper function that can fail.
// Entries for which fn returns ok=false are skipped; the first error stops processing.
//
// Behavior details:
//   - For a nil map, it returns (nil, nil).
//   - For an empty but non-nil map, it returns (empty slice, nil).
//   - The mapper receives both key and value and returns an output element, ok flag, and error.
//   - If fn returns ok=false, that entry is omitted.
//   - If fn returns a non-nil error, ToSliceErr returns (nil, err) immediately.
//   - The output slice order is not deterministic and depends on map iteration order.
//   - The returned slice may have a length greater than the number of accepted entries
//     if the mapper filters items; do not assume a tight length without checking.
//
// Example:
//
//	out, err := ToSliceErr(map[int]string{
//		1: "10",
//		2: "x",
//	}, func(k int, v string) (int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return 0, false, err
//		}
//		return n, k%2 == 0, nil
//	})
//	// err != nil because of "x"
func ToSliceErr[K comparable, V, V1 any](m map[K]V, fn func(K, V) (V1, bool, error)) ([]V1, error) {
	if m == nil {
		return nil, nil
	}
	if len(m) == 0 {
		return []V1{}, nil
	}

	r := make([]V1, 0, len(m))
	for k, v := range m {
		v1, ok, err := fn(k, v)
		if err != nil {
			return nil, err
		}
		if ok {
			r = append(r, v1)
		}
	}

	return r, nil
}
