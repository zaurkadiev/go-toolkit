package tslices

import (
	"slices"
	"testing"
)

func TestNilToEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		slice   []int
		want    []int
		wantNil bool
	}{
		{
			// The input is a nil slice.
			// The function returns an empty, non-nil slice.
			name:    "given_nil_slice_when_normalized_then_returns_empty_slice",
			slice:   nil,
			want:    []int{},
			wantNil: false,
		},
		{
			// The input is an already allocated empty slice.
			// The function returns the slice unchanged and non-nil.
			name:    "given_empty_slice_when_normalized_then_stays_empty",
			slice:   []int{},
			want:    []int{},
			wantNil: false,
		},
		{
			// The input is a slice with values.
			// The function returns the same values and does not reset it to an empty slice.
			name:    "given_slice_with_values_when_normalized_then_unchanged",
			slice:   []int{1, 2, 3},
			want:    []int{1, 2, 3},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Nil2Empty(tt.slice)

			if (got == nil) != tt.wantNil {
				t.Fatalf("NilToEmpty() nil = %v, want %v", got == nil, tt.wantNil)
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("NilToEmpty() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNilToEmptyNamedSlice(t *testing.T) {
	t.Parallel()

	type namedSlice []int

	tests := []struct {
		name  string
		slice namedSlice
		want  namedSlice
	}{
		{
			// The input is a nil slice of a named type.
			// The function returns an empty slice of the same named type.
			name:  "given_named_nil_slice_when_normalized_then_keeps_type",
			slice: nil,
			want:  namedSlice{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got namedSlice = Nil2Empty(tt.slice)

			if got == nil {
				t.Fatal("NilToEmpty() returned nil")
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("NilToEmpty() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
