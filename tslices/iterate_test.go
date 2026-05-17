package tslices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterate(t *testing.T) {
	type args[T any] struct {
		slice []T
		fn    func(T) T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "when_slice_is_empty_then_no_panic",
			args: args[int]{
				slice: []int{},
				fn: func(value int) int {
					value++
					return value
				},
			},
			want: []int{},
		},
		{
			name: "when_slice_has_values_then_apply_fn",
			args: args[int]{
				slice: []int{1, 2},
				fn: func(value int) int {
					value += 1
					return value
				},
			},
			want: func() []int {
				return []int{2, 3}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Iterate(tt.args.slice, tt.args.fn)
			for i, i2 := range tt.want {
				assert.Equal(t, i2, tt.args.slice[i])
			}
		})
	}
}
