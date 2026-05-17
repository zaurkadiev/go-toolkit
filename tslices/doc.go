// Package tslices provides generic helpers for working with slices.
//
// The package focuses on small, composable building blocks for common
// collection workflows such as mapping, filtering, grouping, aggregation,
// and simple searches. Functions are designed to be explicit and predictable:
// mapper callbacks often return an ok flag for filtering, and error-aware
// variants stop at the first error.
//
// Many helpers accept a function to derive keys or values, which makes them
// suitable for both primitive and structured data. The implementations avoid
// hidden side effects, favoring straightforward iteration so behavior is easy
// to reason about and profile.
package tslices
