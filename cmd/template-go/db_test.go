package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDatabaseInfo(t *testing.T) {
	skipIfNoTestcontainers(t)

	t.Run("returns error on context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(1 * time.Millisecond)

		_, err := getDatabaseInfo(ctx, testPool)
		require.Error(t, err)
	})

	t.Run("returns error on cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := getDatabaseInfo(ctx, testPool)
		require.Error(t, err)
	})

	t.Run("returns database info with valid connection", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		dbInfo, err := getDatabaseInfo(ctx, testPool)
		require.NoError(t, err)

		assert.Contains(t, dbInfo, "database")
		assert.Contains(t, dbInfo, "version")

		database, ok := dbInfo["database"].(string)
		assert.True(t, ok, "database should be a string")
		assert.Equal(t, "testdb", database)

		version, ok := dbInfo["version"].(string)
		assert.True(t, ok, "version should be a string")
		assert.Contains(t, version, "PostgreSQL")
	})
}

func TestOpenDB(t *testing.T) {
	t.Run("returns error for invalid DSN", func(t *testing.T) {
		_, err := openDB("invalid-dsn")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse database config")
	})

	t.Run("returns error for unreachable host", func(t *testing.T) {
		_, err := openDB("postgres://user:pass@127.0.0.1:59999/db?sslmode=disable&connect_timeout=1")
		require.Error(t, err)
	})
}
