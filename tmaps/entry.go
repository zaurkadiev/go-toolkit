package tmaps

import (
	"cmp"
	"slices"
)

// GetKeys returns the keys of the map as a slice.
//
// Behavior details:
// - For a nil or empty map, it returns an empty slice.
// - The order of keys is not guaranteed because Go maps are unordered.
// - The returned slice contains each key exactly once.
//
// Example:
//
//	keys := GetKeys(map[int]string{
//		1: "a",
//		2: "b",
//	})
//	// keys contains 1 and 2 (order not guaranteed)
func GetKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// GetVals returns the values of the map as a slice.
//
// Behavior details:
// - For a nil or empty map, it returns an empty slice.
// - The order of values is not guaranteed because Go maps are unordered.
// - The returned slice contains each value exactly once per map entry.
//
// Example:
//
//	values := GetVals(map[int]string{
//		1: "a",
//		2: "b",
//	})
//	// values contains "a" and "b" (order not guaranteed)
func GetVals[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

// GetKeysSorted returns the keys of m sorted in ascending order.
// K must satisfy cmp.Ordered (~string, integer, float types).
//
// Example:
//
//	keys := GetKeysSorted(map[string]int{"b": 2, "a": 1})
//	// keys == ["a", "b"]
func GetKeysSorted[K cmp.Ordered, V any](m map[K]V) []K {
	r := GetKeys[K, V](m)
	slices.Sort(r)
	return r
}

// GetValsUnique returns a slice of unique values from the map.
//
// Behavior details:
// - For a nil or empty map, it returns an empty slice.
// - Duplicate values are removed using value equality.
// - The order of returned values is not guaranteed because Go maps are unordered.
//
// Example:
//
//	values := GetValsUnique(map[int]string{
//		1: "a",
//		2: "a",
//		3: "b",
//	})
//	// values contains "a" and "b" (order not guaranteed)
func GetValsUnique[K, V comparable](m map[K]V) []V {
	m1 := make(map[V]bool)
	for _, v := range m {
		m1[v] = true
	}

	r := make([]V, 0, len(m1))
	for k := range m1 {
		r = append(r, k)
	}

	return r
}
