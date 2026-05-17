package tmaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsure(t *testing.T) {
	type args[K comparable, V any] struct {
		m            map[K]V
		key          K
		defaultValue V
	}
	type testCase[K comparable, V any] struct {
		name    string
		args    args[K, V]
		want    V
		wantMap map[K]V
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_map_when_ensure_then_return_default",
			args: args[int, int]{
				m:            nil,
				key:          1,
				defaultValue: 10,
			},
			want:    10,
			wantMap: nil,
		},
		{
			name: "given_empty_map_when_ensure_then_insert_default",
			args: args[int, int]{
				m:            map[int]int{},
				key:          1,
				defaultValue: 10,
			},
			want: 10,
			wantMap: map[int]int{
				1: 10,
			},
		},
		{
			name: "given_existing_key_when_ensure_then_return_existing",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
				},
				key:          1,
				defaultValue: 20,
			},
			want: 10,
			wantMap: map[int]int{
				1: 10,
			},
		},
		{
			name: "given_missing_key_when_ensure_then_insert_default",
			args: args[int, int]{
				m: map[int]int{
					1: 10,
				},
				key:          2,
				defaultValue: 20,
			},
			want: 20,
			wantMap: map[int]int{
				1: 10,
				2: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ensure(tt.args.m, tt.args.key, tt.args.defaultValue)
			assert.Equalf(t, tt.wantMap, tt.args.m, "Ensure() gotMap = %v, wantMap %v", tt.args.m, tt.wantMap)
			assert.Equalf(t, tt.want, got, "Ensure() got = %v, want %v", got, tt.want)
		})
	}
}
