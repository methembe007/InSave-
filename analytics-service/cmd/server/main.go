package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/insavein/analytics-service/internal/analytics"
	"github.com/insavein/analytics-service/internal/handlers"
	"github.com/insavein/analytics-service/internal/middleware"
	"github.com/insavein/analytics-service/pkg/database"
)

func main() {
	// Load database configuration
	dbConfig := database.LoadConfigFromEnv()
	
	// Connect to database (read replica for analytics)
	replicaHost := os.Getenv("DB_REPLICA_HOST")
	if replicaHost != "" {
		dbConfig.Host = replicaHost
	}
	
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	log.Println("Connected to database")
	
	// Create repository
	repo := analytics.NewPostgresRepository(db)
	
	// Create cache
	cache := analytics.NewMemoryCache()
	
	// Create service
	service := analytics.NewService(repo, cache)
	
	// Create handler
	handler := handlers.NewAnalyticsHandler(service)
	
	// Set up router
	router := mux.NewRouter()
	
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")
	
	// API routes with authentication middleware
	api := router.PathPrefix("/api/analytics").Subrouter()
	api.Use(middleware.AuthMiddleware)
	
	api.HandleFunc("/spending", handler.GetSpendingAnalysis).Methods("GET")
	api.HandleFunc("/patterns", handler.GetSavingsPatterns).Methods("GET")
	api.HandleFunc("/recommendations", handler.GetRecommendations).Methods("GET")
	api.HandleFunc("/health", handler.GetFinancialHealth).Methods("GET")
	
	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}
	
	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Analytics service starting on %s", addr)
	
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
