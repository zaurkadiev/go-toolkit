package tmaps

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	type args[K comparable, V any] struct {
		m  map[K]V
		fn func(K, V) bool
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_map_when_filter_then_nil",
			args: args[int, int]{
				m:  nil,
				fn: func(_, _ int) bool { return true },
			},
			want: nil,
		},
		{
			name: "given_empty_map_when_filter_then_nil",
			args: args[int, int]{
				m:  map[int]int{},
				fn: func(_, _ int) bool { return true },
			},
			want: map[int]int{},
		},
		{
			name: "given_map_when_predicate_always_true_then_same_pairs",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
				},
				fn: func(_, _ int) bool { return true },
			},
			want: map[int]int{
				1: 10,
				2: 20,
			},
		},
		{
			name: "given_map_when_predicate_always_false_then_empty_map",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
				},
				fn: func(_, _ int) bool { return false },
			},
			want: map[int]int{},
		},
		{
			name: "given_map_when_filter_by_key_then_only_matching",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(k, _ int) bool { return k%2 == 1 },
			},
			want: map[int]int{
				1: 10,
				3: 30,
			},
		},
		{
			name: "given_map_when_filter_by_value_then_only_matching",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(_, v int) bool { return v >= 20 },
			},
			want: map[int]int{
				2: 20,
				3: 30,
			},
		},
		{
			name: "given_map_when_filter_by_key_and_value_then_only_matching",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(k, v int) bool { return k == 2 && v == 20 },
			},
			want: map[int]int{
				2: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.m, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterFuncNotZeroNotPointerValue(t *testing.T) {
	type args[K any, V comparable] struct {
		k K
		v V
	}
	type testCase[K any, V comparable] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[int, int]{
		{
			name: "when_value_type_is_not_pointer_and_zero_then_return_false",
			args: args[int, int]{
				k: 0,
				v: 0, // is zero
			},
			want: false,
		},
		{
			name: "when_value_type_is_not_pointer_and_not_zero_then_return_true",
			args: args[int, int]{
				k: 0,
				v: 1, // is not zero
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FilterFuncNotZero(tt.args.k, tt.args.v), "FilterFuncNotZero(%v, %v)", tt.args.k, tt.args.v)
		})
	}
}

func TestFilterFuncNotZeroPointerValue(t *testing.T) {
	type args[K any, V comparable] struct {
		k K
		v V
	}
	type testCase[K any, V comparable] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[int, *int]{
		{
			name: "when_value_type_is_not_pointer_and_zero_then_return_false",
			args: args[int, *int]{
				k: 0,
				v: nil, // is zero
			},
			want: false,
		},
		{
			name: "when_value_type_is_not_pointer_and_not_zero_then_return_true",
			args: args[int, *int]{
				k: 0,
				v: func() *int {
					v := 1
					return &v
				}(), // is zero
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FilterFuncNotZero(tt.args.k, tt.args.v), "FilterFuncNotZero(%v, %v)", tt.args.k, tt.args.v)
		})
	}
}
