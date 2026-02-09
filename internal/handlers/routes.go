package handlers

import (
	"net/http"
)

// RegisterRoutes registers the HTTP routes for the application
func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", Ping)
	return mux
}
