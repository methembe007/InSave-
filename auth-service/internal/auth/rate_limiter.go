package auth

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RateLimiter defines rate limiting operations
type RateLimiter interface {
	// AllowLogin checks if login attempt is allowed
	AllowLogin(ctx context.Context, email string) (bool, error)

	// RecordFailedLogin records a failed login attempt
	RecordFailedLogin(ctx context.Context, email string) error

	// ResetLoginAttempts resets the login attempt counter
	ResetLoginAttempts(ctx context.Context, email string) error
}

// InMemoryRateLimiter implements rate limiting using in-memory storage
type InMemoryRateLimiter struct {
	mu       sync.RWMutex
	attempts map[string]*loginAttempts
}

type loginAttempts struct {
	count     int
	firstAttempt time.Time
	blockedUntil time.Time
}

const (
	// MaxLoginAttempts is the maximum number of login attempts allowed
	MaxLoginAttempts = 5

	// LoginAttemptWindow is the time window for counting attempts (15 minutes)
	LoginAttemptWindow = 15 * time.Minute

	// BlockDuration is how long to block after exceeding attempts
	BlockDuration = 15 * time.Minute
)

// NewInMemoryRateLimiter creates a new in-memory rate limiter
func NewInMemoryRateLimiter() *InMemoryRateLimiter {
	limiter := &InMemoryRateLimiter{
		attempts: make(map[string]*loginAttempts),
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// AllowLogin checks if login attempt is allowed
func (r *InMemoryRateLimiter) AllowLogin(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	attempts, exists := r.attempts[email]
	if !exists {
		return true, nil
	}

	now := time.Now()

	// Check if currently blocked
	if now.Before(attempts.blockedUntil) {
		return false, nil
	}

	// Check if window has expired
	if now.Sub(attempts.firstAttempt) > LoginAttemptWindow {
		return true, nil
	}

	// Check if max attempts exceeded
	if attempts.count >= MaxLoginAttempts {
		return false, nil
	}

	return true, nil
}

// RecordFailedLogin records a failed login attempt
func (r *InMemoryRateLimiter) RecordFailedLogin(ctx context.Context, email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	attempts, exists := r.attempts[email]

	if !exists {
		r.attempts[email] = &loginAttempts{
			count:        1,
			firstAttempt: now,
		}
		return nil
	}

	// Reset if window expired
	if now.Sub(attempts.firstAttempt) > LoginAttemptWindow {
		r.attempts[email] = &loginAttempts{
			count:        1,
			firstAttempt: now,
		}
		return nil
	}

	// Increment count
	attempts.count++

	// Block if max attempts exceeded
	if attempts.count >= MaxLoginAttempts {
		attempts.blockedUntil = now.Add(BlockDuration)
	}

	return nil
}

// ResetLoginAttempts resets the login attempt counter
func (r *InMemoryRateLimiter) ResetLoginAttempts(ctx context.Context, email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.attempts, email)
	return nil
}

// cleanup periodically removes expired entries
func (r *InMemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		now := time.Now()

		for email, attempts := range r.attempts {
			// Remove if window expired and not blocked
			if now.Sub(attempts.firstAttempt) > LoginAttemptWindow && now.After(attempts.blockedUntil) {
				delete(r.attempts, email)
			}
		}

		r.mu.Unlock()
	}
}

// GetAttemptInfo returns current attempt info for an email (for testing/debugging)
func (r *InMemoryRateLimiter) GetAttemptInfo(email string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	attempts, exists := r.attempts[email]
	if !exists {
		return "No attempts recorded"
	}

	now := time.Now()
	if now.Before(attempts.blockedUntil) {
		return fmt.Sprintf("Blocked until %s", attempts.blockedUntil.Format(time.RFC3339))
	}

	return fmt.Sprintf("Attempts: %d/%d, Window started: %s",
		attempts.count, MaxLoginAttempts, attempts.firstAttempt.Format(time.RFC3339))
}
