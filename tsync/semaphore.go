package tsync

import "context"

// Semaphore limits how many goroutines can hold it at the same time.
//
// Use it to cap concurrency, for example the number of parallel requests to an
// external service. Every successful Acquire or TryAcquire must be paired with
// exactly one Release.
//
// Waiters are not served in FIFO order. For weighted acquisition (taking several
// slots at once) use golang.org/x/sync/semaphore.Weighted.
//
// The zero Semaphore has no slots, like NewSemaphore(0).
//
// Example:
//
//	sem := NewSemaphore(3)
//	for _, id := range ids {
//		if err := sem.Acquire(ctx); err != nil {
//			return err
//		}
//		go func() {
//			defer sem.Release()
//			process(ctx, id)
//		}()
//	}
type Semaphore struct {
	ch chan any
}

// NewSemaphore returns a Semaphore that can be held by at most n goroutines at once.
//
// A Semaphore of size 0 has no slots: Acquire waits until ctx is done and
// TryAcquire always returns false.
//
// Example:
//
//	sem := NewSemaphore(10)
func NewSemaphore(n uint) *Semaphore {
	if n == 0 {
		// nil channel: sends never proceed, so the semaphore is never acquired
		return &Semaphore{}
	}

	return &Semaphore{
		ch: make(chan any, n),
	}
}

// Acquire takes a slot, blocking until one is free or ctx is done.
// On success, it returns nil; otherwise it returns ctx.Err() and takes no slot.
// If ctx is already done but a slot is free, Acquire may still succeed.
//
// Example:
//
//	if err := sem.Acquire(ctx); err != nil {
//		return err
//	}
//	defer sem.Release()
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// TryAcquire takes a slot without blocking and reports whether it succeeded.
//
// Example:
//
//	if !sem.TryAcquire() {
//		return ErrBusy
//	}
//	defer sem.Release()
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release frees a slot taken by Acquire or TryAcquire and reports whether a slot
// was freed. It returns false if no slot is held, which means Release was called
// more times than Acquire and TryAcquire succeeded; the call is then a no-op.
//
// The semaphore does not track which goroutine holds a slot, so an extra Release
// while other goroutines hold slots frees one of their slots and returns true.
//
// Example:
//
//	if err := sem.Acquire(ctx); err != nil {
//		return err
//	}
//	defer sem.Release()
func (s *Semaphore) Release() bool {
	select {
	case <-s.ch:
		return true
	default:
		return false
	}
}
