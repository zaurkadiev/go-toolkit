package tslices

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	type args[T comparable, R comparable] struct {
		slice []T
		fn    func(T) R
	}
	type testCase[T comparable, R comparable] struct {
		name string
		args args[T, R]
		want []T
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_empty_then_return_empty",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) int {
					return value
				},
			},
			want: []int{},
		},
		{
			name: "when_slice_has_unique_values_then_return_all",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) int {
					return value
				},
			},
			want: []int{1, 3, 2},
		},
		{
			name: "when_slice_has_duplicates_then_return_unique",
			args: args[int, int]{
				slice: []int{1, 2, 2, 3, 1},
				fn: func(value int) int {
					return value
				},
			},
			want: []int{1, 3, 2},
		},
		{
			name: "when_key_groups_values_then_keep_last_per_key",
			args: args[int, int]{
				slice: []int{1, 2, 3, 4, 4},
				fn: func(value int) int {
					return value % 2
				},
			},
			want: []int{4, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.args.slice, tt.args.fn)
			assert.ElementsMatchf(t, tt.want, got, "UniqueOrdered(%v, %T)", tt.args.slice, tt.args.fn)
		})
	}
}

func TestUniqueFuncAsIs(t *testing.T) {
	type args[T comparable] struct {
		value T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "when_value_is_zero_then_return_zero",
			args: args[int]{
				value: 0,
			},
			want: 0,
		},
		{
			name: "when_value_is_non_zero_then_return_value",
			args: args[int]{
				value: 12,
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueFuncAsIs(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueFuncAsIs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueOrdered(t *testing.T) {
	type args[T comparable, R comparable] struct {
		slice []T
		fn    func(T) R
	}
	type testCase[T comparable, R comparable] struct {
		name string
		args args[T, R]
		want []T
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_empty_then_return_empty",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) int {
					return value
				},
			},
			want: []int{},
		},
		{
			name: "when_slice_has_unique_values_then_return_all",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) int {
					return value
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "when_slice_has_duplicates_then_return_unique",
			args: args[int, int]{
				slice: []int{1, 2, 2, 3, 1},
				fn: func(value int) int {
					return value
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "when_key_groups_values_then_keep_last_first_per_key",
			args: args[int, int]{
				slice: []int{1, 2, 3, 4},
				fn: func(value int) int {
					return value % 2
				},
			},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, UniqueOrdered(tt.args.slice, tt.args.fn), "UniqueOrdered(%v, %v)", tt.args.slice, tt.args.fn)
		})
	}
}
