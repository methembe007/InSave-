package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/insavein/education-service/internal/education"
	"github.com/insavein/education-service/internal/handlers"
	"github.com/insavein/education-service/internal/middleware"
	"github.com/insavein/education-service/pkg/database"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	
	// Connect to primary database
	db, err := database.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to primary database: %v", err)
	}
	defer db.Close()
	
	log.Println("Connected to primary database successfully")
	
	// Connect to read replica
	replicaDB, err := database.ConnectPostgresReplica()
	if err != nil {
		log.Fatalf("Failed to connect to replica database: %v", err)
	}
	defer replicaDB.Close()
	
	log.Println("Connected to replica database successfully")
	
	// Initialize repository and service
	repo := education.NewPostgresRepository(db, replicaDB)
	service := education.NewEducationService(repo)
	
	// Initialize handlers
	educationHandler := handlers.NewEducationHandler(service)
	
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
	
	// Education routes
	api.HandleFunc("/education/lessons", educationHandler.GetLessons).Methods("GET")
	api.HandleFunc("/education/lessons/{id}", educationHandler.GetLesson).Methods("GET")
	api.HandleFunc("/education/lessons/{id}/complete", educationHandler.MarkLessonComplete).Methods("POST")
	api.HandleFunc("/education/progress", educationHandler.GetUserProgress).Methods("GET")
	
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
		port = "8085"
	}
	
	log.Printf("Education Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
