package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/justin/upbeat-monitor/internal/handlers"
)

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Default port if not set in .env
	}

	mux := handlers.RegisterRoutes()

	fmt.Printf("Server starting on port %s\n", port)
	// Test curl: curl http://localhost:8080/ping
	log.Fatal(http.ListenAndServe(":"+port, mux))
}