package worker_pool

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		n int
	}
	type testCase[T any] struct {
		name string
		args args
		want WorkerPool[T]
	}
	tests := []testCase[int]{
		{
			name: "given_positive_size_when_new_then_empty_pool",
			args: args{
				n: 2,
			},
			want: &impl[int]{
				jobs:     []func() (int, error){},
				poolSize: 2,
			},
		},
		{
			name: "given_zero_size_when_new_then_empty_pool",
			args: args{
				n: 0,
			},
			want: &impl[int]{
				jobs:     []func() (int, error){},
				poolSize: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, New[int](tt.args.n), "New(%v)", tt.args.n)
		})
	}
}

func Test_impl_Add(t *testing.T) {
	type fields[T any] struct {
		jobs []func() (T, error)
	}
	type args[T any] struct {
		job func() (T, error)
	}
	type testCase[T any] struct {
		name   string
		fields fields[T]
		args   args[T]
		want   int
	}
	tests := []testCase[int]{
		{
			name: "given_no_jobs_when_add_then_one_job",
			fields: fields[int]{
				jobs: []func() (int, error){},
			},
			args: args[int]{
				job: func() (int, error) { return 1, nil },
			},
			want: 1,
		},
		{
			name: "given_jobs_when_add_then_job_appended",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) { return 1, nil },
				},
			},
			args: args[int]{
				job: func() (int, error) { return 2, nil },
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp := &impl[int]{
				jobs: tt.fields.jobs,
			}
			wp.Add(tt.args.job)
			assert.Lenf(t, wp.jobs, tt.want, "Add(%T)", tt.args.job)
		})
	}
}

func Test_impl_Run(t *testing.T) {
	type fields[T any] struct {
		jobs     []func() (T, error)
		poolSize int
	}
	type args struct {
		ctx context.Context
	}
	type testCase[T any] struct {
		name    string
		fields  fields[T]
		args    args
		want    []T
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[int]{
		{
			name: "given_no_jobs_when_run_then_nil",
			fields: fields[int]{
				jobs:     []func() (int, error){},
				poolSize: 1,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "given_zero_pool_size_when_run_then_error",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) { return 1, nil },
				},
				poolSize: 0,
			},
			args: args{
				ctx: context.Background(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.EqualError(t, err, "worker size must be positive", i...)
			},
		},
		{
			name: "given_cancelled_context_when_run_then_context_error",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) { return 1, nil },
				},
				poolSize: 1,
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, context.Canceled, i...)
			},
		},
		{
			name: "given_job_error_when_run_then_error",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) { return 1, nil },
					func() (int, error) { return 0, assert.AnError },
				},
				poolSize: 2,
			},
			args: args{
				ctx: context.Background(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
		{
			name: "given_job_error_with_queued_jobs_when_run_then_queued_jobs_not_started",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) { return 0, assert.AnError },
					func() (int, error) {
						t.Error("job after error must not start")
						return 2, nil
					},
				},
				poolSize: 1,
			},
			args: args{
				ctx: context.Background(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
		{
			name: "given_job_panics_when_run_then_error",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) { return 1, nil },
					func() (int, error) { panic("panic") },
				},
				poolSize: 2,
			},
			args: args{
				ctx: context.Background(),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.EqualError(t, err, "panic in worker pool job: panic", i...)
			},
		},
		{
			name: "given_single_worker_when_run_then_results",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) { return 1, nil },
					func() (int, error) { return 2, nil },
				},
				poolSize: 1,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []int{1, 2},
			wantErr: assert.NoError,
		},
		{
			name: "given_more_workers_than_jobs_when_run_then_results",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) { return 1, nil },
				},
				poolSize: 4,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []int{1},
			wantErr: assert.NoError,
		},
		{
			name: "given_jobs_finish_in_reverse_order_when_run_then_results_in_add_order",
			fields: fields[int]{
				jobs: []func() (int, error){
					func() (int, error) {
						time.Sleep(30 * time.Millisecond)
						return 1, nil
					},
					func() (int, error) {
						time.Sleep(20 * time.Millisecond)
						return 2, nil
					},
					func() (int, error) {
						time.Sleep(10 * time.Millisecond)
						return 3, nil
					},
					func() (int, error) { return 4, nil },
				},
				poolSize: 4,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []int{1, 2, 3, 4},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp := &impl[int]{
				jobs:     tt.fields.jobs,
				poolSize: tt.fields.poolSize,
			}
			got, err := wp.Run(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("Run(%v)", tt.args.ctx)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Run(%v)", tt.args.ctx)
		})
	}
}

