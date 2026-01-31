package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// handleHealth handles the health check endpoint.
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// handleDatabaseTest handles the test endpoint and returns database information as JSON.
func handleDatabaseTest(w http.ResponseWriter, r *http.Request) {
	dbInfo, err := getDatabaseInfo(r.Context())
	if err != nil {
		slog.Error("Failed to get database info", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbInfo); err != nil {
		slog.Error("Failed to write JSON response", "error", err)
	}
}
