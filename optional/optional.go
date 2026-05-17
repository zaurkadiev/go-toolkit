package optional

import (
	"errors"
	"fmt"
)

// ErrEmpty is returned by Unwrap when the optional has no value.
var ErrEmpty = errors.New("optional is empty")

// O represents an optional value for type T.
//
// O is empty when Value is nil.
type O[T any] struct {
	// Value stores the underlying pointer.
	// Nil means the optional is empty.
	Value *T
}

// New creates an optional from the provided pointer.
//
// Passing nil returns an empty optional.
func New[T any](value *T) O[T] {
	return O[T]{Value: value}
}

// NewZeroAsEmpty creates an optional and treats the zero value as empty.
//
// This helper is useful when values come from sources that use zero values
// (for example 0, "", false) to represent absence.
func NewZeroAsEmpty[T comparable](value T) O[T] {
	var zero T
	if value == zero {
		return O[T]{Value: nil}
	}
	return O[T]{Value: &value}
}

// Empty returns an empty optional for type T.
func Empty[T any]() O[T] {
	return O[T]{}
}

// Map applies f to the contained value and returns the result as an optional.
// Returns an empty optional when t is empty.
func Map[T, R any](t O[T], f func(T) R) O[R] {
	if t.Value == nil {
		return O[R]{}
	}
	value := f(*t.Value)
	return O[R]{Value: &value}
}

// Get returns the contained value.
//
// If the optional is empty, Get returns the zero value of T.
func (o O[T]) Get() T {
	if o.Value == nil {
		return *new(T)
	}

	return *o.Value
}

// IsEmpty reports whether the optional contains no value.
func (o O[T]) IsEmpty() bool {
	return o.Value == nil
}

// GetOrDefault returns the contained value when present.
//
// If the optional is empty, it returns defaultValue.
func (o O[T]) GetOrDefault(defaultValue T) T {
	if o.Value == nil {
		return defaultValue
	}
	return *o.Value
}

// Unwrap returns the contained value.
//
// If the optional is empty, Unwrap returns the zero value of T and ErrEmpty.
func (o O[T]) Unwrap() (T, error) {
	if o.Value == nil {
		return *new(T), ErrEmpty
	}
	return *o.Value, nil
}

// String returns a string representation of the optional value.
//
// It returns "<empty>" when the optional has no value. Otherwise it formats the
// contained value with fmt.Sprint. This representation is useful while
// serializing optional values into logs, debug output, or text-based payloads.
func (o O[T]) String() string {
	if o.Value == nil {
		return "<empty>"
	}

	return fmt.Sprint(*o.Value)
}
