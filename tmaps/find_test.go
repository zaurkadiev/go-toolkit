package tmaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFirst(t *testing.T) {
	type args[K comparable, V any] struct {
		m  map[K]V
		fn func(K, V) bool
	}
	type testCase[K comparable, V any] struct {
		name  string
		args  args[K, V]
		want  K
		want1 V
		want2 bool
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_map_when_find_first_then_zero_values",
			args: args[int, int]{
				m:  nil,
				fn: func(_, _ int) bool { return true },
			},
			want:  0,
			want1: 0,
			want2: false,
		},
		{
			name: "given_empty_map_when_find_first_then_zero_values",
			args: args[int, int]{
				m:  map[int]int{},
				fn: func(_, _ int) bool { return true },
			},
			want:  0,
			want1: 0,
			want2: false,
		},
		{
			name: "given_single_pair_when_predicate_matches_then_return_pair",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
				},
				fn: func(k, v int) bool { return k == 1 && v == 10 },
			},
			want:  1,
			want1: 10,
			want2: true,
		},
		{
			name: "given_single_pair_when_predicate_misses_then_zero_values",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
				},
				fn: func(k, v int) bool { return k == 2 && v == 10 },
			},
			want:  0,
			want1: 0,
			want2: false,
		},
		{
			name: "given_multiple_pairs_when_only_one_matches_then_return_matching",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(k, v int) bool { return k == 2 && v == 20 },
			},
			want:  2,
			want1: 20,
			want2: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := FindFirst(tt.args.m, tt.args.fn)
			assert.Equalf(t, tt.want, got, "FindFirst() = %v, want %v", got, tt.want)
			assert.Equalf(t, tt.want1, got1, "FindFirst() = %v, want %v", got, tt.want)
			assert.Equalf(t, tt.want2, got2, "FindFirst() = %v, want %v", got, tt.want)
		})
	}
}

func TestContains(t *testing.T) {
	type args[K comparable, V any] struct {
		m  map[K]V
		fn func(K, V) bool
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want bool
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_map_when_contains_then_false",
			args: args[int, int]{
				m:  nil,
				fn: func(_, _ int) bool { return true },
			},
			want: false,
		},
		{
			name: "given_empty_map_when_contains_then_false",
			args: args[int, int]{
				m:  map[int]int{},
				fn: func(_, _ int) bool { return true },
			},
			want: false,
		},
		{
			name: "given_single_pair_when_predicate_matches_then_true",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
				},
				fn: func(k, v int) bool { return k == 1 && v == 10 },
			},
			want: true,
		},
		{
			name: "given_single_pair_when_predicate_misses_then_false",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
				},
				fn: func(k, v int) bool { return k == 2 && v == 10 },
			},
			want: false,
		},
		{
			name: "given_multiple_pairs_when_any_matches_by_key_then_true",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(k, _ int) bool { return k == 2 },
			},
			want: true,
		},
		{
			name: "given_multiple_pairs_when_any_matches_by_value_then_true",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(_, v int) bool { return v == 30 },
			},
			want: true,
		},
		{
			name: "given_multiple_pairs_when_none_matches_then_false",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(k, v int) bool { return k == 4 && v == 40 },
			},
			want: false,
		},
		{
			name: "given_multiple_pairs_when_predicate_matches_key_and_value_then_true",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(k, v int) bool { return k == 3 && v == 30 },
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.m, tt.args.fn); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
