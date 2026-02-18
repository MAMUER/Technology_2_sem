package service

import (
	"sync"
)

type User struct {
	Username string
	Password string // В реальном проекте должен быть хэш
}

type AuthService struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewAuthService() *AuthService {
	// Инициализация с тестовыми пользователями
	users := map[string]User{
		"student": {
			Username: "student",
			Password: "student", // В реальном проекте должен быть хэш
		},
		"admin": {
			Username: "admin",
			Password: "admin123",
		},
	}

	return &AuthService{
		users: users,
	}
}

func (s *AuthService) ValidateCredentials(username, password string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return false
	}

	// В реальном проекте здесь должно быть сравнение хэшей
	return user.Password == password
}

func (s *AuthService) ValidateToken(token string) (bool, string) {
	// Упрощенная проверка токена
	// В реальном проекте здесь должна быть проверка JWT
	if token == "demo-token-for-student" {
		return true, "student"
	}
	if token == "demo-token-for-admin" {
		return true, "admin"
	}
	return false, ""
}
