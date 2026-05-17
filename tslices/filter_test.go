package tslices

import (
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	type args[T any] struct {
		slice []T
		fn    func(T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_empty",
			args: args[int]{
				slice: []int{},
				fn: func(value int) bool {
					return value%2 == 0
				},
			},
			want: []int{},
		},
		{
			name: "when_all_values_filtered_out_then_return_empty",
			args: args[int]{
				slice: []int{1, 3, 5},
				fn: func(value int) bool {
					return value%2 == 0
				},
			},
			want: []int{},
		},
		{
			name: "when_all_values_match_then_return_all",
			args: args[int]{
				slice: []int{2, 4, 6},
				fn: func(value int) bool {
					return value%2 == 0
				},
			},
			want: []int{2, 4, 6},
		},
		{
			name: "when_some_values_match_then_return_filtered",
			args: args[int]{
				slice: []int{1, 2, 3, 4},
				fn: func(value int) bool {
					return value%2 == 0
				},
			},
			want: []int{2, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterFuncNotEmpty(t *testing.T) {
	type args[T comparable] struct {
		t T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "when_value_is_zero_then_return_false",
			args: args[int]{
				t: 0,
			},
			want: false,
		},
		{
			name: "when_value_is_non_zero_then_return_true",
			args: args[int]{
				t: 5,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterFuncNotZero(tt.args.t); got != tt.want {
				t.Errorf("FilterFuncNotEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
