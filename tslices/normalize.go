package tslices

// Nil2Empty returns an empty non-nil slice when slice is nil.
// Non-nil slices are returned unchanged.
//
// Behavior details:
// - Nil input returns a new zero-length slice.
// - Non-nil empty and populated slices are returned as-is.
// - Named slice types are preserved.
//
// Example:
//
//	out := NilToEmpty([]int(nil))
//	// out == []int{}
func Nil2Empty[T any](slice []T) []T {
	if slice == nil {
		return make([]T, 0)
	}

	return slice
}
