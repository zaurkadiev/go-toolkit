package tsync

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_Load(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable] struct {
		key K
	}
	type testCase[K comparable, V any] struct {
		name   string
		fields fields[K, V]
		args   args[K]
		want   V
		want1  bool
	}
	tests := []testCase[string, any]{
		{
			name: "given_nil_map_when_load_then_nil_and_false",
			fields: fields[string, any]{
				m: nil,
			},
			args: args[string]{
				key: "key",
			},
			want:  nil,
			want1: false,
		},
		{
			name: "given_stored_key_when_load_then_value_and_true",
			fields: fields[string, any]{
				m: map[string]any{"key": 1},
			},
			args: args[string]{
				key: "key",
			},
			want:  1,
			want1: true,
		},
		{
			name: "given_other_key_when_load_then_nil_and_false",
			fields: fields[string, any]{
				m: map[string]any{"key": 1},
			},
			args: args[string]{
				key: "key-1",
			},
			want:  nil,
			want1: false,
		},
		{
			name: "given_stored_nil_value_when_load_then_nil_and_true",
			fields: fields[string, any]{
				m: map[string]any{"key": nil},
			},
			args: args[string]{
				key: "key",
			},
			want:  nil,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, any]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			got, got1 := m.Load(tt.args.key)
			assert.Equalf(t, tt.want, got, "Load(%v)", tt.args.key)
			assert.Equalf(t, tt.want1, got1, "Load(%v)", tt.args.key)
		})
	}
}

func TestMap_Store(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable, V any] struct {
		key   K
		value V
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		args       args[K, V]
		wantStored map[K]V
	}
	tests := []testCase[string, int]{
		{
			name: "given_nil_map_when_store_then_value_stored",
			fields: fields[string, int]{
				m: nil,
			},
			args: args[string, int]{
				key:   "key",
				value: 1,
			},
			wantStored: map[string]int{"key": 1},
		},
		{
			name: "given_other_key_when_store_then_both_stored",
			fields: fields[string, int]{
				m: map[string]int{"key": 1},
			},
			args: args[string, int]{
				key:   "key-1",
				value: 2,
			},
			wantStored: map[string]int{"key": 1, "key-1": 2},
		},
		{
			name: "given_stored_key_when_store_then_value_overwritten",
			fields: fields[string, int]{
				m: map[string]int{"key": 1},
			},
			args: args[string, int]{
				key:   "key",
				value: 2,
			},
			wantStored: map[string]int{"key": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, int]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			m.Store(tt.args.key, tt.args.value)
			gotStored := map[string]int{}
			m.Range(func(k string, v int) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "Store(%v, %v)", tt.args.key, tt.args.value)
		})
	}
}

func TestMap_LoadOrStore(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable, V any] struct {
		key   K
		value V
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		args       args[K, V]
		want       V
		want1      bool
		wantStored map[K]V
	}
	tests := []testCase[string, int]{
		{
			name: "given_nil_map_when_load_or_store_then_store_and_false",
			fields: fields[string, int]{
				m: nil,
			},
			args: args[string, int]{
				key:   "key",
				value: 2,
			},
			want:       2,
			want1:      false,
			wantStored: map[string]int{"key": 2},
		},
		{
			name: "given_stored_key_when_load_or_store_then_existing_and_true",
			fields: fields[string, int]{
				m: map[string]int{"key": 1},
			},
			args: args[string, int]{
				key:   "key",
				value: 2,
			},
			want:       1,
			want1:      true,
			wantStored: map[string]int{"key": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, int]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			got, got1 := m.LoadOrStore(tt.args.key, tt.args.value)
			assert.Equalf(t, tt.want, got, "LoadOrStore(%v, %v)", tt.args.key, tt.args.value)
			assert.Equalf(t, tt.want1, got1, "LoadOrStore(%v, %v)", tt.args.key, tt.args.value)
			gotStored := map[string]int{}
			m.Range(func(k string, v int) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "LoadOrStore(%v, %v)", tt.args.key, tt.args.value)
		})
	}
}

