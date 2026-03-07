package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"tech-ip-sem2/services/auth/internal/models"
	"tech-ip-sem2/shared/logger"
)

type SessionService struct {
	sessions   map[string]*models.Session
	csrfTokens map[string]string // Добавляем хранилище для CSRF токенов
	mu         sync.RWMutex
	log        *logger.Logger
}

func NewSessionService(log *logger.Logger) *SessionService {
	return &SessionService{
		sessions:   make(map[string]*models.Session),
		csrfTokens: make(map[string]string),
		log:        log,
	}
}

// Генерация безопасного случайного токена
func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Создание новой сессии
func (s *SessionService) CreateSession(username, subject string) (sessionID string, csrfToken string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Генерируем ID сессии
	sessionID, err = generateSecureToken(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Генерируем CSRF токен
	csrfToken, err = generateSecureToken(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}

	session := &models.Session{
		ID:        sessionID,
		Username:  username,
		Subject:   subject,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 часа
	}

	// Сохраняем сессию
	s.sessions[sessionID] = session

	// Сохраняем CSRF токен, связанный с сессией
	s.csrfTokens[sessionID] = csrfToken

	s.log.Info("Session created",
		zap.String("username", username),
		zap.String("subject", subject),
		zap.String("session_id", sessionID[:8]+"..."),
	)

	return sessionID, csrfToken, nil
}

// Получение сессии по ID
func (s *SessionService) GetSession(sessionID string) (*models.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	// Проверяем, не истекла ли сессия
	if time.Now().After(session.ExpiresAt) {
		// Вернем ошибку и позволим вызывающему коду удалить сессию
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

// Удаление сессии (logout)
func (s *SessionService) DeleteSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[sessionID]; !exists {
		return fmt.Errorf("session not found")
	}

	delete(s.sessions, sessionID)
	delete(s.csrfTokens, sessionID)

	s.log.Info("Session deleted", zap.String("session_id", sessionID[:8]+"..."))
	return nil
}

// Получение CSRF токена для сессии
func (s *SessionService) GetCSRFToken(sessionID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token, exists := s.csrfTokens[sessionID]
	if !exists {
		return "", fmt.Errorf("CSRF token not found for session")
	}

	// Проверяем, существует ли сессия и не истекла ли она
	session, exists := s.sessions[sessionID]
	if !exists {
		return "", fmt.Errorf("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		return "", fmt.Errorf("session expired")
	}

	return token, nil
}

// Валидация CSRF токена
func (s *SessionService) ValidateCSRF(sessionID, csrfToken string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Проверяем существование сессии
	session, exists := s.sessions[sessionID]
	if !exists {
		return false
	}

	// Проверяем, не истекла ли сессия
	if time.Now().After(session.ExpiresAt) {
		return false
	}

	// Проверяем CSRF токен
	storedToken, exists := s.csrfTokens[sessionID]
	if !exists {
		return false
	}

	return storedToken == csrfToken
}

// Очистка истекших сессий
func (s *SessionService) CleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, session := range s.sessions {
		if now.After(session.ExpiresAt) {
			delete(s.sessions, id)
			delete(s.csrfTokens, id)
			s.log.Debug("Expired session cleaned up", zap.String("session_id", id[:8]+"..."))
		}
	}
}
