package tmaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeys(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []K
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_map_when_get_keys_then_empty_slice",
			args: args[int, int]{
				m: nil,
			},
			want: []int{},
		},
		{
			name: "given_empty_map_when_get_keys_then_empty_slice",
			args: args[int, int]{
				m: map[int]int{},
			},
			want: []int{},
		},
		{
			name: "given_map_when_get_keys_then_all_keys",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatchf(t, tt.want, GetKeys(tt.args.m), "GetVals(%v)", tt.args.m)
		})
	}
}

func TestGetVals(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want []V
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_map_when_get_vals_then_empty_slice",
			args: args[int, int]{
				m: nil,
			},
			want: []int{},
		},
		{
			name: "given_empty_map_when_get_vals_then_empty_slice",
			args: args[int, int]{
				m: map[int]int{},
			},
			want: []int{},
		},
		{
			name: "given_map_when_get_vals_then_all_values",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
			},
			want: []int{10, 20, 30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, GetVals(tt.args.m))
		})
	}
}

func TestGetUniqueVals(t *testing.T) {
	type args[K comparable, V comparable] struct {
		m map[K]V
	}
	type testCase[K comparable, V comparable] struct {
		name string
		args args[K, V]
		want []V
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_map_when_get_unique_vals_then_empty_slice",
			args: args[int, int]{
				m: nil,
			},
			want: []int{},
		},
		{
			name: "given_empty_map_when_get_unique_vals_then_empty_slice",
			args: args[int, int]{
				m: map[int]int{},
			},
			want: []int{},
		},
		{
			name: "given_map_with_unique_values_when_get_unique_vals_then_all_values",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
			},
			want: []int{10, 20, 30},
		},
		{
			name: "given_map_with_duplicate_values_when_get_unique_vals_then_deduped_values",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 10,
					3: 20,
				},
			},
			want: []int{10, 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, GetValsUnique(tt.args.m))
		})
	}
}

func TestGetKeysSorted(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]int
		want []string
	}{
		{
			name: "empty_map_returns_empty",
			m:    map[string]int{},
			want: []string{},
		},
		{
			name: "keys_returned_sorted",
			m:    map[string]int{"b": 2, "a": 1, "c": 3},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetKeysSorted(tt.m))
		})
	}
}
