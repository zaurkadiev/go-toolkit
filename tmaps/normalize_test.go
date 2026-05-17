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
			// На вход приходит nil-мапа.
			// Функция возвращает пустую, но не nil-мапу.
			name:    "Успех,nil_мапа_становится_пустой",
			m:       nil,
			want:    map[string]int{},
			wantNil: false,
		},
		{
			// На вход приходит уже созданная пустая мапа.
			// Функция возвращает пустую, но не nil-мапу без изменения содержимого.
			name:    "Успех,пустая_мапа_остается_пустой",
			m:       map[string]int{},
			want:    map[string]int{},
			wantNil: false,
		},
		{
			// На вход приходит мапа со значениями.
			// Функция возвращает те же значения и не нормализует ее в пустую мапу.
			name:    "Успех,мапа_со_значениями_не_меняется",
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
			// На вход приходит nil-мапа именованного типа.
			// Функция возвращает пустую мапу того же именованного типа.
			name: "Успех,именованная_nil_мапа_сохраняет_тип",
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
