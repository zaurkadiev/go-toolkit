package tmaps

// Iterate updates each value in m by applying fn to every key/value pair.
//
// It performs in-place transformation: for every entry in m, fn is called
// with the current key and value, and the return value replaces the existing
// value for that key.
//
// The iteration order is not specified because it follows Go's map iteration
// semantics. Do not rely on a particular ordering.
//
// The provided fn should be free of side effects that depend on iteration
// order, and should not add or delete keys in m. Mutating the map's keys while
// ranging over it can lead to unpredictable results.
//
// Example:
//
//	m := map[string]int{"a": 1, "b": 2}
//	maps.Iterate(m, func(_ string, v int) int {
//		return v * 10
//	})
//	// m is now map[string]int{"a": 10, "b": 20}
func Iterate[K comparable, V any](m map[K]V, fn func(K, V) V) {
	for k, v := range m {
		m[k] = fn(k, v)
	}
}
