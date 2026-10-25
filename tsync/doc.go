// Package tsync provides typed wrappers over the standard library sync primitives.
//
// Map is a type-safe wrapper over sync.Map: keys and values are typed, so callers
// do not need type assertions. Semantics, concurrency guarantees and performance
// characteristics are the same as for sync.Map.
//
// RWMap offers the same API over a plain map guarded by sync.RWMutex, plus Len.
// See the Map and RWMap docs for when to choose each.
//
// Semaphore limits how many goroutines can hold it at the same time.
//
// The package is named tsync so it does not clash with the standard library sync
// package when both are imported in one file:
//
//	import "github.com/zaurkadiev/go-toolkit/tsync"
package tsync
