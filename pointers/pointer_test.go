package pointers

import (
	"reflect"
	"testing"
)

func TestOf(t *testing.T) {
	type args[T any] struct {
		value T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *T
	}
	tests := []testCase[int]{
		{
			name: "given_zero_value_when_of_then_pointer_with_zero",
			args: args[int]{value: 0},
			want: func() *int { v := 0; return &v }(),
		},
		{
			name: "given_non_zero_value_when_of_then_pointer_with_value",
			args: args[int]{value: 10},
			want: func() *int { v := 10; return &v }(),
		},
		{
			name: "given_negative_value_when_of_then_pointer_with_value",
			args: args[int]{value: -5},
			want: func() *int { v := -5; return &v }(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Of(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Of() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueOf(t *testing.T) {
	type args[T any] struct {
		ptr *T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "given_nil_pointer_when_value_of_then_zero_value",
			args: args[int]{ptr: nil},
			want: 0,
		},
		{
			name: "given_pointer_when_value_of_then_value",
			args: args[int]{ptr: func() *int { v := 10; return &v }()},
			want: 10,
		},
		{
			name: "given_pointer_to_zero_when_value_of_then_zero_value",
			args: args[int]{ptr: func() *int { v := 0; return &v }()},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValueOf(tt.args.ptr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValueOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
