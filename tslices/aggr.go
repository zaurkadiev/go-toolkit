package tslices

// Aggregate groups items by keyFn and reduces each group with aggFn.
// If keyFn returns ok=false, the element is skipped.
// If aggFn returns ok=false, that group is skipped.
// Returns nil for nil input.
// Returns an empty non-nil map for non-nil empty input.
//
// Behavior details:
// - keyFn is called in order for each element to choose a grouping key.
// - aggFn receives the full slice of elements for a key.
// - If aggFn returns ok=false, the result for that key is omitted.
// - Output order is not guaranteed because maps are unordered.
//
// Example:
//
//	out := Aggregate([]int{1, 2, 3, 4, 5}, func(v int) (int, bool) {
//		return v % 2, true
//	}, func(values []int) (int, bool) {
//		sum := 0
//		for _, v := range values {
//			sum += v
//		}
//		return sum, true
//	})
//	// out == map[int]int{0: 6, 1: 9}
func Aggregate[T any, K comparable, V any](slice []T, keyFn func(T) (K, bool), aggFn func([]T) (V, bool)) map[K]V {
	if slice == nil {
		return nil
	}
	if len(slice) == 0 {
		return map[K]V{}
	}

	m := make(map[K][]T, len(slice))
	for _, v := range slice {
		if key, ok := keyFn(v); ok {
			m[key] = append(m[key], v)
		}
	}

	r := make(map[K]V, len(m))
	for k, v := range m {
		if ar, ok := aggFn(v); ok {
			r[k] = ar
		}
	}

	return r
}

// AggregateFuncKeyAsIs is a key function for Aggregate that returns the input as-is.
//
// It always returns ok=true, so every element is eligible as a key.
//
// Behavior details:
// - The returned key is exactly the input value.
// - Useful when elements themselves should be the grouping key.
//
// Example:
//
//	out := Aggregate([]string{"a", "b", "a"}, AggregateFuncKeyAsIs[string], func(values []string) (int, bool) {
//		return len(values), true
//	})
//	// out == map[string]int{"a": 2, "b": 1}
//
// See also: Aggregate.
func AggregateFuncKeyAsIs[T any](t T) (T, bool) {
	return t, true
}

// AggregateFuncCount returns the count of elements as type R.
// Intended to be used with Aggregate.
//
// It always returns ok=true, even for empty slices.
//
// Example:
//
//	out, ok := AggregateFuncCount[int, int]([]int{1, 2, 3})
//	// out == 3, ok == true
//
//	out, ok = AggregateFuncCount[int, int]([]int{})
//	// out == 0, ok == true
//
// See also: Aggregate.
func AggregateFuncCount[T any, R integer](slice []T) (R, bool) {
	return R(len(slice)), true
}

// AggregateFuncSum sums ordered values.
// Intended to be used with Aggregate.
//
// Returns ok=false for empty slices and zero value in that case.
//
// Example:
//
//	out, ok := AggregateFuncSum([]int{1, 2, 3})
//	// out == 6, ok == true
//
//	out, ok = AggregateFuncSum([]int{})
//	// out == 0, ok == false
//
// See also: Aggregate.
func AggregateFuncSum[T ordered](slice []T) (T, bool) {
	if len(slice) == 0 {
		return *new(T), false
	}

	var r T
	for _, v := range slice {
		r = r + v
	}

	return r, true
}

// AggregateFuncMin returns the minimum value in slice.
// Intended to be used with Aggregate.
//
// Returns ok=false for empty slices and zero value in that case.
//
// Example:
//
//	out, ok := AggregateFuncMin([]int{3, 1, 2})
//	// out == 1, ok == true
//
//	out, ok = AggregateFuncMin([]int{})
//	// out == 0, ok == false
//
// See also: Aggregate.
func AggregateFuncMin[T ordered](slice []T) (T, bool) {
	if len(slice) == 0 {
		return *new(T), false
	}

	r := slice[0]
	for _, v := range slice {
		if v < r {
			r = v
		}
	}

	return r, true
}

// AggregateFuncMax returns the maximum value in slice.
// Intended to be used with Aggregate.
//
// Returns ok=false for empty slices and zero value in that case.
//
// Example:
//
//	out, ok := AggregateFuncMax([]int{3, 1, 2})
//	// out == 3, ok == true
//
//	out, ok = AggregateFuncMax([]int{})
//	// out == 0, ok == false
//
// See also: Aggregate.
func AggregateFuncMax[T ordered](slice []T) (T, bool) {
	if len(slice) == 0 {
		return *new(T), false
	}

	r := slice[0]
	for _, v := range slice {
		if v > r {
			r = v
		}
	}

	return r, true
}

// AggregateFuncFirst returns the first element in slice.
// Intended to be used with Aggregate.
// Returns ok=false for empty slices and zero value in that case.
//
// Example:
//
//	out, ok := AggregateFuncFirst([]int{3, 1, 2})
//	// out == 3, ok == true
//
//	out, ok = AggregateFuncFirst([]int{})
//	// out == 0, ok == false
//
// See also: Aggregate.
func AggregateFuncFirst[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		return *new(T), false
	}

	return slice[0], true
}

// AggregateFuncLast returns the last element in slice.
// Intended to be used with Aggregate.
//
// Returns ok=false for empty slices and zero value in that case.
//
// Example:
//
//	out, ok := AggregateFuncLast([]int{3, 1, 2})
//	// out == 2, ok == true
//
//	out, ok = AggregateFuncLast([]int{})
//	// out == 0, ok == false
//
// See also: Aggregate.
func AggregateFuncLast[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		return *new(T), false
	}

	return slice[len(slice)-1], true
}

// AggregateErr is the error-returning variant of Aggregate.
// It stops on the first error from keyFn or aggFn and returns it.
// Returns (nil, nil) for nil input.
// Returns (empty non-nil map, nil) for non-nil empty input.
//
// Behavior details:
// - keyFn is called in order for each element until an error occurs.
// - aggFn is called per group; if aggFn returns ok=false, that group is omitted.
// - When err != nil from keyFn or aggFn, AggregateErr stops and returns (nil, err).
// - Output order is not guaranteed because maps are unordered.
//
// Example:
//
//	out, err := AggregateErr([]string{"1", "x", "2"}, func(v string) (int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return 0, false, err
//		}
//		return n % 2, true, nil
//	}, func(values []string) (int, bool, error) {
//		return len(values), true, nil
//	})
//	// err != nil (because of "x")
//
//	out, err = AggregateErr([]string{"1", "2", "3"}, func(v string) (int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return 0, false, err
//		}
//		return n % 2, true, nil
//	}, func(values []string) (int, bool, error) {
//		return len(values), true, nil
//	})
//	// out == map[int]int{0: 1, 1: 2}, err == nil
func AggregateErr[T any, K comparable, V any](slice []T, keyFn func(T) (K, bool, error), aggFn func([]T) (V, bool, error)) (map[K]V, error) {
	if slice == nil {
		return nil, nil
	}
	if len(slice) == 0 {
		return map[K]V{}, nil
	}

	m := make(map[K][]T, len(slice))
	for _, v := range slice {
		key, ok, err := keyFn(v)
		if err != nil {
			return nil, err
		}
		if ok {
			m[key] = append(m[key], v)
		}
	}

	r := make(map[K]V, len(m))
	for k, v := range m {
		ar, ok, err := aggFn(v)
		if err != nil {
			return nil, err
		}
		if ok {
			r[k] = ar
		}
	}

	return r, nil
}
