package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Global variables for test environment.
var (
	testDSN       string
	testContext   context.Context
	testContainer testcontainers.Container
)

// TestMain sets up and tears down the shared test environment.
func TestMain(m *testing.M) {
	// Skip container setup if TESTCONTAINERS is not set.
	if os.Getenv("TESTCONTAINERS") != "1" {
		os.Exit(m.Run())
	}

	// Create context
	testContext = context.Background()

	// Create container request
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		),
	}

	// Start container
	pgContainer, err := testcontainers.GenericContainer(testContext, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("Failed to start container: %s\n", err)
		os.Exit(1)
	}
	testContainer = pgContainer

	// Get the mapped port
	port, err := pgContainer.MappedPort(testContext, "5432")
	if err != nil {
		fmt.Printf("Failed to get port: %s\n", err)
		terminateContainer()
		os.Exit(1)
	}

	// Get the host
	host, err := pgContainer.Host(testContext)
	if err != nil {
		fmt.Printf("Failed to get host: %s\n", err)
		terminateContainer()
		os.Exit(1)
	}

	// Construct DSN
	testDSN = fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	fmt.Println("Using test DSN: ", testDSN)

	// Initialize the global database connection pool
	if err := InitDB(testDSN); err != nil {
		fmt.Printf("Failed to initialize database: %s\n", err)
		terminateContainer()
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	// Clean up
	CloseDB()
	terminateContainer()

	// Exit with the appropriate code
	os.Exit(code)
}

// Helper function to terminate the container.
func terminateContainer() {
	if err := testContainer.Terminate(testContext); err != nil {
		fmt.Printf("Error terminating container: %s\n", err)
	}
}

// skipIfNoTestcontainers skips the test if testcontainers are not enabled.
func skipIfNoTestcontainers(t *testing.T) {
	t.Helper()
	if os.Getenv("TESTCONTAINERS") != "1" {
		t.Skip("Skipping integration test, set TESTCONTAINERS=1 to run it.")
	}
}

func TestDatabaseService(t *testing.T) {
	skipIfNoTestcontainers(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test getting database info directly (uses global db pool)
	dbInfo, err := getDatabaseInfo(ctx)
	require.NoError(t, err)

	// Verify expected fields are present
	assert.Contains(t, dbInfo, "database")
	assert.Contains(t, dbInfo, "version")

	// Verify database name is correct
	assert.Equal(t, "testdb", dbInfo["database"])

	// Verify version contains PostgreSQL
	version, ok := dbInfo["version"].(string)
	assert.True(t, ok, "version should be a string")
	assert.Contains(t, version, "PostgreSQL", "version should contain PostgreSQL")
}

func TestSetupServer(t *testing.T) {
	server := setupServer()
	assert.NotNil(t, server)
	assert.Equal(t, ":8000", server.Addr)
	assert.NotNil(t, server.Handler)
}

func TestHandleDatabaseTest_Success(t *testing.T) {
	skipIfNoTestcontainers(t)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handleDatabaseTest(rr, req)

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

	server := setupServer()
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

func TestHealthEndpoint(t *testing.T) {
	server := setupServer()
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+"/health", http.NoBody)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestServerConfiguration tests server timeout configuration.
func TestServerConfiguration(t *testing.T) {
	server := setupServer()

	assert.Equal(t, 15*time.Second, server.ReadTimeout)
	assert.Equal(t, 15*time.Second, server.WriteTimeout)
	assert.Equal(t, 60*time.Second, server.IdleTimeout)
}

// TestHandleDatabaseTestContextCancellation tests context cancellation.
func TestHandleDatabaseTestContextCancellation(t *testing.T) {
	skipIfNoTestcontainers(t)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/test", http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handleDatabaseTest(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Internal server error")
}

// TestSetupServerMiddlewareConfiguration tests that all middleware is properly configured.
func TestSetupServerMiddlewareConfiguration(t *testing.T) {
	server := setupServer()

	// Test that health endpoint is configured
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/health", http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestSetupServerNonExistentRoute tests 404 handling.
func TestSetupServerNonExistentRoute(t *testing.T) {
	server := setupServer()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/nonexistent", http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestHealthEndpointDirect tests the handleHealth function directly.
func TestHealthEndpointDirect(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/health", http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handleHealth(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestSetupServerMethodNotAllowed tests that wrong HTTP methods are rejected.
func TestSetupServerMethodNotAllowed(t *testing.T) {
	server := setupServer()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"POST to /health", http.MethodPost, "/health"},
		{"PUT to /health", http.MethodPut, "/health"},
		{"DELETE to /test", http.MethodDelete, "/test"},
		{"POST to /test", http.MethodPost, "/test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), tt.method, tt.path, http.NoBody)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			server.Handler.ServeHTTP(rr, req)

			// Standard library returns 405 Method Not Allowed for wrong methods on defined routes
			assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		})
	}
}

// TestSetupServerReturnsValidServer tests that setupServer returns a properly configured server.
func TestSetupServerReturnsValidServer(t *testing.T) {
	server := setupServer()

	assert.NotNil(t, server)
	assert.NotNil(t, server.Handler)
	assert.Equal(t, ":8000", server.Addr)
	assert.Greater(t, server.ReadTimeout, time.Duration(0))
	assert.Greater(t, server.WriteTimeout, time.Duration(0))
	assert.Greater(t, server.IdleTimeout, time.Duration(0))
}
