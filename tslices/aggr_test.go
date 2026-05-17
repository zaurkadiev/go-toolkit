package tslices

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregate(t *testing.T) {
	type args[T any, K comparable, V any] struct {
		slice []T
		keyFn func(T) (K, bool)
		aggFn func([]T) (V, bool)
	}
	type testCase[T any, K comparable, V any] struct {
		name string
		args args[T, K, V]
		want map[K]V
	}
	tests := []testCase[int, int, int]{
		{
			name: "when_slice_is_nil_then_return_nil",
			args: args[int, int, int]{
				slice: nil,
				keyFn: func(value int) (int, bool) {
					return value, true
				},
				aggFn: func(values []int) (int, bool) {
					return len(values), true
				},
			},
			want: nil,
		},
		{
			name: "when_slice_is_empty_then_return_emtpy_map",
			args: args[int, int, int]{
				slice: []int{},
				keyFn: func(value int) (int, bool) {
					return value, true
				},
				aggFn: func(values []int) (int, bool) {
					return len(values), true
				},
			},
			want: map[int]int{},
		},
		{
			name: "when_all_keys_filtered_then_return_empty_map",
			args: args[int, int, int]{
				slice: []int{1, 2, 3},
				keyFn: func(value int) (int, bool) {
					return 0, false
				},
				aggFn: func(values []int) (int, bool) {
					return len(values), true
				},
			},
			want: map[int]int{},
		},
		{
			name: "when_agg_returns_false_then_skip_group_and_return_empty_map",
			args: args[int, int, int]{
				slice: []int{1, 2, 3, 4},
				keyFn: func(value int) (int, bool) {
					return value % 2, true
				},
				aggFn: func(values []int) (int, bool) {
					if len(values) == 2 {
						return 0, false
					}
					return len(values), true
				},
			},
			want: map[int]int{},
		},
		{
			name: "when_grouping_then_aggregate_per_key",
			args: args[int, int, int]{
				slice: []int{1, 2, 3, 4, 5},
				keyFn: func(value int) (int, bool) {
					return value % 2, true
				},
				aggFn: func(values []int) (int, bool) {
					sum := 0
					for _, value := range values {
						sum += value
					}
					return sum, true
				},
			},
			want: map[int]int{
				0: 6,
				1: 9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Aggregate(tt.args.slice, tt.args.keyFn, tt.args.aggFn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregateFuncKeyAsIs(t *testing.T) {
	type args[T any] struct {
		t T
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[int]{
		{
			name: "when_value_is_zero_then_return_zero_and_true",
			args: args[int]{
				t: 0,
			},
			want:  0,
			want1: true,
		},
		{
			name: "when_value_is_non_zero_then_return_value_and_true",
			args: args[int]{
				t: 5,
			},
			want:  5,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AggregateFuncKeyAsIs(tt.args.t)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateFuncKeyAsIs() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AggregateFuncKeyAsIs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAggregateFuncCount(t *testing.T) {
	type args[T any] struct {
		slice []T
	}
	type testCase[T any, R integer] struct {
		name  string
		args  args[T]
		want  R
		want1 bool
	}
	tests := []testCase[int, int]{
		{
			name: "when_slice_is_empty_then_return_zero_and_true",
			args: args[int]{
				slice: []int{},
			},
			want:  0,
			want1: true,
		},
		{
			name: "when_slice_has_values_then_return_length",
			args: args[int]{
				slice: []int{1, 2, 3},
			},
			want:  3,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AggregateFuncCount[int, int](tt.args.slice)
			if got != tt.want {
				t.Errorf("AggregateFuncCount() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AggregateFuncCount() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAggregateFuncSum(t *testing.T) {
	type args[T ordered] struct {
		slice []T
	}
	type testCase[T ordered] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_zero_and_false",
			args: args[int]{
				slice: []int{},
			},
			want:  0,
			want1: false,
		},
		{
			name: "when_slice_has_values_then_return_sum",
			args: args[int]{
				slice: []int{1, 2, 3},
			},
			want:  6,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AggregateFuncSum(tt.args.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateFuncSum() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AggregateFuncSum() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAggregateFuncMin(t *testing.T) {
	type args[T ordered] struct {
		slice []T
	}
	type testCase[T ordered] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_zero_and_false",
			args: args[int]{
				slice: []int{},
			},
			want:  0,
			want1: false,
		},
		{
			name: "when_slice_has_values_then_return_min",
			args: args[int]{
				slice: []int{3, 1, 2},
			},
			want:  1,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AggregateFuncMin(tt.args.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateFuncMin() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AggregateFuncMin() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAggregateFuncMax(t *testing.T) {
	type args[T ordered] struct {
		slice []T
	}
	type testCase[T ordered] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_zero_and_false",
			args: args[int]{
				slice: []int{},
			},
			want:  0,
			want1: false,
		},
		{
			name: "when_slice_has_values_then_return_max",
			args: args[int]{
				slice: []int{-1, 3, 1, 2},
			},
			want:  3,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AggregateFuncMax(tt.args.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateFuncMax() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AggregateFuncMax() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAggregateFuncFirst(t *testing.T) {
	type args[T any] struct {
		slice []T
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_zero_and_false",
			args: args[int]{
				slice: []int{},
			},
			want:  0,
			want1: false,
		},
		{
			name: "when_slice_has_values_then_return_first",
			args: args[int]{
				slice: []int{3, 2, 1},
			},
			want:  3,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AggregateFuncFirst(tt.args.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateFuncFirst() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AggregateFuncFirst() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAggregateFuncLast(t *testing.T) {
	type args[T any] struct {
		slice []T
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		want  T
		want1 bool
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_return_zero_and_false",
			args: args[int]{
				slice: []int{},
			},
			want:  0,
			want1: false,
		},
		{
			name: "when_slice_has_values_then_return_last",
			args: args[int]{
				slice: []int{3, 2, 1},
			},
			want:  1,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := AggregateFuncLast(tt.args.slice)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateFuncLast() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AggregateFuncLast() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAggregateErr(t *testing.T) {
	type args[T any, K comparable, V any] struct {
		slice []T
		keyFn func(T) (K, bool, error)
		aggFn func([]T) (V, bool, error)
	}
	type testCase[T any, K comparable, V any] struct {
		name    string
		args    args[T, K, V]
		want    map[K]V
		wantErr bool
	}
	tests := []testCase[int, int, int]{
		{
			name: "when_slice_is_nil_then_return_nil_without_error",
			args: args[int, int, int]{
				slice: nil,
				keyFn: func(value int) (int, bool, error) {
					return value, true, nil
				},
				aggFn: func(values []int) (int, bool, error) {
					return len(values), true, nil
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "when_slice_is_empty_then_return_empty_map",
			args: args[int, int, int]{
				slice: []int{},
				keyFn: func(value int) (int, bool, error) {
					return value, true, nil
				},
				aggFn: func(values []int) (int, bool, error) {
					return len(values), true, nil
				},
			},
			want:    map[int]int{},
			wantErr: false,
		},
		{
			name: "when_key_fn_returns_error_then_fail",
			args: args[int, int, int]{
				slice: []int{1, 2, 3},
				keyFn: func(value int) (int, bool, error) {
					if value == 2 {
						return 0, false, assert.AnError
					}
					return value, true, nil
				},
				aggFn: func(values []int) (int, bool, error) {
					return len(values), true, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when_agg_fn_returns_error_then_fail",
			args: args[int, int, int]{
				slice: []int{1, 2, 3, 4},
				keyFn: func(value int) (int, bool, error) {
					return value % 2, true, nil
				},
				aggFn: func(values []int) (int, bool, error) {
					if len(values) == 2 {
						return 0, false, assert.AnError
					}
					return len(values), true, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when_agg_returns_false_then_skip_group",
			args: args[int, int, int]{
				slice: []int{1, 2, 3, 4},
				keyFn: func(value int) (int, bool, error) {
					return value % 2, true, nil
				},
				aggFn: func(values []int) (int, bool, error) {
					if len(values) == 2 {
						return 0, false, nil
					}
					return len(values), true, nil
				},
			},
			want:    map[int]int{},
			wantErr: false,
		},
		{
			name: "when_grouping_then_aggregate_per_key",
			args: args[int, int, int]{
				slice: []int{1, 2, 3, 4, 5},
				keyFn: func(value int) (int, bool, error) {
					return value % 2, true, nil
				},
				aggFn: func(values []int) (int, bool, error) {
					sum := 0
					for _, value := range values {
						sum += value
					}
					return sum, true, nil
				},
			},
			want: map[int]int{
				0: 6,
				1: 9,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AggregateErr[int, int](tt.args.slice, tt.args.keyFn, tt.args.aggFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("AggregateErr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateErr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
