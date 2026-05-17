package tslices

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	type args[T any, K comparable, V any] struct {
		slice []T
		fn    func(T) (K, V, bool)
	}
	type testCase[T any, K comparable, V any] struct {
		name string
		args args[T, K, V]
		want map[K][]V
	}
	tests := []testCase[int, int, int]{
		{
			name: "when_slice_is_empty_then_return_nil",
			args: args[int, int, int]{
				slice: []int{},
				fn: func(value int) (int, int, bool) {
					return value, value, true
				},
			},
			want: nil,
		},
		{
			name: "when_all_filtered_then_return_empty_map",
			args: args[int, int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) (int, int, bool) {
					return 0, 0, false
				},
			},
			want: map[int][]int{},
		},
		{
			name: "when_grouping_then_collect_values",
			args: args[int, int, int]{
				slice: []int{1, 2, 3, 4, 5},
				fn: func(value int) (int, int, bool) {
					return value % 2, value, true
				},
			},
			want: map[int][]int{
				0: {2, 4},
				1: {1, 3, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Group(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupFuncAsIs(t *testing.T) {
	type args[T any] struct {
		t T
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		want  T
		want1 T
		want2 bool
	}
	tests := []testCase[int]{
		{
			name: "when_value_is_zero_then_return_zero_twice",
			args: args[int]{
				t: 0,
			},
			want:  0,
			want1: 0,
			want2: true,
		},
		{
			name: "when_value_is_non_zero_then_return_value_twice",
			args: args[int]{
				t: 7,
			},
			want:  7,
			want1: 7,
			want2: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := GroupFuncAsIs(tt.args.t)
			assert.Equalf(t, tt.want, got, "GroupFuncAsIs(%v)", tt.args.t)
			assert.Equalf(t, tt.want1, got1, "GroupFuncAsIs(%v)", tt.args.t)
			assert.Equalf(t, tt.want2, got2, "GroupFuncAsIs(%v)", tt.args.t)
		})
	}
}

func TestGroupErr(t *testing.T) {
	type args[T any, K comparable, V any] struct {
		slice []T
		fn    func(T) (K, V, bool, error)
	}
	type testCase[T any, K comparable, V any] struct {
		name    string
		args    args[T, K, V]
		want    map[K][]V
		wantErr bool
	}
	tests := []testCase[int, int, int]{
		{
			name: "when_slice_is_empty_then_return_nil_without_error",
			args: args[int, int, int]{
				slice: []int{},
				fn: func(value int) (int, int, bool, error) {
					return value, value, true, nil
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "when_fn_returns_error_then_fail",
			args: args[int, int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) (int, int, bool, error) {
					if value == 2 {
						return 0, 0, false, assert.AnError
					}
					return value, value, true, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when_all_filtered_then_return_empty_map",
			args: args[int, int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) (int, int, bool, error) {
					return 0, 0, false, nil
				},
			},
			want:    map[int][]int{},
			wantErr: false,
		},
		{
			name: "when_grouping_then_collect_values",
			args: args[int, int, int]{
				slice: []int{1, 2, 3, 4, 5},
				fn: func(value int) (int, int, bool, error) {
					return value % 2, value, true, nil
				},
			},
			want: map[int][]int{
				0: {2, 4},
				1: {1, 3, 5},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GroupErr(tt.args.slice, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GroupErr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupErr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupErrFuncAsIs(t *testing.T) {
	type args[T any] struct {
		t T
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		want    T
		want1   T
		want2   bool
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[int]{
		{
			name: "when_value_is_zero_then_return_zero_twice",
			args: args[int]{
				t: 0,
			},
			want:    0,
			want1:   0,
			want2:   true,
			wantErr: assert.NoError,
		},
		{
			name: "when_value_is_non_zero_then_return_value_twice",
			args: args[int]{
				t: 7,
			},
			want:    7,
			want1:   7,
			want2:   true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := GroupErrFuncAsIs(tt.args.t)
			if !tt.wantErr(t, err, fmt.Sprintf("GroupErrFuncAsIs(%v)", tt.args.t)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GroupErrFuncAsIs(%v)", tt.args.t)
			assert.Equalf(t, tt.want1, got1, "GroupErrFuncAsIs(%v)", tt.args.t)
			assert.Equalf(t, tt.want2, got2, "GroupErrFuncAsIs(%v)", tt.args.t)
		})
	}
}
