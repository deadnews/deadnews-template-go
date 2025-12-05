package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDatabaseInfoWithNilDB(t *testing.T) {
	// Save and restore db
	savedDB := db
	db = nil
	defer func() { db = savedDB }()

	ctx := context.Background()
	_, err := getDatabaseInfo(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "database not initialized")
}

func TestGetDatabaseInfoWithTimeout(t *testing.T) {
	if os.Getenv("TESTCONTAINERS") != "1" {
		t.Skip("Skipping integration test, set TESTCONTAINERS=1 to run it.")
	}

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for the context to be cancelled
	time.Sleep(1 * time.Millisecond)

	_, err := getDatabaseInfo(ctx)
	require.Error(t, err)
}

func TestGetDatabaseInfoWithValidConnection(t *testing.T) {
	// This test uses the real container setup
	if os.Getenv("TESTCONTAINERS") != "1" {
		t.Skip("Skipping integration test, set TESTCONTAINERS=1 to run it.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbInfo, err := getDatabaseInfo(ctx)
	require.NoError(t, err)

	// Verify all expected fields are present
	assert.Contains(t, dbInfo, "database")
	assert.Contains(t, dbInfo, "version")

	// Verify field types
	database, ok := dbInfo["database"].(string)
	assert.True(t, ok, "database should be a string")
	assert.NotEmpty(t, database)

	version, ok := dbInfo["version"].(string)
	assert.True(t, ok, "version should be a string")
	assert.NotEmpty(t, version)
}

func TestInitDBWithInvalidDSN(t *testing.T) {
	err := InitDB("invalid-dsn")
	require.Error(t, err)
}

func TestCloseDBWhenNil(_ *testing.T) {
	// Save and restore db
	savedDB := db
	db = nil
	defer func() { db = savedDB }()

	// Should not panic
	CloseDB()
}
