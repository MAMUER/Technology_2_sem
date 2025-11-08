package repo

import (
	"sync"
	"time"
)

type RefreshStore struct {
	mu      sync.RWMutex
	tokens  map[string]time.Time
	revoked map[string]time.Time
}

func NewRefreshStore() *RefreshStore {
	return &RefreshStore{
		tokens:  make(map[string]time.Time),
		revoked: make(map[string]time.Time),
	}
}

func (s *RefreshStore) Store(token string, exp time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = exp
	return nil
}

func (s *RefreshStore) IsRevoked(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, revoked := s.revoked[token]; revoked {
		return true
	}
	if exp, exists := s.tokens[token]; exists {
		return time.Now().After(exp)
	}
	return true
}

func (s *RefreshStore) Revoke(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revoked[token] = time.Now()
	delete(s.tokens, token)
	return nil
}

func (s *RefreshStore) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, exp := range s.tokens {
		if now.After(exp) {
			delete(s.tokens, token)
		}
	}
	for token, revokedAt := range s.revoked {
		if now.Sub(revokedAt) > 7*24*time.Hour {
			delete(s.revoked, token)
		}
	}
}
