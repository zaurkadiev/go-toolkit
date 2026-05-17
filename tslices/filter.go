package tslices

// Filter returns a new slice containing elements where fn returns true.
// The output preserves the original order.
//
// Behavior details:
// - fn is called in order for each element.
// - Elements for which fn returns false are excluded.
// - The output is a new slice; the input is not modified.
//
// Example:
//
//	out := Filter([]int{1, 2, 3, 4}, func(v int) bool {
//		return v%2 == 0
//	})
//	// out == []int{2, 4}
func Filter[T any](slice []T, fn func(T) bool) []T {
	r := make([]T, 0, len(slice))
	for _, t := range slice {
		if fn(t) {
			r = append(r, t)
		}
	}

	return r
}

// FilterFuncNotZero returns true when t is not zero.
// Intended to be used with Filter.
//
// Behavior details:
// - Intended for use with nil-able element types (pointers, maps, slices, funcs, channels, interfaces).
// - When used with such types, it filters out nil elements.
//
// Example:
//
//	out := Filter([][]int{{1, 2}, nil, {3, 4}}, FilterFuncNotZero[[]int])
//	// out == [][]int{{1, 2}, {3, 4}}
//
// See also: Filter.
func FilterFuncNotZero[T comparable](t T) bool {
	var zero T
	return t != zero
}
