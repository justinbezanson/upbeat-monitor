package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Response is a simple JSON response structure
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message: "Pong",
			Status:  "OK",
		}
		json.NewEncoder(w).Encode(response)
	})

	fmt.Printf("Server starting on port %s\n", port)
	// Test curl: curl http://localhost:8080/ping
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
