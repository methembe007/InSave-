package auth

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	users map[string]*User
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users: make(map[string]*User),
	}
}

func (m *MockRepository) CreateUser(ctx context.Context, user *User) error {
	user.ID = "test-user-id"
	m.users[user.Email] = user
	return nil
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *MockRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	for _, user := range m.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	_, exists := m.users[email]
	return exists, nil
}

func TestRegister(t *testing.T) {
	repo := NewMockRepository()
	rateLimiter := NewInMemoryRateLimiter()
	tokenStore := NewInMemoryTokenStore()
	service := NewAuthService(repo, "test-secret", rateLimiter, tokenStore)

	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		req := RegisterRequest{
			Email:       "test@example.com",
			Password:    "password123",
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: "1990-01-01",
		}

		resp, err := service.Register(ctx, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.AccessToken == "" {
			t.Error("expected access token, got empty string")
		}

		if resp.RefreshToken == "" {
			t.Error("expected refresh token, got empty string")
		}

		if resp.User.Email != req.Email {
			t.Errorf("expected email %s, got %s", req.Email, resp.User.Email)
		}
	})

	t.Run("password too short", func(t *testing.T) {
		req := RegisterRequest{
			Email:       "test2@example.com",
			Password:    "short",
			FirstName:   "Jane",
			LastName:    "Doe",
			DateOfBirth: "1990-01-01",
		}

		_, err := service.Register(ctx, req)
		if err == nil {
			t.Error("expected error for short password, got nil")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		req := RegisterRequest{
			Email:       "test@example.com",
			Password:    "password123",
			FirstName:   "Jane",
			LastName:    "Doe",
			DateOfBirth: "1990-01-01",
		}

		_, err := service.Register(ctx, req)
		if err == nil {
			t.Error("expected error for duplicate email, got nil")
		}
	})
}

func TestLogin(t *testing.T) {
	repo := NewMockRepository()
	rateLimiter := NewInMemoryRateLimiter()
	tokenStore := NewInMemoryTokenStore()
	service := NewAuthService(repo, "test-secret", rateLimiter, tokenStore)

	ctx := context.Background()

	// Create a test user
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), BcryptCost)
	testUser := &User{
		ID:           "test-user-id",
		Email:        "test@example.com",
		PasswordHash: string(passwordHash),
		FirstName:    "John",
		LastName:     "Doe",
		CreatedAt:    time.Now(),
	}
	repo.users[testUser.Email] = testUser

	t.Run("successful login", func(t *testing.T) {
		req := LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		resp, err := service.Login(ctx, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.AccessToken == "" {
			t.Error("expected access token, got empty string")
		}

		if resp.RefreshToken == "" {
			t.Error("expected refresh token, got empty string")
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		req := LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		_, err := service.Login(ctx, req)
		if err == nil {
			t.Error("expected error for invalid password, got nil")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		req := LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		_, err := service.Login(ctx, req)
		if err == nil {
			t.Error("expected error for nonexistent user, got nil")
		}
	})
}

func TestValidateToken(t *testing.T) {
	repo := NewMockRepository()
	rateLimiter := NewInMemoryRateLimiter()
	tokenStore := NewInMemoryTokenStore()
	service := NewAuthService(repo, "test-secret", rateLimiter, tokenStore)

	ctx := context.Background()

	t.Run("valid token", func(t *testing.T) {
		token, err := service.generateAccessToken("test-user-id", "test@example.com")
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		claims, err := service.ValidateToken(ctx, token)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if claims.UserID != "test-user-id" {
			t.Errorf("expected user_id test-user-id, got %s", claims.UserID)
		}

		if claims.Email != "test@example.com" {
			t.Errorf("expected email test@example.com, got %s", claims.Email)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := service.ValidateToken(ctx, "invalid-token")
		if err == nil {
			t.Error("expected error for invalid token, got nil")
		}
	})
}

func TestRateLimiter(t *testing.T) {
	limiter := NewInMemoryRateLimiter()
	ctx := context.Background()
	email := "test@example.com"

	t.Run("allow initial attempts", func(t *testing.T) {
		for i := 0; i < MaxLoginAttempts; i++ {
			allowed, err := limiter.AllowLogin(ctx, email)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !allowed {
				t.Errorf("attempt %d should be allowed", i+1)
			}
			limiter.RecordFailedLogin(ctx, email)
		}
	})

	t.Run("block after max attempts", func(t *testing.T) {
		allowed, err := limiter.AllowLogin(ctx, email)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if allowed {
			t.Error("should be blocked after max attempts")
		}
	})

	t.Run("reset allows login again", func(t *testing.T) {
		limiter.ResetLoginAttempts(ctx, email)
		allowed, err := limiter.AllowLogin(ctx, email)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Error("should be allowed after reset")
		}
	})
}

func TestTokenStore(t *testing.T) {
	store := NewInMemoryTokenStore()
	ctx := context.Background()
	userID := "test-user-id"
	token := "test-token"

	t.Run("store and validate token", func(t *testing.T) {
		err := store.StoreRefreshToken(ctx, userID, token, time.Hour)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		valid, err := store.IsRefreshTokenValid(ctx, userID, token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !valid {
			t.Error("token should be valid")
		}
	})

	t.Run("revoke token", func(t *testing.T) {
		err := store.RevokeRefreshToken(ctx, userID, token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		valid, err := store.IsRefreshTokenValid(ctx, userID, token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if valid {
			t.Error("token should be invalid after revocation")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		expiredToken := "expired-token"
		err := store.StoreRefreshToken(ctx, userID, expiredToken, -time.Hour)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		valid, err := store.IsRefreshTokenValid(ctx, userID, expiredToken)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if valid {
			t.Error("expired token should be invalid")
		}
	})
}
