package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// db holds the database connection pool.
var db *pgxpool.Pool

// InitDB initializes the database connection pool.
func InitDB(dsn string) error {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	config.MaxConns = 4
	config.MinConns = 1
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = time.Minute

	// Create pool with config
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection is working
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db = pool
	return nil
}

// CloseDB closes the database connection pool.
func CloseDB() {
	if db != nil {
		db.Close()
		db = nil
	}
}

// getDatabaseInfo queries the database and returns current database name and version.
func getDatabaseInfo(ctx context.Context) (map[string]any, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	// Get database name and version
	var dbName, version string

	err := db.QueryRow(ctx, "SELECT current_database()").Scan(&dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to get current database name: %w", err)
	}

	err = db.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("failed to get database version: %w", err)
	}

	return map[string]any{
		"database": dbName,
		"version":  version,
	}, nil
}
