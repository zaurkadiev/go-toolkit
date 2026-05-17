package tslices

// Unique de-duplicates slice using fn to produce a key per element.
// When keys repeat, the last value for that key is kept.
// Output order is not guaranteed because map iteration is unordered.
//
// Behavior details:
// - fn is called for each element to produce a comparable key.
// - If multiple elements map to the same key, the last one wins.
// - The output does not preserve input order.
//
// Example:
//
//	out := Unique([]int{1, 2, 2, 3}, func(v int) int { return v })
//	// out contains 1, 2, 3 (order not guaranteed)
func Unique[T, R comparable](slice []T, fn func(T) R) []T {
	m := make(map[R]T)
	for _, t := range slice {
		m[fn(t)] = t
	}

	r := make([]T, 0, len(m))
	for _, t := range m {
		r = append(r, t)
	}

	return r
}

// UniqueOrdered de-duplicates a slice using fn to produce a key per element.
// When keys repeat, the first value for that key is kept.
// Output order matches the input order of the first occurrences.
//
// Behavior details:
//   - fn is called for each element to produce a comparable key.
//   - If multiple elements map to the same key, the first one wins and later
//     duplicates are skipped.
//   - The output preserves the order of first appearances in the input.
//
// Example:
//
//	out := UniqueOrdered([]int{1, 2, 1, 2, 3}, func(v int) int { return v })
//	// out contains 1, 2, 3 in that order
func UniqueOrdered[T, R comparable](slice []T, fn func(T) R) []T {
	m := make(map[R]bool)
	r := make([]T, 0)
	for _, t := range slice {
		value := fn(t)
		if m[value] {
			continue
		}
		m[value] = true
		r = append(r, t)
	}

	return r
}

// UniqueFuncAsIs is an identity key function for Unique or UniqueOrdered.
// Intended to be used with Unique or UniqueOrdered.
//
// Use when the element itself should be the key.
//
// Example:
//
//	out := Unique([]int{1, 2, 2, 3}, UniqueFuncAsIs[int])
//	// out contains 1, 2, 3 (order not guaranteed)
//
// See also: Unique.
func UniqueFuncAsIs[T comparable](value T) T {
	return value
}