func TestMap_LoadAndDelete(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable] struct {
		key K
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		args       args[K]
		want       V
		want1      bool
		wantStored map[K]V
	}
	tests := []testCase[string, int]{
		{
			name: "given_nil_map_when_load_and_delete_then_zero_value_and_false",
			fields: fields[string, int]{
				m: nil,
			},
			args: args[string]{
				key: "key",
			},
			want:       0,
			want1:      false,
			wantStored: map[string]int{},
		},
		{
			name: "given_stored_key_when_load_and_delete_then_value_and_true",
			fields: fields[string, int]{
				m: map[string]int{"key": 1, "key-1": 2},
			},
			args: args[string]{
				key: "key",
			},
			want:       1,
			want1:      true,
			wantStored: map[string]int{"key-1": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, int]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			got, got1 := m.LoadAndDelete(tt.args.key)
			assert.Equalf(t, tt.want, got, "LoadAndDelete(%v)", tt.args.key)
			assert.Equalf(t, tt.want1, got1, "LoadAndDelete(%v)", tt.args.key)
			gotStored := map[string]int{}
			m.Range(func(k string, v int) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "LoadAndDelete(%v)", tt.args.key)
		})
	}
}

func TestMap_Delete(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable] struct {
		key K
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		args       args[K]
		wantStored map[K]V
	}
	tests := []testCase[string, int]{
		{
			name: "given_nil_map_when_delete_then_empty",
			fields: fields[string, int]{
				m: nil,
			},
			args: args[string]{
				key: "key",
			},
			wantStored: map[string]int{},
		},
		{
			name: "given_stored_key_when_delete_then_only_key_removed",
			fields: fields[string, int]{
				m: map[string]int{"key": 1, "key-1": 2},
			},
			args: args[string]{
				key: "key",
			},
			wantStored: map[string]int{"key-1": 2},
		},
		{
			name: "given_other_key_when_delete_then_unchanged",
			fields: fields[string, int]{
				m: map[string]int{"key": 1},
			},
			args: args[string]{
				key: "key-1",
			},
			wantStored: map[string]int{"key": 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, int]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			m.Delete(tt.args.key)
			gotStored := map[string]int{}
			m.Range(func(k string, v int) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "Delete(%v)", tt.args.key)
		})
	}
}

func TestMap_Swap(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable, V any] struct {
		key   K
		value V
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		args       args[K, V]
		want       V
		want1      bool
		wantStored map[K]V
	}
	tests := []testCase[string, int]{
		{
			name: "given_nil_map_when_swap_then_zero_value_and_false",
			fields: fields[string, int]{
				m: nil,
			},
			args: args[string, int]{
				key:   "key",
				value: 2,
			},
			want:       0,
			want1:      false,
			wantStored: map[string]int{"key": 2},
		},
		{
			name: "given_stored_key_when_swap_then_previous_and_true",
			fields: fields[string, int]{
				m: map[string]int{"key": 1},
			},
			args: args[string, int]{
				key:   "key",
				value: 2,
			},
			want:       1,
			want1:      true,
			wantStored: map[string]int{"key": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, int]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			got, got1 := m.Swap(tt.args.key, tt.args.value)
			assert.Equalf(t, tt.want, got, "Swap(%v, %v)", tt.args.key, tt.args.value)
			assert.Equalf(t, tt.want1, got1, "Swap(%v, %v)", tt.args.key, tt.args.value)
			gotStored := map[string]int{}
			m.Range(func(k string, v int) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "Swap(%v, %v)", tt.args.key, tt.args.value)
		})
	}
}