func TestNewSimple(t *testing.T) {
	type args struct {
		n int
	}
	type testCase struct {
		name string
		args args
		want Simple
	}
	tests := []testCase{
		{
			name: "given_positive_size_when_new_simple_then_empty_pool",
			args: args{
				n: 2,
			},
			want: &simpleImpl{
				jobs:     []func() error{},
				poolSize: 2,
			},
		},
		{
			name: "given_zero_size_when_new_simple_then_empty_pool",
			args: args{
				n: 0,
			},
			want: &simpleImpl{
				jobs:     []func() error{},
				poolSize: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewSimple(tt.args.n), "NewSimple(%v)", tt.args.n)
		})
	}
}

func Test_simpleImpl_Add(t *testing.T) {
	type fields struct {
		jobs []func() error
	}
	type args struct {
		jobs []func() error
	}
	type testCase struct {
		name   string
		fields fields
		args   args
		want   int
	}
	tests := []testCase{
		{
			name: "given_no_jobs_when_add_none_then_no_jobs",
			fields: fields{
				jobs: []func() error{},
			},
			args: args{
				jobs: nil,
			},
			want: 0,
		},
		{
			name: "given_no_jobs_when_add_one_then_one_job",
			fields: fields{
				jobs: []func() error{},
			},
			args: args{
				jobs: []func() error{
					func() error { return nil },
				},
			},
			want: 1,
		},
		{
			name: "given_jobs_when_add_many_then_jobs_appended",
			fields: fields{
				jobs: []func() error{
					func() error { return nil },
				},
			},
			args: args{
				jobs: []func() error{
					func() error { return nil },
					func() error { return nil },
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp := &simpleImpl{
				jobs: tt.fields.jobs,
			}
			wp.Add(tt.args.jobs...)
			assert.Lenf(t, wp.jobs, tt.want, "Add(%v)", len(tt.args.jobs))
		})
	}
}

func Test_simpleImpl_Run(t *testing.T) {
	type fields struct {
		jobs     []func() error
		poolSize int
	}
	type args struct {
		ctx context.Context
	}
	type testCase struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase{
		{
			name: "given_no_jobs_when_run_then_nil",
			fields: fields{
				jobs:     []func() error{},
				poolSize: 1,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "given_zero_pool_size_when_run_then_error",
			fields: fields{
				jobs: []func() error{
					func() error { return nil },
				},
				poolSize: 0,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.EqualError(t, err, "worker size must be positive", i...)
			},
		},
		{
			name: "given_cancelled_context_when_run_then_context_error",
			fields: fields{
				jobs: []func() error{
					func() error { return nil },
				},
				poolSize: 1,
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, context.Canceled, i...)
			},
		},
		{
			name: "given_job_error_when_run_then_error",
			fields: fields{
				jobs: []func() error{
					func() error { return nil },
					func() error { return assert.AnError },
				},
				poolSize: 2,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
		{
			name: "given_job_error_with_queued_jobs_when_run_then_queued_jobs_not_started",
			fields: fields{
				jobs: []func() error{
					func() error { return assert.AnError },
					func() error {
						t.Error("job after error must not start")
						return nil
					},
				},
				poolSize: 1,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
		{
			name: "given_job_panics_when_run_then_error",
			fields: fields{
				jobs: []func() error{
					func() error { panic("panic") },
				},
				poolSize: 1,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.EqualError(t, err, "panic in worker pool job: panic", i...)
			},
		},
		{
			name: "given_jobs_succeed_when_run_then_nil",
			fields: fields{
				jobs: []func() error{
					func() error { return nil },
					func() error { return nil },
					func() error { return nil },
				},
				poolSize: 2,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp := &simpleImpl{
				jobs:     tt.fields.jobs,
				poolSize: tt.fields.poolSize,
			}
			tt.wantErr(t, wp.Run(tt.args.ctx), fmt.Sprintf("Run(%v)", tt.args.ctx))
		})
	}
}

func TestCall(t *testing.T) {
	type args[T any] struct {
		job func() (T, error)
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		want    T
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[int]{
		{
			name: "given_job_returns_value_when_call_then_value",
			args: args[int]{
				job: func() (int, error) { return 1, nil },
			},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name: "given_job_returns_error_when_call_then_error",
			args: args[int]{
				job: func() (int, error) { return 0, assert.AnError },
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
		{
			name: "given_job_panics_with_value_when_call_then_error",
			args: args[int]{
				job: func() (int, error) { panic("panic") },
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.EqualError(t, err, "panic in worker pool job: panic", i...)
			},
		},
		{
			name: "given_job_panics_with_error_when_call_then_wrapped_error",
			args: args[int]{
				job: func() (int, error) { panic(assert.AnError) },
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.ErrorIs(t, err, assert.AnError, i...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := call(tt.args.job)
			if !tt.wantErr(t, err, fmt.Sprintf("call(%T)", tt.args.job)) {
				return
			}
			assert.Equalf(t, tt.want, got, "call(%T)", tt.args.job)
		})
	}
}
