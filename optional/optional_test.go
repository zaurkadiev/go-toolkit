package optional

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    O[int]
		f        func(int) string
		expected O[string]
	}{
		{
			// Входное значение присутствует. Функция применяется к значению.
			// Результат содержит преобразованное значение.
			name:     "успех_при_непустом_optional",
			input:    New(ptr(42)),
			f:        func(n int) string { return "42" },
			expected: New(ptr("42")),
		},
		{
			// Входное значение отсутствует. Функция не вызывается.
			// Результат — пустой optional.
			name:     "пустой_результат_при_пустом_входе",
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
