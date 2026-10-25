// Package future runs functions in background goroutines and lets callers await
// their results later.
//
// New starts a function immediately and returns a Future; Await blocks until the
// function returns. Join awaits a slice of futures and collects their results in
// order. Panics inside the function are recovered and returned as errors.
//
// Futures are not cancellable: Await always waits for the function to return, so
// long-running functions should take a context through a closure.
package future
