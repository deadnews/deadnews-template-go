package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	t.Run("returns error for invalid DSN", func(t *testing.T) {
		_, err := NewApp(&Config{DSN: "invalid-dsn"})
		require.Error(t, err)
	})

	t.Run("returns error for unreachable host", func(t *testing.T) {
		_, err := NewApp(&Config{DSN: "postgres://user:pass@127.0.0.1:59999/db?sslmode=disable&connect_timeout=1"})
		require.Error(t, err)
	})

	t.Run("success with valid DSN", func(t *testing.T) {
		skipIfNoTestcontainers(t)

		app, err := NewApp(&Config{DSN: testDSN})
		require.NoError(t, err)
		require.NotNil(t, app)
		require.NotNil(t, app.DB)

		app.Close()
	})
}

func TestAppClose(t *testing.T) {
	t.Run("close with nil DB does not panic", func(t *testing.T) {
		app := &App{DB: nil}
		assert.NotPanics(t, app.Close)
	})

	t.Run("close with valid DB pool", func(t *testing.T) {
		skipIfNoTestcontainers(t)

		app, err := NewApp(&Config{DSN: testDSN})
		require.NoError(t, err)

		assert.NotPanics(t, app.Close)
	})
}
