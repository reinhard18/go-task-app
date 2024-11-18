package handlers

import (
	"net/http"
)

// HealthCheckHandler is a simple health check endpoint.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
