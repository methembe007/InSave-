package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/insavein/goal-service/internal/goal"
	"github.com/insavein/goal-service/internal/handlers"
	"github.com/insavein/goal-service/internal/middleware"
	"github.com/insavein/goal-service/pkg/database"
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
	
	// Initialize repository and service
	repo := goal.NewPostgresRepository(db)
	service := goal.NewGoalService(repo)
	
	// Initialize handlers
	goalHandler := handlers.NewGoalHandler(service)
	
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
	
	// Goal routes
	api.HandleFunc("/goals", goalHandler.CreateGoal).Methods("POST")
	api.HandleFunc("/goals", goalHandler.GetActiveGoals).Methods("GET")
	api.HandleFunc("/goals/{id}", goalHandler.GetGoal).Methods("GET")
	api.HandleFunc("/goals/{id}", goalHandler.UpdateGoal).Methods("PUT")
	api.HandleFunc("/goals/{id}", goalHandler.DeleteGoal).Methods("DELETE")
	api.HandleFunc("/goals/{id}/progress", goalHandler.UpdateProgress).Methods("POST")
	api.HandleFunc("/goals/{id}/milestones", goalHandler.GetMilestones).Methods("GET")
	
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
		port = "8005"
	}
	
	log.Printf("Goal Service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
