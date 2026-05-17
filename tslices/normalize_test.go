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
			// На вход приходит nil-слайс.
			// Функция возвращает пустой, но не nil-слайс.
			name:    "Успех,nil_слайс_становится_пустым",
			slice:   nil,
			want:    []int{},
			wantNil: false,
		},
		{
			// На вход приходит уже созданный пустой слайс.
			// Функция возвращает пустой, но не nil-слайс без изменения содержимого.
			name:    "Успех,пустой_слайс_остается_пустым",
			slice:   []int{},
			want:    []int{},
			wantNil: false,
		},
		{
			// На вход приходит слайс со значениями.
			// Функция возвращает те же значения и не нормализует его в пустой слайс.
			name:    "Успех,слайс_со_значениями_не_меняется",
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
			// На вход приходит nil-слайс именованного типа.
			// Функция возвращает пустой слайс того же именованного типа.
			name:  "Успех,именованный_nil_слайс_сохраняет_тип",
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