func TestMap_CompareAndSwap(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable, V any] struct {
		key K
		old V
		new V
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		args       args[K, V]
		want       bool
		wantPanic  bool
		wantStored map[K]V
	}
	tests := []testCase[string, any]{
		{
			name: "given_nil_map_when_compare_and_swap_then_false",
			fields: fields[string, any]{
				m: nil,
			},
			args: args[string, any]{
				key: "key",
				old: 1,
				new: 2,
			},
			want:       false,
			wantStored: map[string]any{},
		},
		{
			name: "given_equal_old_value_when_compare_and_swap_then_swapped",
			fields: fields[string, any]{
				m: map[string]any{"key": 1},
			},
			args: args[string, any]{
				key: "key",
				old: 1,
				new: 2,
			},
			want:       true,
			wantStored: map[string]any{"key": 2},
		},
		{
			name: "given_different_old_value_when_compare_and_swap_then_not_swapped",
			fields: fields[string, any]{
				m: map[string]any{"key": 1},
			},
			args: args[string, any]{
				key: "key",
				old: 3,
				new: 2,
			},
			want:       false,
			wantStored: map[string]any{"key": 1},
		},
		{
			name: "given_not_comparable_value_when_compare_and_swap_then_panic",
			fields: fields[string, any]{
				m: map[string]any{"key": []int{1}},
			},
			args: args[string, any]{
				key: "key",
				old: []int{1},
				new: []int{2},
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, any]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			if tt.wantPanic {
				assert.Panicsf(t, func() { m.CompareAndSwap(tt.args.key, tt.args.old, tt.args.new) },
					"CompareAndSwap(%v, %v, %v)", tt.args.key, tt.args.old, tt.args.new)
				return
			}
			got := m.CompareAndSwap(tt.args.key, tt.args.old, tt.args.new)
			assert.Equalf(t, tt.want, got, "CompareAndSwap(%v, %v, %v)", tt.args.key, tt.args.old, tt.args.new)
			gotStored := map[string]any{}
			m.Range(func(k string, v any) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "CompareAndSwap(%v, %v, %v)", tt.args.key, tt.args.old, tt.args.new)
		})
	}
}

func TestMap_CompareAndDelete(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable, V any] struct {
		key K
		old V
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		args       args[K, V]
		want       bool
		wantPanic  bool
		wantStored map[K]V
	}
	tests := []testCase[string, any]{
		{
			name: "given_nil_map_when_compare_and_delete_then_false",
			fields: fields[string, any]{
				m: nil,
			},
			args: args[string, any]{
				key: "key",
				old: 1,
			},
			want:       false,
			wantStored: map[string]any{},
		},
		{
			name: "given_equal_old_value_when_compare_and_delete_then_deleted",
			fields: fields[string, any]{
				m: map[string]any{"key": 1},
			},
			args: args[string, any]{
				key: "key",
				old: 1,
			},
			want:       true,
			wantStored: map[string]any{},
		},
		{
			name: "given_different_old_value_when_compare_and_delete_then_kept",
			fields: fields[string, any]{
				m: map[string]any{"key": 1},
			},
			args: args[string, any]{
				key: "key",
				old: 3,
			},
			want:       false,
			wantStored: map[string]any{"key": 1},
		},
		{
			name: "given_not_comparable_value_when_compare_and_delete_then_panic",
			fields: fields[string, any]{
				m: map[string]any{"key": []int{1}},
			},
			args: args[string, any]{
				key: "key",
				old: []int{1},
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, any]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			if tt.wantPanic {
				assert.Panicsf(t, func() { m.CompareAndDelete(tt.args.key, tt.args.old) },
					"CompareAndDelete(%v, %v)", tt.args.key, tt.args.old)
				return
			}
			got := m.CompareAndDelete(tt.args.key, tt.args.old)
			assert.Equalf(t, tt.want, got, "CompareAndDelete(%v, %v)", tt.args.key, tt.args.old)
			gotStored := map[string]any{}
			m.Range(func(k string, v any) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "CompareAndDelete(%v, %v)", tt.args.key, tt.args.old)
		})
	}
}

