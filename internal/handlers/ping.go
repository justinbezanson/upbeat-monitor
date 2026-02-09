package handlers

import (
	"encoding/json"
	"net/http"
)

// Response is a simple JSON response structure
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Ping is a simple health check endpoint
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Pong",
		Status:  "OK",
	}
	json.NewEncoder(w).Encode(response)
}
