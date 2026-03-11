package savings

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type postgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

// CreateTransaction inserts a new savings transaction
func (r *postgresRepository) CreateTransaction(ctx context.Context, tx *SavingsTransaction) error {
	query := `
		INSERT INTO savings_transactions (id, user_id, amount, currency, description, category, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		tx.ID,
		tx.UserID,
		tx.Amount,
		tx.Currency,
		tx.Description,
		tx.Category,
		tx.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	
	return nil
}

// GetTransactionsByUser retrieves all transactions for a user
func (r *postgresRepository) GetTransactionsByUser(ctx context.Context, userID string, limit, offset int) ([]SavingsTransaction, error) {
	query := `
		SELECT id, user_id, amount, currency, description, category, created_at
		FROM savings_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()
	
	var transactions []SavingsTransaction
	for rows.Next() {
		var tx SavingsTransaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Amount,
			&tx.Currency,
			&tx.Description,
			&tx.Category,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}
	
	return transactions, nil
}

// GetTotalSaved calculates the total amount saved by a user
func (r *postgresRepository) GetTotalSaved(ctx context.Context, userID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM savings_transactions
		WHERE user_id = $1
	`
	
	var total float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total saved: %w", err)
	}
	
	return total, nil
}

// GetMonthlyTotal calculates the total saved in a specific month
func (r *postgresRepository) GetMonthlyTotal(ctx context.Context, userID string, month time.Time) (float64, error) {
	// Get first and last day of the month
	firstDay := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second)
	
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM savings_transactions
		WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3
	`
	
	var total float64
	err := r.db.QueryRowContext(ctx, query, userID, firstDay, lastDay).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get monthly total: %w", err)
	}
	
	return total, nil
}

// GetMonthlyAverage calculates the average monthly savings
func (r *postgresRepository) GetMonthlyAverage(ctx context.Context, userID string) (float64, error) {
	query := `
		SELECT COALESCE(AVG(monthly_total), 0)
		FROM (
			SELECT DATE_TRUNC('month', created_at) as month, SUM(amount) as monthly_total
			FROM savings_transactions
			WHERE user_id = $1
			GROUP BY DATE_TRUNC('month', created_at)
		) as monthly_totals
	`
	
	var average float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&average)
	if err != nil {
		return 0, fmt.Errorf("failed to get monthly average: %w", err)
	}
	
	return average, nil
}

// GetLastSavingDate returns the date of the last savings transaction
func (r *postgresRepository) GetLastSavingDate(ctx context.Context, userID string) (time.Time, error) {
	query := `
		SELECT created_at
		FROM savings_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	var lastDate time.Time
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&lastDate)
	if err == sql.ErrNoRows {
		return time.Time{}, nil
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get last saving date: %w", err)
	}
	
	return lastDate, nil
}

// GetAllSavingDates returns all unique dates when user saved money
func (r *postgresRepository) GetAllSavingDates(ctx context.Context, userID string) ([]time.Time, error) {
	query := `
		SELECT DISTINCT DATE(created_at) as save_date
		FROM savings_transactions
		WHERE user_id = $1
		ORDER BY save_date DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query saving dates: %w", err)
	}
	defer rows.Close()
	
	var dates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			return nil, fmt.Errorf("failed to scan date: %w", err)
		}
		dates = append(dates, date)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating dates: %w", err)
	}
	
	return dates, nil
}

// UpdateUserStreak updates the streak information in user preferences
func (r *postgresRepository) UpdateUserStreak(ctx context.Context, userID string, currentStreak, longestStreak int) error {
	query := `
		UPDATE users
		SET preferences = jsonb_set(
			jsonb_set(
				COALESCE(preferences, '{}'::jsonb),
				'{current_streak}',
				$1::text::jsonb
			),
			'{longest_streak}',
			$2::text::jsonb
		)
		WHERE id = $3
	`
	
	_, err := r.db.ExecContext(ctx, query, currentStreak, longestStreak, userID)
	if err != nil {
		return fmt.Errorf("failed to update user streak: %w", err)
	}
	
	return nil
}

// GetUserStreak retrieves the current streak information from user preferences
func (r *postgresRepository) GetUserStreak(ctx context.Context, userID string) (currentStreak, longestStreak int, err error) {
	query := `
		SELECT COALESCE(preferences->>'current_streak', '0')::int,
		       COALESCE(preferences->>'longest_streak', '0')::int
		FROM users
		WHERE id = $1
	`
	
	err = r.db.QueryRowContext(ctx, query, userID).Scan(&currentStreak, &longestStreak)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get user streak: %w", err)
	}
	
	return currentStreak, longestStreak, nil
}
