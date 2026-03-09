package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost is the cost factor for bcrypt hashing
	BcryptCost = 12

	// AccessTokenExpiry is the duration for access token validity (15 minutes)
	AccessTokenExpiry = 15 * time.Minute

	// RefreshTokenExpiry is the duration for refresh token validity (7 days)
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

// AuthService implements the Service interface
type AuthService struct {
	repo         Repository
	jwtSecret    []byte
	rateLimiter  RateLimiter
	tokenStore   TokenStore
}

// NewAuthService creates a new authentication service
func NewAuthService(repo Repository, jwtSecret string, rateLimiter RateLimiter, tokenStore TokenStore) *AuthService {
	return &AuthService{
		repo:        repo,
		jwtSecret:   []byte(jwtSecret),
		rateLimiter: rateLimiter,
		tokenStore:  tokenStore,
	}
}

// Register creates a new user account with hashed password
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Validate password length
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Check if email already exists
	exists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, errors.New("email already in use")
	}

	// Hash password with bcrypt cost factor 12
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("invalid date of birth format, use YYYY-MM-DD: %w", err)
	}

	// Create user
	now := time.Now()
	user := &User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		DateOfBirth:  dob,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token
	if err := s.tokenStore.StoreRefreshToken(ctx, user.ID, refreshToken, RefreshTokenExpiry); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(AccessTokenExpiry.Seconds()),
		User: UserSummary{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// Login authenticates user and returns JWT tokens
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Check rate limit (5 attempts per 15 minutes per email)
	allowed, err := s.rateLimiter.AllowLogin(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}
	if !allowed {
		return nil, errors.New("too many login attempts, please try again later")
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		// Don't reveal whether email exists
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// Record failed attempt
		_ = s.rateLimiter.RecordFailedLogin(ctx, req.Email)
		return nil, errors.New("invalid credentials")
	}

	// Reset rate limit on successful login
	_ = s.rateLimiter.ResetLoginAttempts(ctx, req.Email)

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token
	if err := s.tokenStore.StoreRefreshToken(ctx, user.ID, refreshToken, RefreshTokenExpiry); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(AccessTokenExpiry.Seconds()),
		User: UserSummary{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// RefreshToken issues new access and refresh tokens
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token is in store (not revoked)
	valid, err := s.tokenStore.IsRefreshTokenValid(ctx, claims.UserID, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}
	if !valid {
		return nil, errors.New("refresh token has been revoked")
	}

	// Generate new tokens
	newAccessToken, err := s.generateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Revoke old refresh token and store new one
	if err := s.tokenStore.RevokeRefreshToken(ctx, claims.UserID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	if err := s.tokenStore.StoreRefreshToken(ctx, claims.UserID, newRefreshToken, RefreshTokenExpiry); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(AccessTokenExpiry.Seconds()),
	}, nil
}

// ValidateToken verifies JWT token signature and expiration
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Validate expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("missing expiration claim")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}

	// Extract required fields
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return nil, errors.New("missing or invalid user_id")
	}

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("missing or invalid email")
	}

	// Extract roles (optional)
	var roles []string
	if rolesInterface, ok := claims["roles"].([]interface{}); ok {
		for _, role := range rolesInterface {
			if roleStr, ok := role.(string); ok {
				roles = append(roles, roleStr)
			}
		}
	}

	return &TokenClaims{
		UserID:    userID,
		Email:     email,
		Roles:     roles,
		ExpiresAt: int64(exp),
	}, nil
}

// Logout invalidates the refresh token
func (s *AuthService) Logout(ctx context.Context, userID string, refreshToken string) error {
	if err := s.tokenStore.RevokeRefreshToken(ctx, userID, refreshToken); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

// generateAccessToken creates a new access token
func (s *AuthService) generateAccessToken(userID, email string) (string, error) {
	expiresAt := time.Now().Add(AccessTokenExpiry)

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"roles":   []string{"user"},
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// generateRefreshToken creates a new refresh token
func (s *AuthService) generateRefreshToken(userID, email string) (string, error) {
	expiresAt := time.Now().Add(RefreshTokenExpiry)

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"roles":   []string{"user"},
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
