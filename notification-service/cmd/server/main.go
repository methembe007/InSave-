package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/insavein/notification-service/internal/handlers"
	"github.com/insavein/notification-service/internal/middleware"
	"github.com/insavein/notification-service/internal/notification"
	"github.com/insavein/notification-service/pkg/database"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Initialize providers
	emailProvider := notification.NewEmailProvider()
	pushProvider := notification.NewPushProvider()

	// Initialize repository and service
	repo := notification.NewPostgresRepository(db)
	service := notification.NewNotificationService(repo, emailProvider, pushProvider)

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(service)

	// Setup router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// API routes with authentication middleware
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Notification routes
	// Requirement 12.4: GET /api/notifications handler
	api.HandleFunc("/notifications", notificationHandler.GetUserNotifications).Methods("GET")
	
	// Requirement 12.5: PUT /api/notifications/:id/read handler
	api.HandleFunc("/notifications/{id}/read", notificationHandler.MarkNotificationAsRead).Methods("PUT")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	log.Printf("Notification Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
