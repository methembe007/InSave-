package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostgresRepository implements Repository interface
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetUserByID retrieves a user by ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT id, email, first_name, last_name, date_of_birth, 
		       profile_image_url, preferences, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	var preferencesJSON []byte
	var profileImageURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&profileImageURL,
		&preferencesJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if profileImageURL.Valid {
		user.ProfileImageURL = profileImageURL.String
	}

	// Parse preferences JSON
	if len(preferencesJSON) > 0 {
		if err := json.Unmarshal(preferencesJSON, &user.Preferences); err != nil {
			return nil, fmt.Errorf("failed to parse preferences: %w", err)
		}
	}

	return &user, nil
}

// UpdateUser updates user profile fields
func (r *PostgresRepository) UpdateUser(ctx context.Context, user *User) error {
	// Validate user ID
	if _, err := uuid.Parse(user.ID); err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Marshal preferences to JSON
	preferencesJSON, err := json.Marshal(user.Preferences)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences: %w", err)
	}

	query := `
		UPDATE users
		SET first_name = $1,
		    last_name = $2,
		    date_of_birth = $3,
		    profile_image_url = $4,
		    preferences = $5,
		    updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		sql.NullString{String: user.ProfileImageURL, Valid: user.ProfileImageURL != ""},
		preferencesJSON,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user and all associated data (cascade)
func (r *PostgresRepository) DeleteUser(ctx context.Context, userID string) error {
	// Validate user ID
	if _, err := uuid.Parse(userID); err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Start transaction for atomicity
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete user (cascade will handle related records due to ON DELETE CASCADE in schema)
	query := `DELETE FROM users WHERE id = $1`
	result, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
