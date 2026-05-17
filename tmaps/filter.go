package tmaps

// Filter returns a new map containing only the pairs that satisfy fn.
//
// Behavior details:
// - For a nil map, it returns nil.
// - For an empty map, it returns empty map.
// - The predicate is applied to each key-value pair.
// - Only pairs where fn returns true are included.
// - The output map has no guaranteed ordering.
//
// Example:
//
//	out := Filter(map[int]int{1: 10, 2: 20, 3: 30}, func(k int, v int) bool {
//		return k%2 == 1 && v >= 20
//	})
//	// out == map[int]int{3: 30}
func Filter[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	if m == nil {
		return nil
	}
	if len(m) == 0 {
		return m
	}

	r := make(map[K]V, len(m))
	for k, v := range m {
		if fn(k, v) {
			r[k] = v
		}
	}

	return r
}

// FilterFuncNotZero reports whether v is not the zero value for its type.
// Intended to be used with Filter.
//
// This helper is intended to be passed as the predicate to Filter when you want
// to drop entries whose values are unset/zero. For comparable types, the zero
// value is well-defined (0 for numbers, "" for strings, false for bool, nil for
// pointers/maps/slices/chans/functions/interfaces, and the zero-initialized
// value for structs).
//
// Behavior notes:
// - The key is ignored; only the value is inspected.
// - For pointer-like types (including maps/slices), non-nil values are kept.
// - For value types (numbers, bool, string, structs), non-zero values are kept.
//
// Example:
//
//	out := Filter(map[string]*int{
//		"a": nil,
//		"b": pointer.Of(10),
//	}, FilterFuncNotZero[string, *int])
//	// out == map[string]*int{"b": pointer.Of(10)}
//
//	out := Filter(map[string]int{
//		"a": 0,
//		"b": 10,
//	}, FilterFuncNotZero[string, int])
//	// out == map[string]int{"b": 10}
//
// See also: Filter.
func FilterFuncNotZero[K any, V comparable](_ K, v V) bool {
	var zero V
	return v != zero
}
