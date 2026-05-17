package tmaps

// MapVals transforms map values while preserving keys.
// Elements for which fn returns ok=false are skipped.
//
// Behavior details:
// - For a nil map, it returns nil.
// - For an empty but non-nil map, it returns an empty map.
// - The mapping is applied to each value; keys are preserved.
// - If fn returns ok=false, that key-value pair is omitted.
// - The output map has no guaranteed ordering.
//
// Example:
//
//	out := MapVals(map[int]int{
//		1: 10,
//		2: 20,
//	}, func(v int) (string, bool) {
//		return fmt.Sprintf("v-%d", v), v > 10
//	})
//	// out == map[int]string{2: "v-20"}
func MapVals[K comparable, V, R any](m map[K]V, fn func(V) (R, bool)) map[K]R {
	if m == nil {
		return nil
	}
	if len(m) == 0 {
		return map[K]R{}
	}

	r := make(map[K]R, len(m))
	for k, v := range m {
		if mapped, ok := fn(v); ok {
			r[k] = mapped
		}
	}

	return r
}

// MapValsErr transforms map values with an error-returning mapper.
// Elements for which fn returns ok=false are skipped; the first error stops processing.
//
// Behavior details:
// - For a nil map, it returns (nil, nil).
// - For an empty but non-nil map, it returns (empty map, nil).
// - The mapper is applied to each value; keys are preserved.
// - If fn returns ok=false, that key-value pair is omitted.
// - If fn returns a non-nil error, MapValsErr returns (nil, err) immediately.
//
// Example:
//
//	out, err := MapValsErr(map[int]string{
//		1: "10",
//		2: "x",
//	}, func(v string) (int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return 0, false, err
//		}
//		return n, n > 10, nil
//	})
//	// err != nil because of "x"
func MapValsErr[K comparable, V, R any](m map[K]V, fn func(V) (R, bool, error)) (map[K]R, error) {
	if m == nil {
		return nil, nil
	}
	if len(m) == 0 {
		return map[K]R{}, nil
	}

	r := make(map[K]R, len(m))
	for k, v := range m {
		mapped, ok, err := fn(v)
		if err != nil {
			return nil, err
		}
		if ok {
			r[k] = mapped
		}
	}

	return r, nil
}
