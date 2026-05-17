package tmaps

// NilToEmpty returns an empty non-nil map when m is nil.
// Non-nil maps are returned unchanged.
//
// Behavior details:
// - Nil input returns a new empty map.
// - Non-nil empty and populated maps are returned as-is.
// - Named map types are preserved.
//
// Example:
//
//	out := NilToEmpty(map[string]int(nil))
//	// out == map[string]int{}
func NilToEmpty[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return make(map[K]V)
	}

	return m
}
