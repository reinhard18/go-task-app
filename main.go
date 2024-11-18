package main

import (
	"log"
	"net/http"
	"os"

	"task-management/db"
	"task-management/handlers"
	"task-management/middleware"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	dbURL := os.Getenv("DATABASE_URL")

	db, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// CORS configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		Debug:          true,
	})

	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", handlers.AuthenticateHandler(db, jwtKey))
	mux.HandleFunc("/api/tasks", middleware.Authenticate(handlers.TaskHandler(db), jwtKey))
	mux.HandleFunc("/api/task/", middleware.Authenticate(handlers.TaskHandler(db), jwtKey))
	mux.HandleFunc("/api/register", handlers.RegisterHandler(db, jwtKey))
	mux.HandleFunc("/api/health", handlers.HealthCheckHandler)
	handler := corsHandler.Handler(mux)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
