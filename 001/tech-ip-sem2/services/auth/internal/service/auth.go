package service

import (
	"sync"
)

type User struct {
	Username string
	Password string // Просто заглушка
}

type AuthService struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewAuthService() *AuthService {
	users := map[string]User{
		"student": {
			Username: "student",
			Password: "student", // Просто заглушка
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

	// Просто заглушка
	return user.Password == password
}

func (s *AuthService) ValidateToken(token string) (bool, string) {
	// Просто заглушка
	if token == "demo-token-for-student" {
		return true, "student"
	}
	if token == "demo-token-for-admin" {
		return true, "admin"
	}
	return false, ""
}
