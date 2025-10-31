package mathx

import (
	"errors"
	"fmt"
)

func Sum(a, b int) int {
	return a + b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("divide by zero")
	}
	return a / b, nil
}

// Factorial вычисляет факториал числа
// Panic-ит если n < 0
func Factorial(n int) int {
	if n < 0 {
		panic("factorial of negative number")
	}
	if n == 0 {
		return 1
	}
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
	}
	return result
}

// ValidatePositive проверяет что число положительное
// Panic-ит если число <= 0
func ValidatePositive(n int) {
	if n <= 0 {
		panic(fmt.Sprintf("number must be positive, got %d", n))
	}
}
