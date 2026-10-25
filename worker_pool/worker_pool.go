package worker_pool

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// WorkerPool runs jobs that return a result on a fixed number of workers.
//
// Add all jobs before calling Run. Add and Run are not safe for concurrent use.
type WorkerPool[T any] interface {
	Add(job func() (T, error))
	Run(ctx context.Context) ([]T, error)
}

// impl manages a pool of workers to execute jobs concurrently.
type impl[T any] struct {
	jobs     []func() (T, error)
	poolSize int
}

// New returns a WorkerPool that runs jobs on n workers.
//
// n is validated in Run: Run returns an error if n is not positive.
//
// Example:
//
//	wp := New[int](2)
//	wp.Add(func() (int, error) { return 1, nil })
//	wp.Add(func() (int, error) { return 2, nil })
//	results, err := wp.Run(ctx)
//	// results == []int{1, 2}, err == nil
func New[T any](n int) WorkerPool[T] {
	return &impl[T]{
		jobs:     make([]func() (T, error), 0),
		poolSize: n,
	}
}

// Add adds a new task to the worker pool.
// It must not be called concurrently with Run.
//
// Example:
//
//	wp := New[string](4)
//	for _, id := range ids {
//		wp.Add(func() (string, error) { return fetchName(ctx, id) })
//	}
func (wp *impl[T]) Add(job func() (T, error)) {
	wp.jobs = append(wp.jobs, job)
}

// Run executes all added jobs on the pool's workers and returns their results in
// the order the jobs were added.
//
// Behavior details:
//   - With no jobs, it returns (nil, nil).
//   - If the pool size is not positive, it returns an error.
//   - If ctx is already done, it returns ctx.Err() without running jobs.
//   - On the first job error, no new jobs start and Run returns (nil, err).
//     Jobs that are already running are not interrupted: a job has no context,
//     so Run waits for it to return. Pass a context into jobs through a closure
//     if they must stop early.
//   - A panic in a job is recovered and returned as an error. If the panic value
//     is an error, it is wrapped, so errors.Is and errors.As match it.
//   - Run can be called again; it runs the same jobs again.
//
// Example:
//
//	wp := New[int](2)
//	for _, n := range []int{1, 2, 3} {
//		wp.Add(func() (int, error) { return n * n, nil })
//	}
//	squares, err := wp.Run(ctx)
//	// squares == []int{1, 4, 9}, err == nil
func (wp *impl[T]) Run(ctx context.Context) ([]T, error) {
	if len(wp.jobs) == 0 {
		return nil, nil
	}
	if wp.poolSize <= 0 {
		return nil, fmt.Errorf("worker size must be positive")
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	result := make([]T, len(wp.jobs))
	// errgroup cancels ctx after the first error is recorded, so Run returns the job error rather than context.Canceled
	eg, ctx := errgroup.WithContext(ctx)
	jobChan := make(chan int)

	// Start the workers
	for i := 0; i < wp.poolSize; i++ {
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case idx, ok := <-jobChan:
					if !ok {
						// if jobChan is closed, the worker exits
						return nil
					}
					res, err := call(wp.jobs[idx])
					if err != nil {
						return err
					}
					// each worker writes only to its own index; eg.Wait synchronizes the read
					result[idx] = res
				}
			}
		})
	}

	// A goroutine sends job indexes to the workers through jobChan
	eg.Go(func() error {
		defer close(jobChan)
		for idx := range wp.jobs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case jobChan <- idx:
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

// Simple runs jobs without results on a fixed number of workers.
//
// Add all jobs before calling Run. Add and Run are not safe for concurrent use.
type Simple interface {
	Add(job ...func() error)
	Run(ctx context.Context) error
}

// simpleImpl WorkerPoolSimple manages a pool of workers to execute jobs concurrently.
type simpleImpl struct {
	jobs     []func() error
	poolSize int
}

// NewSimple returns a Simple pool that runs jobs on n workers.
//
// n is validated in Run: Run returns an error if n is not positive.
//
// Example:
//
//	wp := NewSimple(2)
//	wp.Add(sendEmail, sendSMS)
//	err := wp.Run(ctx)
func NewSimple(n int) Simple {
	return &simpleImpl{
		jobs:     make([]func() error, 0),
		poolSize: n,
	}
}

// Add adds one or more tasks to the worker pool.
// It must not be called concurrently with Run.
//
// Example:
//
//	wp := NewSimple(4)
//	for _, id := range ids {
//		wp.Add(func() error { return notify(ctx, id) })
//	}
func (wp *simpleImpl) Add(jobs ...func() error) {
	wp.jobs = append(wp.jobs, jobs...)
}

// Run executes all added jobs on the pool's workers.
//
// Behavior details:
//   - With no jobs, it returns nil.
//   - If the pool size is not positive, it returns an error.
//   - If ctx is already done, it returns ctx.Err() without running jobs.
//   - On the first job error, no new jobs start and Run returns that error.
//     Jobs that are already running are not interrupted: a job has no context,
//     so Run waits for it to return. Pass a context into jobs through a closure
//     if they must stop early.
//   - A panic in a job is recovered and returned as an error. If the panic value
//     is an error, it is wrapped, so errors.Is and errors.As match it.
//   - Run can be called again; it runs the same jobs again.
//
// Example:
//
//	wp := NewSimple(2)
//	for _, id := range ids {
//		wp.Add(func() error { return notify(ctx, id) })
//	}
//	if err := wp.Run(ctx); err != nil {
//		return err
//	}
func (wp *simpleImpl) Run(ctx context.Context) error {
	if len(wp.jobs) == 0 {
		return nil
	}
	if wp.poolSize <= 0 {
		return fmt.Errorf("worker size must be positive")
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}

	jobChan := make(chan func() error)
	// errgroup cancels ctx after the first error is recorded, so Run returns the job error rather than context.Canceled
	eg, ctx := errgroup.WithContext(ctx)

	// Start the workers
	for i := 0; i < wp.poolSize; i++ {
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case job, ok := <-jobChan:
					if !ok {
						// if jobChan is closed, the worker exits
						return nil
					}
					if _, err := call(func() (struct{}, error) { return struct{}{}, job() }); err != nil {
						return err
					}
				}
			}
		})
	}

	// A goroutine sends jobs to the workers through jobChan
	eg.Go(func() error {
		defer close(jobChan)
		for _, job := range wp.jobs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case jobChan <- job:
			}
		}
		return nil
	})

	return eg.Wait()
}

// call runs job and converts a panic into an error, so a panicking job fails Run
// instead of crashing the program.
func call[T any](job func() (T, error)) (res T, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = fmt.Errorf("panic in worker pool job: %w", e)
			} else {
				err = fmt.Errorf("panic in worker pool job: %v", r)
			}
		}
	}()

	return job()
}
