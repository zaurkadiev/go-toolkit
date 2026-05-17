package tmaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	type args[K comparable, V any] struct {
		maps []map[K]V
	}
	type testCase[K comparable, V any] struct {
		name string
		args args[K, V]
		want map[K]V
	}
	tests := []testCase[int, int]{
		{
			name: "given_nil_maps_when_merge_then_nil",
			args: args[int, int]{
				maps: nil,
			},
			want: nil,
		},
		{
			name: "given_empty_maps_when_merge_then_empty_map",
			args: args[int, int]{
				maps: make([]map[int]int, 0),
			},
			want: map[int]int{},
		},
		{
			name: "given_all_empty_maps_when_merge_then_empty_map",
			args: args[int, int]{
				maps: []map[int]int{
					{},
					{},
				},
			},
			want: map[int]int{},
		},
		{
			name: "given_single_map_when_merge_then_same",
			args: args[int, int]{
				maps: []map[int]int{
					{
						1: 10,
						2: 20,
					},
				},
			},
			want: map[int]int{
				1: 10,
				2: 20,
			},
		},
		{
			name: "given_multiple_maps_when_merge_then_later_stays",
			args: args[int, int]{
				maps: []map[int]int{
					{
						1: 10,
						2: 20,
					},
					{
						2: 200,
						3: 30,
					},
				},
			},
			want: map[int]int{
				1: 10,
				2: 200,
				3: 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.args.maps...)
			assert.Equalf(t, tt.want, got, "Merge() = %v, want %v", got, tt.want)
		})
	}
}
