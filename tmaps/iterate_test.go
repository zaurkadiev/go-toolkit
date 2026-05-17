package tmaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterate(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		wantMap map[K]V
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_map_when_iterate_then_no_panic",
			args: args[int, int]{
				m: nil,
			},
			wantMap: nil,
		},
		{
			name: "given_empty_map_when_iterate_then_no_change",
			args: args[int, int]{
				m: map[int]int{},
			},
			wantMap: map[int]int{},
		},
		{
			name: "given_map_when_iterate_then_call_fn_for_each_pair",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
					2: 20,
				},
			},
			wantMap: map[int]int{
				1: 11,
				2: 22,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Iterate(tt.args.m, func(k, v int) int {
				return k + v
			})
			assert.Equalf(t, tt.args.m, tt.wantMap, "Iterate() map = %v, want %v", tt.args.m, tt.wantMap)
		})
	}
}
