package future

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args[T any] struct {
		f func() (T, error)
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		want    T
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[int]{
		{
			name: "given_f_returns_value_when_await_then_value",
			args: args[int]{
				f: func() (int, error) { return 1, nil },
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "given_f_returns_error_when_await_then_error",
			args: args[int]{
				f: func() (int, error) { return 0, assert.AnError },
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
		{
			name: "given_f_panics_with_value_when_await_then_error",
			args: args[int]{
				f: func() (int, error) { panic("panic") },
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.EqualError(t, err, "panic in future: panic", i...)
			},
		},
		{
			name: "given_f_panics_with_error_when_await_then_wrapped_error",
			args: args[int]{
				f: func() (int, error) { panic(assert.AnError) },
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.f).Await()
			if !tt.wantErr(t, err, fmt.Sprintf("New(%T).Await()", tt.args.f)) {
				return
			}
			assert.Equalf(t, tt.want, got, "New(%T).Await()", tt.args.f)
		})
	}
}

func Test_impl_Await(t *testing.T) {
	type fields[T any] struct {
		f func() (T, error)
	}
	type args struct {
		calls int
	}
	type testCase[T any] struct {
		name    string
		fields  fields[T]
		args    args
		want    T
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[int]{
		{
			name: "given_single_await_when_f_returns_then_value",
			fields: fields[int]{
				f: func() (int, error) { return 1, nil },
			},
			args: args{
				calls: 1,
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "given_concurrent_awaits_when_f_returns_then_same_value",
			fields: fields[int]{
				f: func() (int, error) {
					time.Sleep(10 * time.Millisecond)
					return 1, nil
				},
			},
			args: args{
				calls: 10,
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "given_concurrent_awaits_when_f_returns_error_then_same_error",
			fields: fields[int]{
				f: func() (int, error) {
					time.Sleep(10 * time.Millisecond)
					return 0, assert.AnError
				},
			},
			args: args{
				calls: 10,
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(tt.fields.f)
			got := make([]int, tt.args.calls)
			errs := make([]error, tt.args.calls)
			wg := sync.WaitGroup{}
			for i := range tt.args.calls {
				wg.Go(func() { got[i], errs[i] = f.Await() })
			}
			wg.Wait()
			for i := range tt.args.calls {
				if !tt.wantErr(t, errs[i], fmt.Sprintf("Await() call %d", i)) {
					return
				}
				assert.Equalf(t, tt.want, got[i], "Await() call %d", i)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type args[T any] struct {
		futures []Future[T]
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		want    []T
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[int]{
		{
			name: "given_nil_futures_when_join_then_nil",
			args: args[int]{
				futures: nil,
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "given_empty_futures_when_join_then_empty_slice",
			args: args[int]{
				futures: []Future[int]{},
			},
			want:    []int{},
			wantErr: assert.NoError,
		},
		{
			name: "given_futures_when_join_then_results_in_order",
			args: args[int]{
				futures: []Future[int]{
					New(func() (int, error) {
						time.Sleep(20 * time.Millisecond)
						return 1, nil
					}),
					New(func() (int, error) {
						time.Sleep(10 * time.Millisecond)
						return 2, nil
					}),
					New(func() (int, error) { return 3, nil }),
				},
			},
			want:    []int{1, 2, 3},
			wantErr: assert.NoError,
		},
		{
			name: "given_future_with_error_when_join_then_error",
			args: args[int]{
				futures: []Future[int]{
					New(func() (int, error) { return 1, nil }),
					New(func() (int, error) { return 0, assert.AnError }),
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
		{
			name: "given_future_with_panic_when_join_then_error",
			args: args[int]{
				futures: []Future[int]{
					New(func() (int, error) { return 1, nil }),
					New(func() (int, error) { panic("panic") }),
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Join(tt.args.futures)
			if !tt.wantErr(t, err, fmt.Sprintf("Join(%v)", tt.args.futures)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Join(%v)", tt.args.futures)
		})
	}
}
