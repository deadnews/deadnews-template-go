// Package main is the entry point for the template-go application.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Setup structured logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Get database DSN from environment
	dsn := os.Getenv("SERVICE_DSN")
	if dsn == "" {
		slog.Error("SERVICE_DSN environment variable is not set")
		os.Exit(1)
	}

	// Initialize database connection pool
	if err := InitDB(dsn); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer CloseDB()

	// Create server
	server := setupServer()

	// Create context that cancels on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Graceful shutdown goroutine
	go func() {
		// Wait for termination signal
		<-ctx.Done()
		slog.Info("Shutdown signal received")

		// Graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown error", "error", err)
		} else {
			slog.Info("Server shutdown completed")
		}
	}()

	// Start server
	slog.Info("Starting server", "port", "8000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server error", "error", err)
	}
}

// setupServer creates a configured HTTP server.
func setupServer() *http.Server {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /test", handleDatabaseTest)

	// Apply middlewares
	handler := Recoverer(Logger(mux))

	return &http.Server{
		Addr:         ":8000",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

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
