package main

import (
	"fmt"
	"log"
	"manara/database"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local" // Default to local environment
	}

	envFile := fmt.Sprintf(".env.%s", env)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("❌ Error loading %s file: %v", envFile, err)
	}

	log.Printf("🚀 Starting Manara Backend in %s mode...", env)

	// Connect to database
	err = database.Connect()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer database.Close()

	// Initialize router
	router := mux.NewRouter()

	// Test route
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Welcome to Manara API", "status": "running"}`))
	}).Methods("GET")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "database": "connected"}`))
	}).Methods("GET")

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CORS_ALLOWED_ORIGINS")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Get port from environment
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("✅ Server is running on port %s", port)
	log.Printf("🌍 Environment: %s", env)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler.Handler(router)))
}
