package tslices

// Contains reports whether any element satisfies fn.
//
// Behavior details:
// - Iterates in order and stops at the first match.
// - Returns false if no elements match.
//
// Example:
//
//	ok := Contains([]int{1, 2, 3}, func(v int) bool {
//		return v == 2
//	})
//	// ok == true
func Contains[T any](slice []T, fn func(T) bool) bool {
	for _, t := range slice {
		if fn(t) {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first element that satisfies fn.
//
// Behavior details:
// - Iterates in order and stops at the first match.
// - Returns -1 if no element matches or slice is empty.
//
// Example:
//
//	i := IndexOf([]int{1, 2, 3}, func(v int) bool {
//		return v == 2
//	})
//	// i == 1
func IndexOf[T any](slice []T, fn func(T) bool) int {
	for i, t := range slice {
		if fn(t) {
			return i
		}
	}
	return -1
}

// FindFirst returns the first element satisfying fn and its index.
//
// Behavior details:
// - Iterates in order and stops at the first match.
// - For empty slices, returns the zero value and index -1.
// - If no element matches, returns the zero value and index -1.
//
// Example:
//
//	out, idx := FindFirst([]int{5, 8, 9}, func(v int) bool {
//		return v%2 == 0
//	})
//	// out == 8, idx == 1
func FindFirst[T any](slice []T, fn func(T) bool) (T, int) {
	if len(slice) == 0 {
		return *new(T), -1
	}

	for i := 0; i < len(slice); i++ {
		t := slice[i]
		if fn(t) {
			return t, i
		}
	}

	return *new(T), -1
}

// FindLast returns the last element satisfying fn and its index.
//
// Behavior details:
// - Iterates from the end and stops at the first match.
// - For empty slices, returns the zero value and index -1.
// - If no element matches, returns the zero value and index -1.
//
// Example:
//
//	out, idx := FindLast([]int{1, 2, 3, 2}, func(v int) bool {
//		return v == 2
//	})
//	// out == 2, idx == 3
func FindLast[T any](slice []T, fn func(T) bool) (T, int) {
	if len(slice) == 0 {
		return *new(T), -1
	}

	for i := len(slice) - 1; i >= 0; i-- {
		t := slice[i]
		if fn(t) {
			return t, i
		}
	}

	return *new(T), -1
}
