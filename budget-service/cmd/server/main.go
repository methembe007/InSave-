package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/insavein/budget-service/internal/budget"
	"github.com/insavein/budget-service/internal/handlers"
	"github.com/insavein/budget-service/internal/middleware"
	"github.com/insavein/budget-service/pkg/database"
)

func main() {
	// Load configuration from environment variables
	cfg := loadConfig()

	// Initialize database connection
	db, err := database.NewPostgresConnection(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Initialize dependencies
	repo := budget.NewPostgresRepository(db)
	budgetService := budget.NewService(repo)

	// Initialize handlers
	budgetHandler := handlers.NewBudgetHandler(budgetService)
	authMiddleware := middleware.NewAuthMiddleware()

	// Setup routes
	mux := http.NewServeMux()

	// Protected routes (require authentication)
	mux.HandleFunc("/api/budget", authMiddleware.Authenticate(budgetHandler.CreateBudget))
	mux.HandleFunc("/api/budget/current", authMiddleware.Authenticate(budgetHandler.GetCurrentBudget))
	mux.HandleFunc("/api/budget/", authMiddleware.Authenticate(budgetHandler.UpdateBudget))
	mux.HandleFunc("/api/budget/spending", authMiddleware.Authenticate(budgetHandler.RecordSpending))
	mux.HandleFunc("/api/budget/alerts", authMiddleware.Authenticate(budgetHandler.GetBudgetAlerts))
	mux.HandleFunc("/api/budget/categories", authMiddleware.Authenticate(budgetHandler.GetCategories))
	mux.HandleFunc("/api/budget/summary", authMiddleware.Authenticate(budgetHandler.GetSpendingSummary))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","database":"disconnected"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"budget-service"}`))
	})

	// Liveness probe
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"alive"}`))
	})

	// Readiness probe
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not ready"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Budget service starting on port %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// Config holds application configuration
type Config struct {
	Port       int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	return Config{
		Port:       getEnvAsInt("PORT", 8083),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "insavein_user"),
		DBPassword: getEnv("DB_PASSWORD", "insavein_password"),
		DBName:     getEnv("DB_NAME", "insavein"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
