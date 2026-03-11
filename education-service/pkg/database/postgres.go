package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// ConnectPostgres creates a connection to PostgreSQL primary database
func ConnectPostgres() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if sslmode == "" {
		sslmode = "disable"
	}
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	// Set connection pool settings
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	
	return db, nil
}

// ConnectPostgresReplica creates a connection to PostgreSQL read replica
func ConnectPostgresReplica() (*sql.DB, error) {
	host := os.Getenv("DB_REPLICA_HOST")
	port := os.Getenv("DB_REPLICA_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	
	// Fall back to primary if replica not configured
	if host == "" {
		host = os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
	}
	if port == "" {
		port = os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
	}
	if sslmode == "" {
		sslmode = "disable"
	}
	
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open replica database: %w", err)
	}
	
	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping replica database: %w", err)
	}
	
	// Set connection pool settings
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	
	return db, nil
}
