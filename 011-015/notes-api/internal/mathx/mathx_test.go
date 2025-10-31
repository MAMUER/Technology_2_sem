package mathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSum_Table - табличные тесты с чистым testing
func TestSum_Table(t *testing.T) {
	cases := []struct {
		a, b, want int
	}{
		{2, 3, 5},
		{10, -5, 5},
		{0, 0, 0},
		{-1, -1, -2},
		{100, 200, 300},
		{-10, 15, 5},
	}

	for _, c := range cases {
		got := Sum(c.a, c.b)
		if got != c.want {
			t.Errorf("Sum(%d, %d) = %d; want %d", c.a, c.b, got, c.want)
		}
	}
}

// TestDivide_Table - табличные тесты для деления
func TestDivide_Table(t *testing.T) {
	cases := []struct {
		a, b    int
		want    int
		wantErr bool
	}{
		{10, 2, 5, false},
		{15, 3, 5, false},
		{8, 4, 2, false},
		{10, 0, 0, true},    // ошибка деления на ноль
		{0, 5, 0, false},    // ноль делится на любое число
		{-10, 2, -5, false}, // отрицательные числа
	}

	for _, c := range cases {
		got, err := Divide(c.a, c.b)

		if c.wantErr {
			if err == nil {
				t.Errorf("Divide(%d, %d) expected error, but got none", c.a, c.b)
			}
		} else {
			if err != nil {
				t.Errorf("Divide(%d, %d) unexpected error: %v", c.a, c.b, err)
			}
			if got != c.want {
				t.Errorf("Divide(%d, %d) = %d; want %d", c.a, c.b, got, c.want)
			}
		}
	}
}

// TestDivide_OkAndError - отдельные проверки с testify
func TestDivide_OkAndError(t *testing.T) {
	t.Run("NormalDivision", func(t *testing.T) {
		// Нормальное деление
		got, err := Divide(10, 2)
		require.NoError(t, err)
		assert.Equal(t, 5, got)

		got, err = Divide(15, 3)
		require.NoError(t, err)
		assert.Equal(t, 5, got)
	})

	t.Run("DivisionByZero", func(t *testing.T) {
		// Деление на ноль
		_, err := Divide(10, 0)
		assert.Error(t, err)
		assert.Equal(t, "divide by zero", err.Error())

		_, err = Divide(0, 0)
		assert.Error(t, err)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		// Граничные случаи
		got, err := Divide(0, 5)
		require.NoError(t, err)
		assert.Equal(t, 0, got)

		got, err = Divide(-10, 2)
		require.NoError(t, err)
		assert.Equal(t, -5, got)
	})
}

// TestSum_WithTestify - тесты суммы с использованием testify
func TestSum_WithTestify(t *testing.T) {
	assert.Equal(t, 5, Sum(2, 3))
	assert.Equal(t, 5, Sum(10, -5))
	assert.Equal(t, 0, Sum(0, 0))
	assert.Equal(t, -2, Sum(-1, -1))
	assert.Equal(t, 300, Sum(100, 200))
}

// TestDivide_WithTestify - комплексные тесты деления с testify
func TestDivide_WithTestify(t *testing.T) {
	t.Run("SuccessfulDivision", func(t *testing.T) {
		result, err := Divide(20, 4)
		assert.NoError(t, err)
		assert.Equal(t, 5, result)
	})

	t.Run("DivisionByZero", func(t *testing.T) {
		result, err := Divide(10, 0)
		assert.Error(t, err)
		assert.Equal(t, 0, result)
		assert.Contains(t, err.Error(), "divide by zero")
	})

	t.Run("NegativeNumbers", func(t *testing.T) {
		result, err := Divide(-15, 3)
		assert.NoError(t, err)
		assert.Equal(t, -5, result)
	})
}

// Тесты для Factorial
func TestFactorial(t *testing.T) {
	t.Run("PositiveNumbers", func(t *testing.T) {
		assert.Equal(t, 1, Factorial(0))
		assert.Equal(t, 1, Factorial(1))
		assert.Equal(t, 2, Factorial(2))
		assert.Equal(t, 6, Factorial(3))
		assert.Equal(t, 24, Factorial(4))
		assert.Equal(t, 120, Factorial(5))
	})

	t.Run("PanicOnNegative", func(t *testing.T) {
		require.Panics(t, func() {
			Factorial(-1)
		})

		require.Panics(t, func() {
			Factorial(-10)
		})

		// Проверка текста паники
		require.PanicsWithValue(t, "factorial of negative number", func() {
			Factorial(-5)
		})
	})
}

// Тесты для ValidatePositive
func TestValidatePositive(t *testing.T) {
	t.Run("ValidNumbers", func(t *testing.T) {
		// Не должно паниковать
		assert.NotPanics(t, func() {
			ValidatePositive(1)
		})

		assert.NotPanics(t, func() {
			ValidatePositive(100)
		})
	})

	t.Run("PanicOnNonPositive", func(t *testing.T) {
		require.Panics(t, func() {
			ValidatePositive(0)
		})

		require.Panics(t, func() {
			ValidatePositive(-1)
		})

		require.Panics(t, func() {
			ValidatePositive(-100)
		})

		// Проверка текста паники
		require.PanicsWithValue(t, "number must be positive, got 0", func() {
			ValidatePositive(0)
		})

		require.PanicsWithValue(t, "number must be positive, got -5", func() {
			ValidatePositive(-5)
		})
	})
}

// Табличные тесты для Factorial
func TestFactorial_TableDriven(t *testing.T) {
	testCases := []struct {
		name        string
		input       int
		expected    int
		shouldPanic bool
		panicValue  string
	}{
		{"Zero", 0, 1, false, ""},
		{"One", 1, 1, false, ""},
		{"Two", 2, 2, false, ""},
		{"Three", 3, 6, false, ""},
		{"Four", 4, 24, false, ""},
		{"Five", 5, 120, false, ""},
		{"NegativeOne", -1, 0, true, "factorial of negative number"},
		{"NegativeTen", -10, 0, true, "factorial of negative number"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPanic {
				require.PanicsWithValue(t, tc.panicValue, func() {
					Factorial(tc.input)
				})
			} else {
				result := Factorial(tc.input)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

// Табличные тесты для ValidatePositive
func TestValidatePositive_TableDriven(t *testing.T) {
	testCases := []struct {
		name        string
		input       int
		shouldPanic bool
		panicValue  string
	}{
		{"PositiveOne", 1, false, ""},
		{"PositiveHundred", 100, false, ""},
		{"Zero", 0, true, "number must be positive, got 0"},
		{"NegativeOne", -1, true, "number must be positive, got -1"},
		{"NegativeFive", -5, true, "number must be positive, got -5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPanic {
				require.PanicsWithValue(t, tc.panicValue, func() {
					ValidatePositive(tc.input)
				})
			} else {
				assert.NotPanics(t, func() {
					ValidatePositive(tc.input)
				})
			}
		})
	}
}

// Бенчмарки
func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Sum(123, 456)
	}
}

func BenchmarkDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Divide(100, 5)
	}
}

func BenchmarkDivide_Error(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Divide(100, 0)
	}
}

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Factorial(10)
	}
}

func BenchmarkValidatePositive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidatePositive(1)
	}
}