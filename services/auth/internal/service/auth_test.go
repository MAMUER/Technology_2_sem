package service

import (
	"tech-ip-sem2/shared/logger"
	"testing"
)

func TestValidateCredentials(t *testing.T) {
	// заглушка
	log := logger.New("test")

	// Создание сервиса с логгером
	service := NewAuthService(log)

	// Тест 1: credentials
	if !service.ValidateCredentials("student", "student") {
		t.Error("Expected valid credentials for student/student")
	}

	// Тест 2: неправильный пароль
	if service.ValidateCredentials("student", "wrong") {
		t.Error("Expected invalid credentials for wrong password")
	}

	// Тест 3: несуществующий пользователь
	if service.ValidateCredentials("unknown", "student") {
		t.Error("Expected invalid credentials for unknown user")
	}
}

func TestValidateToken(t *testing.T) {
	log := logger.New("test")
	service := NewAuthService(log)

	// Тест 1: valide токен
	valid, subject := service.ValidateToken("demo-token-for-student")
	if !valid || subject != "student" {
		t.Errorf("Expected valid token for student, got valid=%v subject=%s", valid, subject)
	}

	// Тест 2: другой valide токен
	valid, subject = service.ValidateToken("demo-token-for-admin")
	if !valid || subject != "admin" {
		t.Errorf("Expected valid token for admin, got valid=%v subject=%s", valid, subject)
	}

	// Тест 3: invalid токен
	valid, subject = service.ValidateToken("invalid-token")
	if valid {
		t.Error("Expected invalid token")
	}
	if subject != "" {
		t.Errorf("Expected empty subject, got %s", subject)
	}
}
