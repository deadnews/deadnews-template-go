package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleDatabaseTest_Success(t *testing.T) {
	skipIfNoTestcontainers(t)

	app := testApp(t)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	app.handleDatabaseTest(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "application/json")

	var response map[string]any
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	// Verify expected fields are present
	assert.Contains(t, response, "database")
	assert.Contains(t, response, "version")

	// Verify database name is correct
	assert.Equal(t, "testdb", response["database"])

	// Verify version contains PostgreSQL
	version, ok := response["version"].(string)
	assert.True(t, ok, "version should be a string")
	assert.Contains(t, version, "PostgreSQL", "version should contain PostgreSQL")
}

func TestHandleDatabaseTest_ViaServer_Success(t *testing.T) {
	skipIfNoTestcontainers(t)

	app := testApp(t)
	server := setupServer(app)
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/test", http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	var response map[string]any
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Contains(t, response, "database")
	assert.Contains(t, response, "version")
}

// TestHandleDatabaseTestContextCancellation tests context cancellation.
func TestHandleDatabaseTestContextCancellation(t *testing.T) {
	skipIfNoTestcontainers(t)

	app := testApp(t)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	app.handleDatabaseTest(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHealthEndpoint(t *testing.T) {
	app := testApp(t)
	server := setupServer(app)

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	rec := httptest.NewRecorder()

	server.Handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
