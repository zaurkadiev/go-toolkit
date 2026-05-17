package tslices

// Sum maps each element to a value and returns the total sum.
//
// Behavior details:
// - fn is called for each element to produce a summand.
// - The summands are accumulated in order.
// - For empty slices, it returns the zero value for R.
//
// Example:
//
//	out := Sum([]int{1, 2, 3}, func(v int) int { return v })
//	// out == 6
func Sum[T any, R ordered](slice []T, fn func(T) R) R {
	var r R
	for _, t := range slice {
		r += fn(t)
	}

	return r
}

// SumFuncAsIs is an identity mapper for Sum.
// It returns the input value unchanged.
//
// Example:
//
//	out := Sum([]int{1, 2, 3}, SumFuncAsIs[int])
//	// out == 6
//
// See also: Sum.
func SumFuncAsIs[T ordered](t T) T {
	return t
}

// Max maps elements and returns the maximum mapped value.
//
// Behavior details:
// - fn is called for each element to produce a comparable value.
// - The maximum of the mapped values is returned.
// - For empty slices, the zero value for R is returned.
//
// Example:
//
//	out := Max([]int{1, 2, 3}, func(v int) int { return v })
//	// out == 3
func Max[T any, R ordered](slice []T, fn func(T) R) R {
	if len(slice) == 0 {
		return *new(R)
	}
	r := fn(slice[0])
	for _, t := range slice {
		m := fn(t)
		if m > r {
			r = m
		}
	}

	return r
}

// MaxFuncAsIs is an identity mapper for Max.
// It returns the input value unchanged.
//
// Example:
//
//	out := Max([]int{1, 2, 3}, MaxFuncAsIs[int])
//	// out == 3
//
// See also: Max.
func MaxFuncAsIs[T ordered](t T) T {
	return t
}

// Min maps elements and returns the minimum mapped value.
//
// Behavior details:
// - fn is called for each element to produce a comparable value.
// - The minimum of the mapped values is returned.
// - For empty slices, the zero value for R is returned.
//
// Example:
//
//	out := Min([]int{1, 2, 3}, func(v int) int { return v })
//	// out == 1
func Min[T any, R ordered](slice []T, fn func(T) R) R {
	if len(slice) == 0 {
		return *new(R)
	}
	r := fn(slice[0])
	for _, t := range slice {
		m := fn(t)
		if m < r {
			r = m
		}
	}

	return r
}

// MinFuncAsIs is an identity mapper for Min.
// It returns the input value unchanged.
//
// Example:
//
//	out := Min([]int{1, 2, 3}, MinFuncAsIs[int])
//	// out == 1
//
// See also: Min.
func MinFuncAsIs[T ordered](t T) T {
	return t
}
