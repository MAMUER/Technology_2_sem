package stringsx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClip(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		result := Clip("", 5)
		assert.Equal(t, "", result)

		result = Clip("", 0)
		assert.Equal(t, "", result)

		result = Clip("", -1)
		assert.Equal(t, "", result)
	})

	t.Run("MaxZero", func(t *testing.T) {
		// max = 0
		result := Clip("hello", 0)
		assert.Equal(t, "", result)

		result = Clip("test string", 0)
		assert.Equal(t, "", result)

		result = Clip("a", 0)
		assert.Equal(t, "", result)
	})

	t.Run("MaxNegative", func(t *testing.T) {
		// max < 0
		result := Clip("hello", -1)
		assert.Equal(t, "", result)

		result = Clip("test", -5)
		assert.Equal(t, "", result)

		result = Clip("abc", -10)
		assert.Equal(t, "", result)
	})

	t.Run("MaxEqualsLength", func(t *testing.T) {
		// max == len(s)
		result := Clip("hello", 5)
		assert.Equal(t, "hello", result)

		result = Clip("test", 4)
		assert.Equal(t, "test", result)

		result = Clip("a", 1)
		assert.Equal(t, "a", result)
	})

	t.Run("MaxGreaterThanLength", func(t *testing.T) {
		// max > len(s)
		result := Clip("hello", 10)
		assert.Equal(t, "hello", result)

		result = Clip("test", 100)
		assert.Equal(t, "test", result)

		result = Clip("", 5)
		assert.Equal(t, "", result)
	})

	t.Run("MaxLessThanLength", func(t *testing.T) {
		// max < len(s)
		result := Clip("hello world", 5)
		assert.Equal(t, "hello", result)

		result = Clip("testing", 4)
		assert.Equal(t, "test", result)

		result = Clip("abcdef", 3)
		assert.Equal(t, "abc", result)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		// Граничные случаи
		result := Clip("a", 1)
		assert.Equal(t, "a", result)

		result = Clip("a", 0)
		assert.Equal(t, "", result)

		result = Clip("ab", 1)
		assert.Equal(t, "a", result)

		result = Clip("abc", 2)
		assert.Equal(t, "ab", result)
	})
}

// Дополнительные тесты
func TestClip_TableDriven(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		max      int
		expected string
	}{
		{"EmptyString_Max5", "", 5, ""},
		{"EmptyString_Max0", "", 0, ""},
		{"EmptyString_MaxNegative", "", -1, ""},
		{"Hello_Max0", "hello", 0, ""},
		{"Hello_Max3", "hello", 3, "hel"},
		{"Hello_Max5", "hello", 5, "hello"},
		{"Hello_Max10", "hello", 10, "hello"},
		{"Test_MaxNegative", "test", -5, ""},
		{"SingleChar_Max1", "a", 1, "a"},
		{"SingleChar_Max0", "a", 0, ""},
		{"LongString_Max10", "this is a long string", 10, "this is a "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Clip(tc.input, tc.max)
			assert.Equal(t, tc.expected, result,
				"Clip(%q, %d) = %q, expected %q",
				tc.input, tc.max, result, tc.expected)
		})
	}
}
