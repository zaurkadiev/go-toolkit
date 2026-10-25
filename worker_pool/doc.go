// Package worker_pool runs jobs concurrently on a fixed number of workers.
//
// New returns a pool for jobs that produce results; Run returns them in the order
// the jobs were added. NewSimple returns a pool for jobs that only return an error.
//
// Both pools stop starting new jobs after the first error and return it. Jobs take
// no context, so jobs that are already running finish before Run returns. Panics in
// jobs are recovered and returned as errors.
//
// A pool is not safe for concurrent use: add all jobs, then call Run.
package worker_pool
