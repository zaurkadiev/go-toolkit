package pointers

// Of returns a pointer to the provided value.
//
// Behavior details:
// - The returned pointer refers to a copy of value.
// - This is useful when you need an addressable literal.
//
// Example:
//
//	ptr := Of(42)
//	// *ptr == 42
func Of[T any](value T) *T {
	return &value
}

// ValueOf safely dereferences ptr and returns the zero value if ptr is nil.
//
// Behavior details:
// - If ptr is nil, it returns the zero value for T.
// - Otherwise it returns the value pointed to by ptr.
//
// Example:
//
//	var ptr *int
//	value := ValueOf(ptr)
//	// value == 0
//
//	ptr = Of(5)
//	value = ValueOf(ptr)
//	// value == 5
func ValueOf[T any](ptr *T) T {
	if ptr == nil {
		return *new(T)
	}

	return *ptr
}