func TestMap_Range(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type args[K comparable, V any] struct {
		f func(m *Map[K, V]) func(key K, value V) bool
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		args       args[K, V]
		want       int
		wantStored map[K]V
	}
	tests := []testCase[string, int]{
		{
			name: "given_nil_map_when_range_then_no_calls",
			fields: fields[string, int]{
				m: nil,
			},
			args: args[string, int]{
				f: func(_ *Map[string, int]) func(string, int) bool {
					return func(string, int) bool { return true }
				},
			},
			want:       0,
			wantStored: map[string]int{},
		},
		{
			name: "given_entries_when_range_then_visit_all",
			fields: fields[string, int]{
				m: map[string]int{"key": 1, "key-1": 2, "key-2": 3},
			},
			args: args[string, int]{
				f: func(_ *Map[string, int]) func(string, int) bool {
					return func(string, int) bool { return true }
				},
			},
			want:       3,
			wantStored: map[string]int{"key": 1, "key-1": 2, "key-2": 3},
		},
		{
			name: "given_entries_when_f_returns_false_then_stop",
			fields: fields[string, int]{
				m: map[string]int{"key": 1, "key-1": 2, "key-2": 3},
			},
			args: args[string, int]{
				f: func(_ *Map[string, int]) func(string, int) bool {
					return func(string, int) bool { return false }
				},
			},
			want:       1,
			wantStored: map[string]int{"key": 1, "key-1": 2, "key-2": 3},
		},
		{
			name: "given_entries_when_f_modifies_map_then_no_deadlock",
			fields: fields[string, int]{
				m: map[string]int{"key": 1, "key-1": 2},
			},
			args: args[string, int]{
				f: func(m *Map[string, int]) func(string, int) bool {
					return func(k string, _ int) bool {
						m.Delete(k)
						return true
					}
				},
			},
			want:       2,
			wantStored: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, int]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			f := tt.args.f(m)
			got := 0
			m.Range(func(k string, v int) bool {
				got++
				return f(k, v)
			})
			assert.Equalf(t, tt.want, got, "Range(%T)", tt.args.f)
			gotStored := map[string]int{}
			m.Range(func(k string, v int) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "Range(%T)", tt.args.f)
		})
	}
}

func TestMap_Clear(t *testing.T) {
	type fields[K comparable, V any] struct {
		m map[K]V
	}
	type testCase[K comparable, V any] struct {
		name       string
		fields     fields[K, V]
		wantStored map[K]V
	}
	tests := []testCase[string, int]{
		{
			name: "given_nil_map_when_clear_then_empty",
			fields: fields[string, int]{
				m: nil,
			},
			wantStored: map[string]int{},
		},
		{
			name: "given_entries_when_clear_then_empty",
			fields: fields[string, int]{
				m: map[string]int{"key": 1, "key-1": 2},
			},
			wantStored: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[string, int]{}
			for k, v := range tt.fields.m {
				m.Store(k, v)
			}
			m.Clear()
			gotStored := map[string]int{}
			m.Range(func(k string, v int) bool {
				gotStored[k] = v
				return true
			})
			assert.Equalf(t, tt.wantStored, gotStored, "Clear()")
		})
	}
}

func TestCast(t *testing.T) {
	err := errors.New("error")
	type args struct {
		v any
	}
	type testCase[T any] struct {
		name string
		args args
		want T
	}
	tests := []testCase[error]{
		{
			name: "given_nil_when_cast_then_zero_value",
			args: args{
				v: nil,
			},
			want: nil,
		},
		{
			name: "given_value_of_type_when_cast_then_value",
			args: args{
				v: err,
			},
			want: err,
		},
		{
			name: "given_value_of_other_type_when_cast_then_zero_value",
			args: args{
				v: "value",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, cast[error](tt.args.v), "cast(%v)", tt.args.v)
		})
	}
}
