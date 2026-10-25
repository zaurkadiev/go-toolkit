package tsync

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSemaphore(t *testing.T) {
	type args struct {
		n uint
	}
	type testCase struct {
		name string
		args args
		want int
	}
	tests := []testCase{
		{
			name: "given_positive_size_when_new_semaphore_then_capacity",
			args: args{
				n: 3,
			},
			want: 3,
		},
		{
			name: "given_zero_size_when_new_semaphore_then_no_slots",
			args: args{
				n: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, cap(NewSemaphore(tt.args.n).ch), "NewSemaphore(%v)", tt.args.n)
		})
	}
}

func TestSemaphore_Acquire(t *testing.T) {
	type fields struct {
		n        uint
		acquired int
	}
	type args struct {
		ctx context.Context
	}
	type testCase struct {
		name         string
		fields       fields
		args         args
		wantAcquired int
		wantErr      assert.ErrorAssertionFunc
	}
	tests := []testCase{
		{
			name: "given_free_slot_when_acquire_then_acquired",
			fields: fields{
				n:        2,
				acquired: 1,
			},
			args: args{
				ctx: context.Background(),
			},
			wantAcquired: 2,
			wantErr:      assert.NoError,
		},
		{
			name: "given_no_free_slot_when_context_cancelled_then_context_error",
			fields: fields{
				n:        1,
				acquired: 1,
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
			},
			wantAcquired: 1,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, context.Canceled, i...)
			},
		},
		{
			name: "given_no_free_slot_when_context_times_out_then_deadline_error",
			fields: fields{
				n:        1,
				acquired: 1,
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
					t.Cleanup(cancel)
					return ctx
				}(),
			},
			wantAcquired: 1,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, context.DeadlineExceeded, i...)
			},
		},
		{
			name: "given_zero_size_when_context_times_out_then_deadline_error",
			fields: fields{
				n:        0,
				acquired: 0,
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
					t.Cleanup(cancel)
					return ctx
				}(),
			},
			wantAcquired: 0,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, context.DeadlineExceeded, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSemaphore(tt.fields.n)
			for range tt.fields.acquired {
				s.ch <- struct{}{}
			}
			if !tt.wantErr(t, s.Acquire(tt.args.ctx), fmt.Sprintf("Acquire(%v)", tt.args.ctx)) {
				return
			}
			assert.Equalf(t, tt.wantAcquired, len(s.ch), "Acquire(%v)", tt.args.ctx)
		})
	}
}

func TestSemaphore_TryAcquire(t *testing.T) {
	type fields struct {
		n        uint
		acquired int
	}
	type testCase struct {
		name         string
		fields       fields
		want         bool
		wantAcquired int
	}
	tests := []testCase{
		{
			name: "given_free_slot_when_try_acquire_then_true",
			fields: fields{
				n:        2,
				acquired: 1,
			},
			want:         true,
			wantAcquired: 2,
		},
		{
			name: "given_no_free_slot_when_try_acquire_then_false",
			fields: fields{
				n:        1,
				acquired: 1,
			},
			want:         false,
			wantAcquired: 1,
		},
		{
			name: "given_zero_size_when_try_acquire_then_false",
			fields: fields{
				n:        0,
				acquired: 0,
			},
			want:         false,
			wantAcquired: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSemaphore(tt.fields.n)
			for range tt.fields.acquired {
				s.ch <- struct{}{}
			}
			assert.Equalf(t, tt.want, s.TryAcquire(), "TryAcquire()")
			assert.Equalf(t, tt.wantAcquired, len(s.ch), "TryAcquire()")
		})
	}
}

func TestSemaphore_Release(t *testing.T) {
	type fields struct {
		n        uint
		acquired int
	}
	type testCase struct {
		name         string
		fields       fields
		want         bool
		wantAcquired int
	}
	tests := []testCase{
		{
			name: "given_held_slot_when_release_then_true_and_slot_freed",
			fields: fields{
				n:        2,
				acquired: 2,
			},
			want:         true,
			wantAcquired: 1,
		},
		{
			name: "given_no_held_slot_when_release_then_false",
			fields: fields{
				n:        1,
				acquired: 0,
			},
			want:         false,
			wantAcquired: 0,
		},
		{
			name: "given_zero_size_when_release_then_false",
			fields: fields{
				n:        0,
				acquired: 0,
			},
			want:         false,
			wantAcquired: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSemaphore(tt.fields.n)
			for range tt.fields.acquired {
				s.ch <- struct{}{}
			}
			assert.Equalf(t, tt.want, s.Release(), "Release()")
			assert.Equalf(t, tt.wantAcquired, len(s.ch), "Release()")
		})
	}
}

func TestSemaphore_ConcurrentAcquire(t *testing.T) {
	type args struct {
		n          uint
		goroutines int
	}
	type testCase struct {
		name string
		args args
		want int32
	}
	tests := []testCase{
		{
			name: "given_single_slot_when_goroutines_acquire_then_one_at_a_time",
			args: args{
				n:          1,
				goroutines: 10,
			},
			want: 1,
		},
		{
			name: "given_three_slots_when_goroutines_acquire_then_at_most_three",
			args: args{
				n:          3,
				goroutines: 20,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSemaphore(tt.args.n)
			var current, peak atomic.Int32
			wg := sync.WaitGroup{}
			for range tt.args.goroutines {
				wg.Go(func() {
					assert.NoError(t, s.Acquire(context.Background()))
					defer s.Release()
					c := current.Add(1)
					for {
						p := peak.Load()
						if c <= p || peak.CompareAndSwap(p, c) {
							break
						}
					}
					time.Sleep(5 * time.Millisecond)
					current.Add(-1)
				})
			}
			wg.Wait()
			assert.Equalf(t, tt.want, peak.Load(), "peak holders with NewSemaphore(%v)", tt.args.n)
		})
	}
}
