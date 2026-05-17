package tmaps

// FindFirst returns the first key-value pair that satisfies fn and a found flag.
// If no pair matches, it returns zero values and false.
//
// Behavior details:
// - For a nil or empty map, it returns zero values and false.
// - Iteration stops at the first matching pair.
// - Because map iteration order is random, "first" is not deterministic.
//
// Example:
//
//	key, val, ok := FindFirst(map[int]int{1: 10, 2: 20}, func(k int, v int) bool {
//		return v == 20
//	})
//	// ok == true, key == 2, val == 20
func FindFirst[K comparable, V any](m map[K]V, fn func(K, V) bool) (K, V, bool) {
	for k, v := range m {
		if fn(k, v) {
			return k, v, true
		}
	}

	return *new(K), *new(V), false
}

// Contains reports whether any key-value pair satisfies fn.
//
// Behavior details:
// - For a nil or empty map, it returns false.
// - Iteration stops at the first matching pair.
//
// Example:
//
//	ok := Contains(map[int]int{1: 10, 2: 20}, func(k int, v int) bool {
//		return v > 15
//	})
//	// ok == true
func Contains[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if fn(k, v) {
			return true
		}
	}

	return false
}
