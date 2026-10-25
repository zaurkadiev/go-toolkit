package future

import (
	"fmt"
	"sync"
)

// Future is the result of a function running in its own goroutine.
//
// Await blocks until the function returns and then returns its result. It is safe
// to call Await many times and from many goroutines: every call returns the same
// result.
type Future[T any] interface {
	Await() (T, error)
}

type impl[T any] struct {
	once   *sync.Once
	wg     *sync.WaitGroup
	result T
	err    error
	resCh  chan T
	errCh  chan error
}

// Await blocks until the function passed to New returns and returns its result.
// Repeated calls return the same result without blocking.
//
// Example:
//
//	f := New(func() (int, error) { return 42, nil })
//	v, err := f.Await()
//	// v == 42, err == nil
func (i *impl[T]) Await() (T, error) {
	i.once.Do(func() {
		i.wg.Add(1)
		defer i.wg.Done()
		i.result = <-i.resCh
		i.err = <-i.errCh
		close(i.resCh)
		close(i.errCh)
	})
	i.wg.Wait()

	return i.result, i.err
}

// New starts f in a new goroutine and returns a Future for its result.
//
// Behavior details:
//   - f starts immediately, not on the first Await.
//   - If f panics, the panic is recovered and Await returns the zero value of T
//     and an error. If the panic value is an error, it is wrapped, so errors.Is and
//     errors.As match it.
//   - f cannot be cancelled: there is no context, and Await always waits for f to
//     return. Pass a context into f through a closure if it must stop early.
//   - If Await is never called, the goroutine still finishes and does not leak.
//
// Example:
//
//	f := New(func() (string, error) {
//		return fetchName(ctx, id)
//	})
//	// ... do other work ...
//	name, err := f.Await()
func New[T any](f func() (T, error)) Future[T] {
	resCh := make(chan T, 1)
	errCh := make(chan error, 1)

	i := &impl[T]{
		resCh: resCh,
		errCh: errCh,
		once:  &sync.Once{},
		wg:    &sync.WaitGroup{},
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					errCh <- fmt.Errorf("panic in future: %w", err)
				} else {
					errCh <- fmt.Errorf("panic in future: %v", r)
				}
				resCh <- *new(T)
			}
		}()
		res, err := f()
		resCh <- res
		errCh <- err
	}()

	return i
}

// Join awaits futures in slice order and returns their results in the same order.
//
// Behavior details:
//   - For a nil slice, it returns (nil, nil).
//   - For an empty non-nil slice, it returns an empty slice and nil.
//   - It stops at the first error in slice order and returns (nil, err). An error
//     from a later future is reported only after the earlier ones finish.
//   - Futures that are not awaited keep running; Join does not cancel them.
//
// Example:
//
//	results, err := Join([]Future[int]{
//		New(func() (int, error) { return 1, nil }),
//		New(func() (int, error) { return 2, nil }),
//	})
//	// results == []int{1, 2}, err == nil
func Join[T any](futures []Future[T]) ([]T, error) {
	if futures == nil {
		return nil, nil
	}
	if len(futures) == 0 {
		return []T{}, nil
	}

	result := make([]T, 0, len(futures))
	for _, f := range futures {
		r, err := f.Await()
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}
