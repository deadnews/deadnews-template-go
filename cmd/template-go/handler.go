package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// handleHealth handles the health check endpoint.
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// handleDatabaseTest handles the test endpoint and returns database information as JSON.
func (app *App) handleDatabaseTest(w http.ResponseWriter, r *http.Request) {
	dbInfo, err := getDatabaseInfo(r.Context(), app.DB)
	if err != nil {
		slog.Error("Failed to get database info", "error", err)
		http.Error(w, fmt.Sprintf("Failed to get database info: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dbInfo); err != nil {
		slog.Error("Failed to write JSON response", "error", err)
	}
}
