package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// JSON response structs.
type errorResponse struct {
	Error string `json:"error"`
}

// respondJSON encodes data as JSON and writes it with the given status code.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to write JSON response", "error", err)
	}
}

// respondError writes a JSON error response and logs server errors.
func respondError(w http.ResponseWriter, status int, message string) {
	if status >= 500 {
		slog.Error("handler error", "status", status, "error", message)
	}
	respondJSON(w, status, errorResponse{Error: message})
}

// handleHealth handles the health check endpoint.
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// handleDatabaseTest handles the test endpoint and returns database information as JSON.
func (app *App) handleDatabaseTest(w http.ResponseWriter, r *http.Request) {
	dbInfo, err := getDatabaseInfo(r.Context(), app.DB)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dbInfo)
}
