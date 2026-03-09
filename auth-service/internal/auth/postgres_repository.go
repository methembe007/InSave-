package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *PostgresRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	user.ID = uuid.New().String()
	user.CreatedAt = user.CreatedAt
	user.UpdatedAt = user.UpdatedAt

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, date_of_birth, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, date_of_birth, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// EmailExists checks if an email is already registered
func (r *PostgresRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
