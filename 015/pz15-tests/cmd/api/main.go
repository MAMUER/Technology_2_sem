package main

import (
	"fmt"
	"log"

	"pz15-tests/internal/mathx"
	"pz15-tests/internal/service"
	"pz15-tests/internal/stringsx"
)

// Простая реализация репозитория для демонстрации
type demoRepo struct{}

func (d demoRepo) ByEmail(email string) (service.UserRecord, error) {
	// Демо данные
	users := map[string]service.UserRecord{
		"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: "demo_hash"},
		"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: "demo_hash"},
	}

	user, exists := users[email]
	if !exists {
		return service.UserRecord{}, service.ErrNotFound
	}
	return user, nil
}

func main() {
	fmt.Println("=== Math Operations ===")

	// Тестируем mathx
	sum := mathx.Sum(10, 5)
	fmt.Printf("Sum(10, 5) = %d\n", sum)

	result, err := mathx.Divide(20, 4)
	if err != nil {
		log.Printf("Division error: %v", err)
	} else {
		fmt.Printf("Divide(20, 4) = %d\n", result)
	}

	// Тестируем деление на ноль
	_, err = mathx.Divide(10, 0)
	if err != nil {
		fmt.Printf("Divide(10, 0) error: %v\n", err)
	}

	fmt.Println("\n=== String Operations ===")

	// Тестируем stringsx.Clip
	clipped := stringsx.Clip("hello world", 5)
	fmt.Printf("Clip('hello world', 5) = '%s'\n", clipped)

	clipped = stringsx.Clip("test", 10)
	fmt.Printf("Clip('test', 10) = '%s'\n", clipped)

	clipped = stringsx.Clip("short", 3)
	fmt.Printf("Clip('short', 3) = '%s'\n", clipped)

	fmt.Println("\n=== Service Operations ===")

	// Тестируем service
	repo := demoRepo{}
	svc := service.New(repo)

	// Поиск существующих пользователей
	id, err := svc.FindIDByEmail("admin@example.com")
	if err != nil {
		fmt.Printf("Error finding admin: %v\n", err)
	} else {
		fmt.Printf("Found admin with ID: %d\n", id)
	}

	id, err = svc.FindIDByEmail("user@example.com")
	if err != nil {
		fmt.Printf("Error finding user: %v\n", err)
	} else {
		fmt.Printf("Found user with ID: %d\n", id)
	}

	// Поиск несуществующего пользователя
	id, err = svc.FindIDByEmail("nonexistent@example.com")
	if err != nil {
		fmt.Printf("User not found (expected): %v\n", err)
	} else {
		fmt.Printf("Found user with ID: %d\n", id)
	}

	fmt.Println("\n=== Application Started ===")
}
