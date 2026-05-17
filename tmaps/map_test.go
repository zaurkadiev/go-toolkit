package tmaps

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	type args[K comparable, V any, K1 comparable, V1 any] struct {
		m  map[K]V
		fn func(K, V) (K1, V1, bool)
	}
	type testCase[K comparable, V any, K1 comparable, V1 any] struct {
		name string
		args args[K, V, K1, V1]
		want map[K1]V1
	}
	tests := []testCase[int, int, int, int]{
		{
			name: "given_nil_map_when_map_then_nil",
			args: args[int, int, int, int]{
				m:  nil,
				fn: func(k, v int) (int, int, bool) { return k, v, true },
			},
			want: nil,
		},
		{
			name: "given_empty_map_when_map_then_empty_map",
			args: args[int, int, int, int]{
				m:  map[int]int{},
				fn: func(k, v int) (int, int, bool) { return k, v, true },
			},
			want: map[int]int{},
		},
		{
			name: "given_map_when_map_then_transform_key_and_value",
			args: args[int, int, int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
				},
				fn: func(k, v int) (int, int, bool) {
					return k + 10, k, true
				},
			},
			want: map[int]int{
				11: 1,
				12: 2,
			},
		},
		{
			name: "given_map_when_map_then_skip_by_ok",
			args: args[int, int, int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(k, v int) (int, int, bool) {
					if k == 2 {
						return 0, 0, false
					}
					return k, v * 2, true
				},
			},
			want: map[int]int{
				1: 20,
				3: 60,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Map(tt.args.m, tt.args.fn), "Map(%v, %T)", tt.args.m, tt.args.fn)
		})
	}
}

func TestMapErr(t *testing.T) {
	type args[K comparable, V any, K1 comparable, V1 any] struct {
		m  map[K]V
		fn func(K, V) (K1, V1, bool, error)
	}
	type testCase[K comparable, V any, K1 comparable, V1 any] struct {
		name    string
		args    args[K, V, K1, V1]
		want    map[K1]V1
		wantErr bool
	}
	tests := []testCase[int, int, int, int]{
		{
			name: "given_nil_map_when_map_err_then_nil",
			args: args[int, int, int, int]{
				m:  nil,
				fn: func(k, v int) (int, int, bool, error) { return k, v, true, nil },
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "given_empty_map_when_map_err_then_empty_map",
			args: args[int, int, int, int]{
				m:  map[int]int{},
				fn: func(k, v int) (int, int, bool, error) { return k, v, true, nil },
			},
			want:    map[int]int{},
			wantErr: false,
		},
		{
			name: "given_map_when_map_err_then_transform_key_and_value",
			args: args[int, int, int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
				},
				fn: func(k, v int) (int, int, bool, error) {
					return k + 10, k, true, nil
				},
			},
			want: map[int]int{
				11: 1,
				12: 2,
			},
			wantErr: false,
		},
		{
			name: "given_map_when_map_err_then_skip_by_ok",
			args: args[int, int, int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(k, v int) (int, int, bool, error) {
					if k == 2 {
						return 0, 0, false, nil
					}
					return k, v * 2, true, nil
				},
			},
			want: map[int]int{
				1: 20,
				3: 60,
			},
			wantErr: false,
		},
		{
			name: "given_map_when_map_err_then_return_error",
			args: args[int, int, int, int]{
				m: map[int]int{
					1: 10,
				},
				fn: func(_, _ int) (int, int, bool, error) {
					return 0, 0, false, assert.AnError
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapErr(tt.args.m, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapErr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "MapErr(%v, %T)", tt.args.m, tt.args.fn)
		})
	}
}

func TestToSlice(t *testing.T) {
	type args[K comparable, V any, V1 any] struct {
		m  map[K]V
		fn func(K, V) (V1, bool)
	}
	type testCase[K comparable, V any, V1 any] struct {
		name string
		args args[K, V, V1]
		want []V1
	}
	tests := []testCase[int, int, int]{
		{
			name: "given_nil_map_when_mapped_then_nil",
			args: args[int, int, int]{
				m:  nil,
				fn: func(k, v int) (int, bool) { return k + v, true },
			},
			want: nil,
		},
		{
			name: "given_empty_map_when_mapped_then_empty_slice",
			args: args[int, int, int]{
				m:  map[int]int{},
				fn: func(k, v int) (int, bool) { return k + v, true },
			},
			want: []int{},
		},
		{
			name: "given_map_when_mapped_then_return_slice",
			args: args[int, int, int]{
				m:  map[int]int{1: 10, 2: 20},
				fn: func(k, v int) (int, bool) { return k + v, true },
			},
			want: []int{11, 22},
		},
		{
			name: "given_map_when_mapped_with_when_skip_by_ok",
			args: args[int, int, int]{
				m: map[int]int{1: 10, 2: 20},
				fn: func(k, v int) (int, bool) {
					if k%2 == 0 {
						return 0, false
					}
					return k + v, true
				},
			},
			want: []int{11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatchf(t, tt.want, ToSlice(tt.args.m, tt.args.fn), "ToSlice(%v, %v)", tt.args.m, tt.args.fn)
		})
	}
}

func TestToSliceErr(t *testing.T) {
	type args[K comparable, V any, V1 any] struct {
		m  map[K]V
		fn func(K, V) (V1, bool, error)
	}
	type testCase[K comparable, V any, V1 any] struct {
		name    string
		args    args[K, V, V1]
		want    []V1
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[int, int, int]{
		{
			name: "given_nil_map_when_mapped_then_nil",
			args: args[int, int, int]{
				m:  nil,
				fn: func(k, v int) (int, bool, error) { return k + v, true, nil },
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "given_empty_map_when_mapped_then_empty_slice",
			args: args[int, int, int]{
				m:  map[int]int{},
				fn: func(k, v int) (int, bool, error) { return k + v, true, nil },
			},
			want:    []int{},
			wantErr: assert.NoError,
		},
		{
			name: "given_map_when_mapped_then_return_slice",
			args: args[int, int, int]{
				m:  map[int]int{1: 10, 2: 20},
				fn: func(k, v int) (int, bool, error) { return k + v, true, nil },
			},
			want:    []int{11, 22},
			wantErr: assert.NoError,
		},
		{
			name: "given_map_when_mapped_with_when_skip_by_ok",
			args: args[int, int, int]{
				m: map[int]int{1: 10, 2: 20},
				fn: func(k, v int) (int, bool, error) {
					if k%2 == 0 {
						return 0, false, nil
					}
					return k + v, true, nil
				},
			},
			want:    []int{11},
			wantErr: assert.NoError,
		},
		{
			name: "given_map_when_mapped_with_error_then_return_error",
			args: args[int, int, int]{
				m: map[int]int{1: 10, 2: 20},
				fn: func(k, v int) (int, bool, error) {
					return 0, false, errors.New("error")
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToSliceErr(tt.args.m, tt.args.fn)
			if !tt.wantErr(t, err, fmt.Sprintf("ToSliceErr(%v, %T)", tt.args.m, tt.args.fn)) {
				return
			}
			assert.ElementsMatchf(t, tt.want, got, "ToSliceErr(%v, %T)", tt.args.m, tt.args.fn)
		})
	}
}
