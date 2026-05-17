package tslices

import (
	"math/rand/v2"
)

// Chunk splits slice into sub-slices of the given size.
//
// Behavior details:
// - For size <= 0 or empty input, it returns an empty slice of slices.
// - The last chunk may be smaller than size.
// - The returned sub-slices share the underlying array with the input.
//
// Example:
//
//	out := Chunk([]int{1, 2, 3, 4, 5}, 2)
//	// out == [][]int{{1, 2}, {3, 4}, {5}}
func Chunk[T any](slice []T, size int) [][]T {
	if len(slice) == 0 || size <= 0 {
		return [][]T{}
	}

	chunks := (len(slice) + size - 1) / size
	r := make([][]T, 0, chunks)
	for len(slice) != 0 {
		if len(slice) < size {
			size = len(slice)
		}

		r = append(r, slice[0:size])
		slice = slice[size:]
	}

	return r
}

// ToSet builds a set (map with bool values) from slice using fn as a key.
// Duplicate keys are coalesced to a single entry.
//
// Behavior details:
// - fn is called for each element to produce a set key.
// - The resulting map contains each key exactly once with value true.
// - The order of keys is not defined because maps are unordered.
//
// Example:
//
//	out := ToSet([]int{1, 2, 2, 3}, func(v int) int { return v })
//	// out == map[int]bool{1: true, 2: true, 3: true}
func ToSet[T, R comparable](slice []T, fn func(T) R) map[R]bool {
	m := make(map[R]bool, len(slice))
	for _, t := range slice {
		m[fn(t)] = true
	}

	return m
}

// ToSetFuncAsIs is an identity mapper for ToSet.
// It returns the input value unchanged.
//
// Example:
//
//	out := ToSet([]int{1, 2}, ToSetFuncAsIs[int])
//	// out == map[int]bool{1: true, 2: true}
//
// See also: ToSet.
func ToSetFuncAsIs[T any](t T) T {
	return t
}

// Shuffle randomizes the order of elements in slice in place.
//
// Behavior details:
// - The input slice is modified in place and also returned for convenience.
// - For slices with fewer than two elements, it returns the input unchanged.
// - The randomization uses math/rand/v2's global source.
//
// Example:
//
//	values := []int{1, 2, 3}
//	out := Shuffle(values)
//	// values and out contain the same elements in random order
func Shuffle[T any](slice []T) []T {
	if len(slice) < 2 {
		return slice
	}

	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice
}

// FirstN returns the first n elements of slice.
//
// Behavior details:
// - If n <= 0 or slice is empty, it returns an empty slice.
// - If n >= len(slice), it returns the original slice.
// - The returned slice shares the underlying array with the input.
//
// Example:
//
//	out := FirstN([]int{1, 2, 3}, 2)
//	// out == []int{1, 2}
func FirstN[T any](slice []T, n int) []T {
	if n <= 0 || len(slice) == 0 {
		return []T{}
	}
	if len(slice) < n {
		return slice
	}
	return slice[:n]
}

// LastN returns the last n elements of slice.
//
// Behavior details:
// - If n <= 0 or slice is empty, it returns an empty slice.
// - If n >= len(slice), it returns the original slice.
// - The returned slice shares the underlying array with the input.
//
// Example:
//
//	out := LastN([]int{1, 2, 3}, 2)
//	// out == []int{2, 3}
func LastN[T any](slice []T, n int) []T {
	if n <= 0 || len(slice) == 0 {
		return []T{}
	}
	if len(slice) < n {
		return slice
	}
	return slice[len(slice)-n:]
}

// Reduce folds slice into an accumulator starting from a.
//
// Behavior details:
// - fn is called in order for each element.
// - The accumulator value is updated on each step.
// - For empty slices, it returns the initial accumulator unchanged.
//
// Example:
//
//	sum := Reduce([]int{1, 2, 3}, 0, func(acc, v int) int {
//		return acc + v
//	})
//	// sum == 6
func Reduce[T, A any](slice []T, a A, fn func(A, T) A) A {
	if len(slice) == 0 {
		return a
	}

	for _, t := range slice {
		a = fn(a, t)
	}

	return a
}

// ReduceErr is the error-returning variant of Reduce.
// It stops on the first error from fn and returns it.
//
// Behavior details:
// - fn is called in order for each element.
// - The accumulator value is updated on each step.
// - If fn returns an error, ReduceErr stops and returns (a, error).
// - For empty slices, it returns the initial accumulator unchanged with nil error.
//
// Example:
//
//	sum, err := ReduceErr([]string{"1", "2", "3"}, 0, func(acc int, v string) (int, error) {
//		n, err := strconv.Atoi(v)
//		if err != nil {
//			return acc, err
//		}
//		return acc + n, nil
//	})
//	// sum == 6, err == nil
func ReduceErr[T, A any](slice []T, a A, fn func(A, T) (A, error)) (A, error) {
	if len(slice) == 0 {
		return a, nil
	}

	var err error
	for _, t := range slice {
		a, err = fn(a, t)
		if err != nil {
			return a, err
		}
	}

	return a, nil
}

// ReduceFuncSum returns the sum of t and a.
// It is intended to be used as a reducer function with Reduce for ordered
// numeric types, mirroring the signature expected by Reduce's fn parameter.
//
// Example:
//
//	sum := Reduce([]int{1, 2, 3}, 0, ReduceFuncSum[int])
//	// sum == 6
//
// See also: Reduce.
func ReduceFuncSum[A ordered](t, a A) A {
	return a + t
}

// Zip pairs elements from left and right by index.
// If lengths differ, missing values are filled with zero values.
//
// Behavior details:
// - The output length equals max(len(left), len(right)).
// - Missing elements are replaced with the zero value of the corresponding type.
// - The output preserves index alignment between left and right.
//
// Example:
//
//	out := Zip([]int{1}, []string{"a", "b"})
//	// out == []Pair[int, string]{{Left: 1, Right: "a"}, {Left: 0, Right: "b"}}
func Zip[L, R any](left []L, right []R) []Pair[L, R] {
	size := max(len(left), len(right))
	r := make([]Pair[L, R], 0, size)
	for i := 0; i < size; i++ {
		p := Pair[L, R]{
			Left:  byIndex(left, i),
			Right: byIndex(right, i),
		}
		r = append(r, p)
	}

	return r
}

func byIndex[T any](slice []T, index int) T {
	if index < 0 {
		return *new(T)
	}

	if index < len(slice) {
		return slice[index]
	}

	return *new(T)
}
