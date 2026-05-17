// Package tmaps provides generic helpers for working with Go maps.
//
// The package includes utilities for extracting keys and values, transforming
// maps, filtering entries, and running simple searches over map contents.
// Transformation helpers accept callbacks that can both map and filter pairs,
// while error-aware variants stop on the first error.
//
// These helpers aim to keep map handling concise without hiding semantics.
// Be aware that map iteration order in Go is not deterministic, so functions that
// return slices or the "first" match do not guarantee ordering.
package tmaps
