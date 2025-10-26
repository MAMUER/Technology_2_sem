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
        {10, 0, 0, true},  // ошибка деления на ноль
        {0, 5, 0, false},  // ноль делится на любое число
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