package tmaps

import (
	"maps"
	"testing"
)

func TestNilToEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		m       map[string]int
		want    map[string]int
		wantNil bool
	}{
		{
			// The input is a nil map.
			// The function returns an empty, non-nil map.
			name:    "given_nil_map_when_normalized_then_returns_empty_map",
			m:       nil,
			want:    map[string]int{},
			wantNil: false,
		},
		{
			// The input is an already allocated empty map.
			// The function returns the map unchanged and non-nil.
			name:    "given_empty_map_when_normalized_then_stays_empty",
			m:       map[string]int{},
			want:    map[string]int{},
			wantNil: false,
		},
		{
			// The input is a map with values.
			// The function returns the same values and does not reset it to an empty map.
			name:    "given_map_with_values_when_normalized_then_unchanged",
			m:       map[string]int{"first": 1, "second": 2},
			want:    map[string]int{"first": 1, "second": 2},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NilToEmpty(tt.m)

			if (got == nil) != tt.wantNil {
				t.Fatalf("NilToEmpty() nil = %v, want %v", got == nil, tt.wantNil)
			}
			if !maps.Equal(got, tt.want) {
				t.Fatalf("NilToEmpty() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNilToEmptyNamedMap(t *testing.T) {
	t.Parallel()

	type namedMap map[string]int

	tests := []struct {
		name string
		m    namedMap
		want namedMap
	}{
		{
			// The input is a nil map of a named type.
			// The function returns an empty map of the same named type.
			name: "given_named_nil_map_when_normalized_then_keeps_type",
			m:    nil,
			want: namedMap{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got namedMap = NilToEmpty(tt.m)

			if got == nil {
				t.Fatal("NilToEmpty() returned nil")
			}
			if !maps.Equal(got, tt.want) {
				t.Fatalf("NilToEmpty() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
