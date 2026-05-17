package tslices

import (
	"reflect"
	"testing"
)

func TestIndexOf(t *testing.T) {
	type args[T any] struct {
		slice []T
		fn    func(T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_minus_one",
			args: args[int]{
				slice: []int{},
				fn: func(value int) bool {
					return value == 1
				},
			},
			want: -1,
		},
		{
			name: "when_no_match_then_return_minus_one",
			args: args[int]{
				slice: []int{1, 2, 3},
				fn: func(value int) bool {
					return value == 4
				},
			},
			want: -1,
		},
		{
			name: "when_match_exists_then_return_first_index",
			args: args[int]{
				slice: []int{1, 2, 3, 2},
				fn: func(value int) bool {
					return value == 2
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexOf(tt.args.slice, tt.args.fn); got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args[T any] struct {
		slice []T
		fn    func(T) bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_false",
			args: args[int]{
				slice: []int{},
				fn: func(value int) bool {
					return value == 1
				},
			},
			want: false,
		},
		{
			name: "when_value_matches_then_return_true",
			args: args[int]{
				slice: []int{1, 2, 3},
				fn: func(value int) bool {
					return value == 2
				},
			},
			want: true,
		},
		{
			name: "when_no_values_match_then_return_false",
			args: args[int]{
				slice: []int{1, 2, 3},
				fn: func(value int) bool {
					return value == 4
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.slice, tt.args.fn); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindFirst(t *testing.T) {
	type args[T any] struct {
		slice []T
		fn    func(T) bool
	}
	type testCase[T any] struct {
		name      string
		args      args[T]
		want      T
		wantIndex int
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_zero_and_-1",
			args: args[int]{
				slice: []int{},
				fn: func(value int) bool {
					return value == 1
				},
			},
			want:      0,
			wantIndex: -1,
		},
		{
			name: "when_no_match_then_return_zero_and_-1",
			args: args[int]{
				slice: []int{1, 2, 3},
				fn: func(value int) bool {
					return value == 4
				},
			},
			want:      0,
			wantIndex: -1,
		},
		{
			name: "when_match_exists_then_return_first_match",
			args: args[int]{
				slice: []int{1, 2, 3, 2},
				fn: func(value int) bool {
					return value == 2
				},
			},
			want:      2,
			wantIndex: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindFirst(tt.args.slice, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindFirst() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantIndex {
				t.Errorf("FindFirst() got1 = %v, want %v", got1, tt.wantIndex)
			}
		})
	}
}

func TestFindLast(t *testing.T) {
	type args[T any] struct {
		slice []T
		fn    func(T) bool
	}
	type testCase[T any] struct {
		name      string
		args      args[T]
		want      T
		wantIndex int
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_zero_and_false",
			args: args[int]{
				slice: []int{},
				fn: func(value int) bool {
					return value == 1
				},
			},
			want:      0,
			wantIndex: -1,
		},
		{
			name: "when_no_match_then_return_zero_and_false",
			args: args[int]{
				slice: []int{1, 2, 3},
				fn: func(value int) bool {
					return value == 4
				},
			},
			want:      0,
			wantIndex: -1,
		},
		{
			name: "when_match_exists_then_return_last_match",
			args: args[int]{
				slice: []int{1, 2, 3, 2},
				fn: func(value int) bool {
					return value == 2
				},
			},
			want:      2,
			wantIndex: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := FindLast(tt.args.slice, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindLast() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantIndex {
				t.Errorf("FindLast() got1 = %v, want %v", got1, tt.wantIndex)
			}
		})
	}
}
