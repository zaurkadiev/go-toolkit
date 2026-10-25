package optional

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	type testCase[T any] struct {
		name string
		want O[int]
	}
	tests := []testCase[int]{
		{
			name: "given_nothing_when_empty_then_returns_empty_optional",
			want: O[int]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Empty[int]())
		})
	}
}

func TestOptional_GetValueOrDefault(t *testing.T) {
	type args[T any] struct {
		defaultValue T
	}
	type testCase[T any] struct {
		name string
		o    O[T]
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "given_empty_optional_when_get_or_default_then_returns_default",
			o:    O[int]{},
			args: args[int]{
				defaultValue: 5,
			},
			want: 5,
		},
		{
			name: "given_non_empty_optional_when_get_or_default_then_returns_value",
			o: O[int]{
				Value: func() *int {
					i := 10
					return &i
				}(),
			},
			args: args[int]{
				defaultValue: 5,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.GetOrDefault(tt.args.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptional_IsEmpty(t *testing.T) {
	type testCase[T any] struct {
		name string
		o    O[T]
		want bool
	}
	tests := []testCase[int]{
		{
			name: "given_empty_optional_when_is_empty_then_true",
			o:    O[int]{},
			want: true,
		},
		{
			name: "given_non_empty_optional_when_is_empty_then_false",
			o: O[int]{
				Value: func() *int {
					i := 10
					return &i
				}(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptional_Unwrap(t *testing.T) {
	type testCase[T any] struct {
		name    string
		o       O[T]
		want    T
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name:    "given_empty_optional_when_unwrap_then_error",
			o:       O[int]{},
			want:    0,
			wantErr: true,
		},
		{
			name: "given_non_empty_optional_when_unwrap_then_value",
			o: O[int]{
				Value: func() *int {
					i := 10
					return &i
				}(),
			},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.o.Unwrap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unwrap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args[T any] struct {
		value *T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want O[int]
	}
	valPtr := func() *int {
		i := 1
		return &i
	}()
	tests := []testCase[int]{
		{
			name: "given_value_when_new_then_returns_optional_with_value",
			args: args[int]{
				value: valPtr,
			},
			want: O[int]{
				Value: valPtr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, New(tt.args.value))
		})
	}
}

func TestOptional_String(t *testing.T) {
	type testCase[T any] struct {
		name string
		o    O[T]
		want string
	}
	valPtr := func() *string {
		s := "cluster-1"
		return &s
	}()
	tests := []testCase[string]{
		{
			// The optional holds no value.
			// The method must return an explicit empty marker so logs do not look like a pointer dump.
			name: "given_empty_optional_when_string_then_empty_marker",
			o:    O[string]{},
			want: "<empty>",
		},
		{
			// The optional holds a string.
			// The method must return the value itself without wrapping the Value field.
			name: "given_non_empty_optional_when_string_then_value",
			o:    O[string]{Value: valPtr},
			want: "cluster-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.o.String())
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    O[int]
		f        func(int) string
		expected O[string]
	}{
		{
			// The input value is present. The function is applied to it.
			// The result holds the transformed value.
			name:     "given_non_empty_optional_when_map_then_mapped_value",
			input:    New(ptr(42)),
			f:        func(n int) string { return "42" },
			expected: New(ptr("42")),
		},
		{
			// The input value is absent. The function is not called.
			// The result is an empty optional.
			name:     "given_empty_optional_when_map_then_empty_result",
			input:    Empty[int](),
			f:        func(n int) string { return "should not be called" },
			expected: Empty[string](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.f)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestNewZeroAsEmpty(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want O[int]
	}{
		{
			name: "zero_value_returns_empty",
			in:   0,
			want: O[int]{},
		},
		{
			name: "non_zero_returns_value",
			in:   7,
			want: New(ptr(7)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewZeroAsEmpty(tt.in))
		})
	}
}

func TestOptional_Get(t *testing.T) {
	tests := []struct {
		name string
		o    O[int]
		want int
	}{
		{
			name: "empty_returns_zero",
			o:    Empty[int](),
			want: 0,
		},
		{
			name: "non_empty_returns_value",
			o:    New(ptr(42)),
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.o.Get())
		})
	}
}
