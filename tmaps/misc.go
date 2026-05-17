package tmaps

// Ensure returns the value for key, inserting defaultValue if missing.
//
// Behavior details:
// - If the key exists, its value is returned and the map is unchanged.
// - If the key is missing and the map is non-nil, defaultValue is stored and returned.
// - If the map is nil, defaultValue is returned and no insertion happens.
//
// Example:
//
//	m := map[string]int{}
//	v := Ensure(m, "a", 10)
//	// v == 10, m["a"] == 10
func Ensure[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}

	if m != nil {
		m[key] = defaultValue
	}

	return defaultValue
}
