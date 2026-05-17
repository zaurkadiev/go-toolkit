package tslices

// Iterate updates each element in slice by applying fn to every element.
//
// It performs in-place transformation: for every index in slice, fn is
// called with the current value, and the return value replaces the existing
// element.
//
// The iteration order is left-to-right, following the natural slice order.
// The function does not change the slice length or capacity.
//
// The provided fn should avoid relying on external mutation of the same slice
// during iteration. To transform elements based on position, capture the index
// in a closure by iterating with a separate index slice or by precomputing
// derived values.
//
// Example:
//
//	values := []int{1, 2, 3}
//	slice.Iterate(values, func(v int) int {
//		return v * v
//	})
//	// values is now []int{1, 4, 9}
func Iterate[T any](slice []T, fn func(T) T) {
	for i, t := range slice {
		slice[i] = fn(t)
	}
}
