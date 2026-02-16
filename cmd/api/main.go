package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/justin/upbeat-monitor/internal/database"
	"github.com/justin/upbeat-monitor/internal/handlers"
)

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Default port if not set in .env
	}

	// Determine database connection string
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "local.db" // Default to local SQLite file for development
	}

	db, err := database.NewDBConnection(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close() // Ensure the database connection is closed when main exits

	mux := handlers.RegisterRoutes(db) // Pass the database connection to RegisterRoutes

	fmt.Printf("Server starting on port %s\n", port)
	// Test curl: curl http://localhost:8080/ping
	log.Fatal(http.ListenAndServe(":"+port, mux))
}