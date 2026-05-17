package tslices

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	type args[T any, R ordered] struct {
		slice []T
		fn    func(T) R
	}
	type testCase[T any, R ordered] struct {
		name string
		args args[T, R]
		want R
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_empty_then_return_zero_value",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) int {
					return value
				},
			},
			want: 0,
		},
		{
			name: "when_all_values_mapped_then_return_sum",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) int {
					return value
				},
			},
			want: 6,
		},
		{
			name: "when_mapped_values_change_then_sum_mapped",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) int {
					return value * 2
				},
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumFuncAsIs(t *testing.T) {
	type args[T ordered] struct {
		t T
	}
	type testCase[T ordered] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "when_value_is_zero_then_return_zero",
			args: args[int]{
				t: 0,
			},
			want: 0,
		},
		{
			name: "when_value_is_non_zero_then_return_value",
			args: args[int]{
				t: 4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SumFuncAsIs(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumFuncAsIs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args[T any, R ordered] struct {
		slice []T
		fn    func(T) R
	}
	type testCase[T any, R ordered] struct {
		name string
		args args[T, R]
		want R
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_empty_then_return_zero_value",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) int {
					return value
				},
			},
			want: 0,
		},
		{
			name: "when_single_value_then_return_mapped_value",
			args: args[int, int]{
				slice: []int{2},
				fn: func(value int) int {
					return value * 2
				},
			},
			want: 4,
		},
		{
			name: "when_multiple_values_then_return_max",
			args: args[int, int]{
				slice: []int{1, 3, 2},
				fn: func(value int) int {
					return value
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxFuncAsIs(t *testing.T) {
	type args[T ordered] struct {
		t T
	}
	type testCase[T ordered] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "when_value_is_zero_then_return_zero",
			args: args[int]{
				t: 0,
			},
			want: 0,
		},
		{
			name: "when_value_is_non_zero_then_return_value",
			args: args[int]{
				t: 9,
			},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxFuncAsIs(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaxFuncAsIs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args[T any, R ordered] struct {
		slice []T
		fn    func(T) R
	}
	type testCase[T any, R ordered] struct {
		name string
		args args[T, R]
		want R
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_empty_then_return_zero_value",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) int {
					return value
				},
			},
			want: 0,
		},
		{
			name: "when_single_value_then_return_mapped_value",
			args: args[int, int]{
				slice: []int{2},
				fn: func(value int) int {
					return value * 2
				},
			},
			want: 4,
		},
		{
			name: "when_multiple_values_then_return_min",
			args: args[int, int]{
				slice: []int{3, 1, 2},
				fn: func(value int) int {
					return value
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinFuncAsIs(t *testing.T) {
	type args[T ordered] struct {
		t T
	}
	type testCase[T ordered] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "when_value_is_zero_then_return_zero",
			args: args[int]{
				t: 0,
			},
			want: 0,
		},
		{
			name: "when_value_is_non_zero_then_return_value",
			args: args[int]{
				t: 11,
			},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinFuncAsIs(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MinFuncAsIs() = %v, want %v", got, tt.want)
			}
		})
	}
}
