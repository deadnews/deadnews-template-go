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
		ctx, cancel := context.WithTimeout(t.Context(), 1*time.Nanosecond)
		defer cancel()
		time.Sleep(1 * time.Millisecond)

		_, err := getDatabaseInfo(ctx, testPool)
		require.Error(t, err)
	})

	t.Run("returns error on cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		_, err := getDatabaseInfo(ctx, testPool)
		require.Error(t, err)
	})

	t.Run("returns database info with valid connection", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()

		dbInfo, err := getDatabaseInfo(ctx, testPool)
		require.NoError(t, err)

		assert.Equal(t, "testdb", dbInfo.Database)
		assert.Contains(t, dbInfo.Version, "PostgreSQL")
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

	t.Run("success with valid DSN", func(t *testing.T) {
		skipIfNoTestcontainers(t)

		pool, err := openDB(testDSN)
		require.NoError(t, err)
		require.NotNil(t, pool)

		pool.Close()
	})
}
