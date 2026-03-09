package storage

import (
	"sync"
	"time"
)

type ProcessedMessages struct {
	mu        sync.RWMutex
	messages  map[string]bool
	createdAt map[string]time.Time
	ttl       time.Duration
}

func NewProcessedMessages(ttl time.Duration) *ProcessedMessages {
	store := &ProcessedMessages{
		messages:  make(map[string]bool),
		createdAt: make(map[string]time.Time),
		ttl:       ttl,
	}

	// Запускаем горутину для очистки старых записей
	go store.cleanup()

	return store
}

func (s *ProcessedMessages) IsProcessed(messageID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	processed, exists := s.messages[messageID]
	if !exists {
		return false
	}

	// Проверяем TTL
	if createdAt, ok := s.createdAt[messageID]; ok {
		if time.Since(createdAt) > s.ttl {
			// Запись устарела, удалим её при следующей очистке
			return false
		}
	}

	return processed
}

func (s *ProcessedMessages) MarkProcessed(messageID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages[messageID] = true
	s.createdAt[messageID] = time.Now()
}

func (s *ProcessedMessages) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, createdAt := range s.createdAt {
			if now.Sub(createdAt) > s.ttl {
				delete(s.messages, id)
				delete(s.createdAt, id)
			}
		}
		s.mu.Unlock()
	}
}
