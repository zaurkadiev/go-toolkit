// Package pointers provides generic helpers for working with pointers.
//
// The package includes small utilities for constructing pointers from values
// and safely dereferencing possibly nil pointers without repetitive
// boilerplate. These helpers are useful when code needs an addressable value
// or wants to treat nil as an explicit "unset" state.
//
// Functions are intentionally minimal and side-effect free, which keeps their
// behavior predictable across value types.
package pointers
