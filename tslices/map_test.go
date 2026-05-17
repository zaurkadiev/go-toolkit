package tslices

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	type args[T any, R any] struct {
		slice []T
		fn    func(T) (R, bool)
	}
	type testCase[T any, R any] struct {
		name string
		args args[T, R]
		want []R
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_nil_then_return_nil",
			args: args[int, int]{
				slice: nil,
				fn: func(value int) (int, bool) {
					return value, true
				},
			},
			want: nil,
		},
		{
			name: "when_slice_is_empty_then_return_empty_slice",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) (int, bool) {
					return value, true
				},
			},
			want: []int{},
		},
		{
			name: "when_filter_and_map_then_get_doubled_evens",
			args: args[int, int]{
				slice: []int{1, 2, 3, 4},
				fn: func(value int) (int, bool) {
					if value%2 == 0 {
						return value * 2, true
					}
					return 0, false
				},
			},
			want: []int{4, 8},
		},
		{
			name: "when_all_filtered_then_return_empty_slice",
			args: args[int, int]{
				slice: []int{1, 3, 5},
				fn: func(value int) (int, bool) {
					return value, false
				},
			},
			want: []int{},
		},
		{
			name: "when_all_mapped_then_return_all_values",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) (int, bool) {
					return value * 3, true
				},
			},
			want: []int{3, 6, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapErr(t *testing.T) {
	type args[T any, R any] struct {
		slice []T
		fn    func(T) (R, bool, error)
	}
	type testCase[T any, R any] struct {
		name    string
		args    args[T, R]
		want    []R
		wantErr bool
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_nil_then_return_nil_slice_without_error",
			args: args[int, int]{
				slice: nil,
				fn: func(value int) (int, bool, error) {
					return value, true, nil
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "when_slice_is_empty_then_return_empty_slice_without_error",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) (int, bool, error) {
					return value, true, nil
				},
			},
			want:    []int{},
			wantErr: false,
		},
		{
			name: "when_fn_returns_error_then_fail",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) (int, bool, error) {
					if value == 2 {
						return 0, false, assert.AnError
					}
					return value * 2, true, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when_all_filtered_out_then_return_empty",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) (int, bool, error) {
					return value, false, nil
				},
			},
			want:    []int{},
			wantErr: false,
		},
		{
			name: "when_some_filtered_then_return_only_mapped",
			args: args[int, int]{
				slice: []int{1, 2, 3, 4},
				fn: func(value int) (int, bool, error) {
					if value%2 == 0 {
						return value * 2, true, nil
					}
					return 0, false, nil
				},
			},
			want:    []int{4, 8},
			wantErr: false,
		},
		{
			name: "when_all_mapped_then_return_all_values",
			args: args[int, int]{
				slice: []int{1, 2},
				fn: func(value int) (int, bool, error) {
					return value + 1, true, nil
				},
			},
			want:    []int{2, 3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapErr(tt.args.slice, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapErr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapErr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap2Many(t *testing.T) {
	type args[T any, R any] struct {
		slice []T
		fn    func(T) ([]R, bool)
	}
	type testCase[T any, R any] struct {
		name string
		args args[T, R]
		want []R
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_nil_then_return_nil",
			args: args[int, int]{
				slice: nil,
				fn: func(value int) ([]int, bool) {
					return []int{value}, true
				},
			},
			want: nil,
		},
		{
			name: "when_slice_is_empty_then_return_empty_slice",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) ([]int, bool) {
					return []int{value}, true
				},
			},
			want: []int{},
		},
		{
			name: "when_all_filtered_out_then_return_empty",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) ([]int, bool) {
					return []int{value}, false
				},
			},
			want: []int{},
		},
		{
			name: "when_some_filtered_then_flatten_results",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) ([]int, bool) {
					if value%2 == 0 {
						return []int{value, value * 2}, true
					}
					return nil, false
				},
			},
			want: []int{2, 4},
		},
		{
			name: "when_all_mapped_then_flatten_results",
			args: args[int, int]{
				slice: []int{1, 2},
				fn: func(value int) ([]int, bool) {
					return []int{value, value + 1}, true
				},
			},
			want: []int{1, 2, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map2Many(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map2Many() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap2ManyErr(t *testing.T) {
	type args[T any, R any] struct {
		slice []T
		fn    func(T) ([]R, bool, error)
	}
	type testCase[T any, R any] struct {
		name    string
		args    args[T, R]
		want    []R
		wantErr bool
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_nil_then_return_nil_without_error",
			args: args[int, int]{
				slice: nil,
				fn: func(value int) ([]int, bool, error) {
					return []int{value}, true, nil
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "when_slice_is_empty_then_return_empty_without_error",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) ([]int, bool, error) {
					return []int{value}, true, nil
				},
			},
			want:    []int{},
			wantErr: false,
		},
		{
			name: "when_fn_returns_error_then_fail",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) ([]int, bool, error) {
					if value == 2 {
						return nil, false, assert.AnError
					}
					return []int{value}, true, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when_all_filtered_out_then_return_empty",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) ([]int, bool, error) {
					return []int{value}, false, nil
				},
			},
			want:    []int{},
			wantErr: false,
		},
		{
			name: "when_some_filtered_then_flatten_results",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) ([]int, bool, error) {
					if value%2 == 0 {
						return []int{value, value * 2}, true, nil
					}
					return nil, false, nil
				},
			},
			want:    []int{2, 4},
			wantErr: false,
		},
		{
			name: "when_all_mapped_then_flatten_results",
			args: args[int, int]{
				slice: []int{1, 2},
				fn: func(value int) ([]int, bool, error) {
					return []int{value, value + 1}, true, nil
				},
			},
			want:    []int{1, 2, 2, 3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Map2ManyErr(tt.args.slice, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Map2ManyErr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map2ManyErr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
