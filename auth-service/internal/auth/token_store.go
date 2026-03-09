package auth

import (
	"context"
	"sync"
	"time"
)

// TokenStore defines operations for managing refresh tokens
type TokenStore interface {
	// StoreRefreshToken stores a refresh token with expiry
	StoreRefreshToken(ctx context.Context, userID, token string, expiry time.Duration) error

	// IsRefreshTokenValid checks if a refresh token is valid (not revoked)
	IsRefreshTokenValid(ctx context.Context, userID, token string) (bool, error)

	// RevokeRefreshToken invalidates a refresh token
	RevokeRefreshToken(ctx context.Context, userID, token string) error

	// RevokeAllUserTokens invalidates all tokens for a user
	RevokeAllUserTokens(ctx context.Context, userID string) error
}

// InMemoryTokenStore implements token storage using in-memory map
type InMemoryTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]map[string]*tokenEntry // userID -> token -> entry
}

type tokenEntry struct {
	token     string
	expiresAt time.Time
}

// NewInMemoryTokenStore creates a new in-memory token store
func NewInMemoryTokenStore() *InMemoryTokenStore {
	store := &InMemoryTokenStore{
		tokens: make(map[string]map[string]*tokenEntry),
	}

	// Start cleanup goroutine
	go store.cleanup()

	return store
}

// StoreRefreshToken stores a refresh token with expiry
func (s *InMemoryTokenStore) StoreRefreshToken(ctx context.Context, userID, token string, expiry time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tokens[userID] == nil {
		s.tokens[userID] = make(map[string]*tokenEntry)
	}

	s.tokens[userID][token] = &tokenEntry{
		token:     token,
		expiresAt: time.Now().Add(expiry),
	}

	return nil
}

// IsRefreshTokenValid checks if a refresh token is valid (not revoked)
func (s *InMemoryTokenStore) IsRefreshTokenValid(ctx context.Context, userID, token string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userTokens, exists := s.tokens[userID]
	if !exists {
		return false, nil
	}

	entry, exists := userTokens[token]
	if !exists {
		return false, nil
	}

	// Check if expired
	if time.Now().After(entry.expiresAt) {
		return false, nil
	}

	return true, nil
}

// RevokeRefreshToken invalidates a refresh token
func (s *InMemoryTokenStore) RevokeRefreshToken(ctx context.Context, userID, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userTokens, exists := s.tokens[userID]; exists {
		delete(userTokens, token)

		// Clean up user entry if no tokens left
		if len(userTokens) == 0 {
			delete(s.tokens, userID)
		}
	}

	return nil
}

// RevokeAllUserTokens invalidates all tokens for a user
func (s *InMemoryTokenStore) RevokeAllUserTokens(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, userID)
	return nil
}

// cleanup periodically removes expired tokens
func (s *InMemoryTokenStore) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()

		for userID, userTokens := range s.tokens {
			for token, entry := range userTokens {
				if now.After(entry.expiresAt) {
					delete(userTokens, token)
				}
			}

			// Remove user entry if no tokens left
			if len(userTokens) == 0 {
				delete(s.tokens, userID)
			}
		}

		s.mu.Unlock()
	}
}

// GetTokenCount returns the number of active tokens for a user (for testing/debugging)
func (s *InMemoryTokenStore) GetTokenCount(userID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if userTokens, exists := s.tokens[userID]; exists {
		return len(userTokens)
	}
	return 0
}
