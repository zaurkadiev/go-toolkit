package tmaps

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapVals(t *testing.T) {
	type args[K comparable, V any, R any] struct {
		m  map[K]V
		fn func(V) (R, bool)
	}
	type testCase[K comparable, V any, R any] struct {
		name string
		args args[K, V, R]
		want map[K]R
	}
	tests := []testCase[int, int, int]{
		{
			name: "given_nil_map_when_map_vals_then_nil",
			args: args[int, int, int]{
				m:  nil,
				fn: func(v int) (int, bool) { return v, true },
			},
			want: nil,
		},
		{
			name: "given_empty_map_when_map_vals_then_empty_map",
			args: args[int, int, int]{
				m:  map[int]int{},
				fn: func(v int) (int, bool) { return v, true },
			},
			want: map[int]int{},
		},
		{
			name: "given_map_when_map_vals_then_transform_values",
			args: args[int, int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
				},
				fn: func(v int) (int, bool) { return 2 * v, true },
			},
			want: map[int]int{
				1: 20,
				2: 40,
			},
		},
		{
			name: "given_map_when_map_vals_then_skip_by_ok",
			args: args[int, int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(v int) (int, bool) {
					if v == 20 {
						return 0, false
					}
					return v * 2, true
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
			got := MapVals(tt.args.m, tt.args.fn)
			assert.Equalf(t, tt.want, got, "MapVals(%v, %T)", tt.args.m, tt.args.fn)
		})
	}
}

func TestMapValsErr(t *testing.T) {
	type args[K comparable, V any, R any] struct {
		m  map[K]V
		fn func(V) (R, bool, error)
	}
	type testCase[K comparable, V any, R any] struct {
		name    string
		args    args[K, V, R]
		want    map[K]R
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[int, int, int]{
		{
			name: "given_nil_map_when_map_vals_err_then_nil",
			args: args[int, int, int]{
				m:  nil,
				fn: func(v int) (int, bool, error) { return v, true, nil },
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "given_empty_map_when_map_vals_err_then_empty_map",
			args: args[int, int, int]{
				m:  map[int]int{},
				fn: func(v int) (int, bool, error) { return v, true, nil },
			},
			want:    map[int]int{},
			wantErr: assert.NoError,
		},
		{
			name: "given_map_when_map_vals_err_then_transform_values",
			args: args[int, int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
				},
				fn: func(v int) (int, bool, error) { return 2 * v, true, nil },
			},
			want: map[int]int{
				1: 20,
				2: 40,
			},
			wantErr: assert.NoError,
		},
		{
			name: "given_map_when_map_vals_err_then_skip_by_ok",
			args: args[int, int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
					3: 30,
				},
				fn: func(v int) (int, bool, error) {
					if v == 20 {
						return 0, false, nil
					}
					return v * 2, true, nil
				},
			},
			want: map[int]int{
				1: 20,
				3: 60,
			},
			wantErr: assert.NoError,
		},
		{
			name: "given_map_when_map_vals_err_then_error",
			args: args[int, int, int]{
				m: map[int]int{
					1: 10,
				},
				fn: func(_ int) (int, bool, error) {
					return 0, false, assert.AnError
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapValsErr(tt.args.m, tt.args.fn)
			if !tt.wantErr(t, err, fmt.Sprintf("MapValsErr(%v, %T)", tt.args.m, tt.args.fn)) {
				return
			}
			assert.Equalf(t, tt.want, got, "MapValsErr(%v, %T)", tt.args.m, tt.args.fn)
		})
	}
}
