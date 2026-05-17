package tslices

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	type args[T any] struct {
		slice []T
		size  int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want [][]T
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_empty",
			args: args[int]{
				slice: []int{},
				size:  2,
			},
			want: [][]int{},
		},
		{
			name: "when_size_is_zero_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				size:  0,
			},
			want: [][]int{},
		},
		{
			name: "when_size_is_negative_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				size:  -1,
			},
			want: [][]int{},
		},
		{
			name: "when_size_is_one_then_split_each_element",
			args: args[int]{
				slice: []int{1, 2, 3},
				size:  1,
			},
			want: [][]int{{1}, {2}, {3}},
		},
		{
			name: "when_size_equals_length_then_single_chunk",
			args: args[int]{
				slice: []int{1, 2, 3},
				size:  3,
			},
			want: [][]int{{1, 2, 3}},
		},
		{
			name: "when_size_exceeds_length_then_single_chunk",
			args: args[int]{
				slice: []int{1, 2, 3},
				size:  5,
			},
			want: [][]int{{1, 2, 3}},
		},
		{
			name: "when_size_is_two_then_split_into_pairs",
			args: args[int]{
				slice: []int{1, 2, 3, 4, 5},
				size:  2,
			},
			want: [][]int{{1, 2}, {3, 4}, {5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Chunk(tt.args.slice, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSet(t *testing.T) {
	type args[T comparable, R comparable] struct {
		slice []T
		fn    func(T) R
	}
	type testCase[T comparable, R comparable] struct {
		name string
		args args[T, R]
		want map[R]bool
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_empty_then_return_empty_map",
			args: args[int, int]{
				slice: []int{},
				fn: func(value int) int {
					return value
				},
			},
			want: map[int]bool{},
		},
		{
			name: "when_slice_has_unique_values_then_set_contains_all",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				fn: func(value int) int {
					return value
				},
			},
			want: map[int]bool{
				1: true,
				2: true,
				3: true,
			},
		},
		{
			name: "when_slice_has_duplicates_then_set_contains_unique",
			args: args[int, int]{
				slice: []int{1, 2, 2, 3, 1},
				fn: func(value int) int {
					return value
				},
			},
			want: map[int]bool{
				1: true,
				2: true,
				3: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSet(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSetFuncAsIs(t *testing.T) {
	type args[T any] struct {
		t T
	}
	type testCase[T any] struct {
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
				t: 8,
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSetFuncAsIs(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSetFuncAsIs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	type args[T any] struct {
		slice []T
	}
	type testCase[T comparable] struct {
		name  string
		args  args[T]
		want  map[T]int
		wantL int
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_empty",
			args: args[int]{
				slice: []int{},
			},
			want:  map[int]int{},
			wantL: 0,
		},
		{
			name: "when_slice_has_one_element_then_keep_element",
			args: args[int]{
				slice: []int{1},
			},
			want: map[int]int{
				1: 1,
			},
			wantL: 1,
		},
		{
			name: "when_slice_has_multiple_elements_then_keep_all_elements",
			args: args[int]{
				slice: []int{1, 2, 2, 3},
			},
			want: map[int]int{
				1: 1,
				2: 2,
				3: 1,
			},
			wantL: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Shuffle(tt.args.slice)
			if len(got) != tt.wantL {
				t.Errorf("Shuffle() len = %v, want %v", len(got), tt.wantL)
				return
			}
			gotCounts := make(map[int]int, len(got))
			for _, value := range got {
				gotCounts[value]++
			}
			if !reflect.DeepEqual(gotCounts, tt.want) {
				t.Errorf("Shuffle() = %v, want %v", gotCounts, tt.want)
			}
		})
	}
}

func TestFirstN(t *testing.T) {
	type args[T any] struct {
		slice []T
		n     int
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
				n:     1,
			},
			want: []int{},
		},
		{
			name: "when_n_is_zero_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     0,
			},
			want: []int{},
		},
		{
			name: "when_n_is_negative_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     -1,
			},
			want: []int{},
		},
		{
			name: "when_n_is_one_then_return_first",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     1,
			},
			want: []int{1},
		},
		{
			name: "when_n_equals_length_then_return_slice",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     3,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "when_n_exceeds_length_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     4,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "when_n_less_than_length_then_return_prefix",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     2,
			},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstN(tt.args.slice, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FirstN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastN(t *testing.T) {
	type args[T any] struct {
		slice []T
		n     int
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
				n:     1,
			},
			want: []int{},
		},
		{
			name: "when_n_is_zero_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     0,
			},
			want: []int{},
		},
		{
			name: "when_n_is_negative_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     -1,
			},
			want: []int{},
		},
		{
			name: "when_n_is_one_then_return_last",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     1,
			},
			want: []int{3},
		},
		{
			name: "when_n_equals_length_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     3,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "when_n_exceeds_length_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     4,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "when_n_less_than_length_then_return_suffix",
			args: args[int]{
				slice: []int{1, 2, 3},
				n:     2,
			},
			want: []int{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastN(tt.args.slice, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LastN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args[T any, R any] struct {
		slice []T
		a     R
		fn    func(R, T) R
	}
	type testCase[T any, R any] struct {
		name string
		args args[T, R]
		want R
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_empty_then_return_initial_value",
			args: args[int, int]{
				slice: []int{},
				a:     10,
				fn: func(acc, value int) int {
					return acc + value
				},
			},
			want: 10,
		},
		{
			name: "when_slice_has_values_then_reduce",
			args: args[int, int]{
				slice: []int{1, 2, 3},
				a:     0,
				fn: func(acc, value int) int {
					return acc + value
				},
			},
			want: 6,
		},
		{
			name: "when_initial_value_non_zero_then_used_in_result",
			args: args[int, int]{
				slice: []int{1, 2},
				a:     5,
				fn: func(acc, value int) int {
					return acc + value
				},
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.slice, tt.args.a, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByIndex(t *testing.T) {
	type args[T any] struct {
		slice []T
		index int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "when_index_is_negative_then_return_zero_value",
			args: args[int]{
				slice: []int{1, 2, 3},
				index: -1,
			},
			want: 0,
		},
		{
			name: "when_index_is_zero_then_return_first",
			args: args[int]{
				slice: []int{1, 2, 3},
				index: 0,
			},
			want: 1,
		},
		{
			name: "when_index_is_last_then_return_last",
			args: args[int]{
				slice: []int{1, 2, 3},
				index: 2,
			},
			want: 3,
		},
		{
			name: "when_index_exceeds_slice_len_then_return_empty",
			args: args[int]{
				slice: []int{1, 2, 3},
				index: 6,
			},
			want: 0,
		},
		{
			name: "when_index_out_of_range_then_return_zero_value",
			args: args[int]{
				slice: []int{1, 2, 3},
				index: 3,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byIndex(tt.args.slice, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("byIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZip(t *testing.T) {
	type args[L any, R any] struct {
		left  []L
		right []R
	}
	type testCase[L any, R any] struct {
		name string
		args args[L, R]
		want []Pair[L, R]
	}
	tests := []testCase[int, int]{
		{
			name: "when_both_slices_empty_then_return_empty",
			args: args[int, int]{
				left:  []int{},
				right: []int{},
			},
			want: []Pair[int, int]{},
		},
		{
			name: "when_left_empty_then_pad_left_with_zero",
			args: args[int, int]{
				left:  []int{},
				right: []int{1, 2},
			},
			want: []Pair[int, int]{
				{Left: 0, Right: 1},
				{Left: 0, Right: 2},
			},
		},
		{
			name: "when_right_empty_then_pad_right_with_zero",
			args: args[int, int]{
				left:  []int{1, 2},
				right: []int{},
			},
			want: []Pair[int, int]{
				{Left: 1, Right: 0},
				{Left: 2, Right: 0},
			},
		},
		{
			name: "when_slices_equal_length_then_zip_pairs",
			args: args[int, int]{
				left:  []int{1, 2},
				right: []int{3, 4},
			},
			want: []Pair[int, int]{
				{Left: 1, Right: 3},
				{Left: 2, Right: 4},
			},
		},
		{
			name: "when_lengths_differ_then_pad_shorter",
			args: args[int, int]{
				left:  []int{1},
				right: []int{2, 3},
			},
			want: []Pair[int, int]{
				{Left: 1, Right: 2},
				{Left: 0, Right: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Zip(tt.args.left, tt.args.right); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Zip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduceFuncSum(t *testing.T) {
	type args[A ordered] struct {
		t A
		a A
	}
	type testCase[A ordered] struct {
		name string
		args args[A]
		want A
	}
	tests := []testCase[int]{
		{
			name: "when_a_and_t_are_set_return_its_sum",
			args: args[int]{
				t: 10,
				a: 0,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ReduceFuncSum(tt.args.t, tt.args.a), "ReduceFuncSum(%v, %v)", tt.args.t, tt.args.a)
		})
	}
}

func TestReduceErr(t *testing.T) {
	boom := errors.New("boom")
	tests := []struct {
		name    string
		slice   []int
		fn      func(int, int) (int, error)
		want    int
		wantErr error
	}{
		{
			name:    "empty_slice_returns_initial_no_error",
			slice:   nil,
			fn:      func(a, t int) (int, error) { return a + t, nil },
			want:    0,
			wantErr: nil,
		},
		{
			name:    "sums_all_when_no_error",
			slice:   []int{1, 2, 3},
			fn:      func(a, t int) (int, error) { return a + t, nil },
			want:    6,
			wantErr: nil,
		},
		{
			name:  "returns_error_and_stops",
			slice: []int{1, 2, 3},
			fn: func(a, t int) (int, error) {
				if t == 2 {
					return a, boom
				}
				return a + t, nil
			},
			want:    1,
			wantErr: boom,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReduceErr(tt.slice, 0, tt.fn)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
