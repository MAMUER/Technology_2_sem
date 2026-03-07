package service

import (
	"sync"

	"go.uber.org/zap"
	"tech-ip-sem2/shared/logger"
)

type User struct {
	Username string
	Password string
}

type AuthService struct {
	users map[string]User
	mu    sync.RWMutex
	log   *logger.Logger
}

func NewAuthService(log *logger.Logger) *AuthService {
	users := map[string]User{
		"student": {
			Username: "student",
			Password: "student",
		},
		"admin": {
			Username: "admin",
			Password: "admin123",
		},
	}

	log.Info("Auth service initialized", zap.Int("users_count", len(users)))

	return &AuthService{
		users: users,
		log:   log,
	}
}

func (s *AuthService) ValidateCredentials(username, password string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		s.log.Debug("user not found", zap.String("username", username))
		return false
	}

	valid := user.Password == password
	if !valid {
		s.log.Debug("invalid password", zap.String("username", username))
	}

	return valid
}

func (s *AuthService) ValidateToken(token string) (bool, string) {
	switch token {
	case "demo-token-for-student":
		s.log.Debug("token validated", zap.String("subject", "student"))
		return true, "student"
	case "demo-token-for-admin":
		s.log.Debug("token validated", zap.String("subject", "admin"))
		return true, "admin"
	default:
		s.log.Debug("invalid token", zap.String("token_prefix", token[:min(10, len(token))]))
		return false, ""
	}
}
