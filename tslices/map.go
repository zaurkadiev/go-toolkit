package tslices

// Map applies the mapper function to every element of the input slice
// and returns a new slice containing only the elements for which mapper
// returned ok == true.
//
// Signature matches the requirements from utils.md:
//
//	slices.Map[T, R]([]T, func(T) (R, bool)) []R
//
// Example:
//
//	out := Map([]int{1, 2, 3}, func(v int) (int, bool) {
//	    if v%2 == 0 {
//	        return v * 2, true
//	    }
//	    return 0, false
//	})
//	// out == []int{4}
func Map[T, R any](slice []T, fn func(T) (R, bool)) []R {
	if slice == nil {
		return nil
	}
	if len(slice) == 0 {
		return []R{}
	}

	r := make([]R, 0, len(slice))
	for _, v := range slice {
		if mapped, ok := fn(v); ok {
			r = append(r, mapped)
		}
	}

	return r
}

// MapErr applies a mapper that can return a value, a keep flag, and an error.
// It returns a new slice of mapped values where ok=true, or the first error
// encountered. If slice is empty, it returns (nil, nil).
//
// Behavior details:
// - The mapper is called in order for each element until an error occurs.
// - When err != nil, MapErr stops and returns (nil, err).
// - When ok == false, the mapped value is discarded.
// - The output preserves the order of elements for which ok == true.
//
// Example:
//
//	out, err := MapErr([]string{"1", "x", "2"}, func(v string) (int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return 0, false, err
//		}
//		return n, n%2 == 0, nil
//	})
//	// err != nil (because of "x")
//
//	out, err = MapErr([]string{"1", "2", "3"}, func(v string) (int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return 0, false, err
//		}
//		return n, n%2 == 0, nil
//	})
//	// out == []int{2}, err == nil
func MapErr[T, R any](slice []T, fn func(T) (R, bool, error)) ([]R, error) {
	if slice == nil {
		return nil, nil
	}
	if len(slice) == 0 {
		return []R{}, nil
	}

	r := make([]R, 0, len(slice))
	for _, v := range slice {
		mapped, ok, err := fn(v)
		if err != nil {
			return nil, err
		}
		if ok {
			r = append(r, mapped)
		}
	}

	return r, nil
}

// Map2Many maps each input element to a slice of outputs and flattens them.
// Only results with ok=true are appended. Returns nil for empty input.
//
// Behavior details:
// - The mapper is called in order for each element.
// - When ok == false, the mapped slice is discarded.
// - The output preserves the relative order of appended mapped slices.
//
// Example:
//
//	out := Map2Many([]int{1, 2, 3}, func(v int) ([]string, bool) {
//		if v%2 == 0 {
//			return []string{strconv.Itoa(v), strconv.Itoa(v * 10)}, true
//		}
//		return nil, false
//	})
//	// out == []string{"2", "20"}
func Map2Many[T, R any](slice []T, fn func(T) ([]R, bool)) []R {
	if slice == nil {
		return nil
	}
	if len(slice) == 0 {
		return []R{}
	}

	r := make([]R, 0, len(slice))
	for _, v := range slice {
		if mapped, ok := fn(v); ok {
			r = append(r, mapped...)
		}
	}

	return r
}

// Map2ManyErr is like Map2Many but the mapper can return an error.
// It stops on the first error and returns it. For empty input, returns (nil, nil).
//
// Behavior details:
// - The mapper is called in order for each element until an error occurs.
// - When err != nil, Map2ManyErr stops and returns (nil, err).
// - When ok == false, the mapped slice is discarded.
// - The output preserves the relative order of appended mapped slices.
//
// Example:
//
//	out, err := Map2ManyErr([]string{"1", "x", "2"}, func(v string) ([]int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return nil, false, err
//		}
//		if n%2 == 0 {
//			return []int{n, n * 10}, true, nil
//		}
//		return nil, false, nil
//	})
//	// err != nil (because of "x")
//
//	out, err = Map2ManyErr([]string{"1", "2", "3"}, func(v string) ([]int, bool, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return nil, false, err
//		}
//		if n%2 == 0 {
//			return []int{n, n * 10}, true, nil
//		}
//		return nil, false, nil
//	})
//	// out == []int{2, 20}, err == nil
func Map2ManyErr[T, R any](slice []T, fn func(T) ([]R, bool, error)) ([]R, error) {
	if slice == nil {
		return nil, nil
	}
	if len(slice) == 0 {
		return []R{}, nil
	}

	r := make([]R, 0, len(slice))
	for _, v := range slice {
		mapped, ok, err := fn(v)
		if err != nil {
			return nil, err
		}
		if ok {
			r = append(r, mapped...)
		}
	}

	return r, nil
}
