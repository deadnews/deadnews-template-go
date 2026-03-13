package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Global variables for test environment.
var testPool *pgxpool.Pool

// testApp creates an App using the shared test pool.
func testApp(t *testing.T) *App {
	t.Helper()
	return &App{DB: testPool}
}

// TestMain sets up and tears down the shared test environment.
func TestMain(m *testing.M) {
	if os.Getenv("TESTCONTAINERS") != "1" {
		os.Exit(m.Run())
	}

	ctx := context.Background()

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

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Printf("Failed to start container: %s\n", err)
		os.Exit(1)
	}

	terminate := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			fmt.Printf("Error terminating container: %s\n", err)
		}
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Printf("Failed to get port: %s\n", err)
		terminate()
		os.Exit(1)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		fmt.Printf("Failed to get host: %s\n", err)
		terminate()
		os.Exit(1)
	}

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	fmt.Println("Using test DSN:", dsn)

	testPool, err = openDB(dsn)
	if err != nil {
		fmt.Printf("Failed to initialize database: %s\n", err)
		terminate()
		os.Exit(1)
	}

	code := m.Run()

	testPool.Close()
	terminate()
	os.Exit(code)
}

// skipIfNoTestcontainers skips the test if testcontainers are not enabled.
func skipIfNoTestcontainers(t *testing.T) {
	t.Helper()
	if os.Getenv("TESTCONTAINERS") != "1" {
		t.Skip("Skipping integration test, set TESTCONTAINERS=1 to run it.")
	}
}

func TestSetupServer(t *testing.T) {
	app := testApp(t)
	server := setupServer(app)

	assert.NotNil(t, server)
	assert.NotNil(t, server.Handler)
}

// TestSetupServerNonExistentRoute tests 404 handling.
func TestSetupServerNonExistentRoute(t *testing.T) {
	app := testApp(t)
	server := setupServer(app)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/nonexistent", http.NoBody)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	server.Handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// TestSetupServerMethodNotAllowed tests that wrong HTTP methods are rejected.
func TestSetupServerMethodNotAllowed(t *testing.T) {
	app := testApp(t)
	server := setupServer(app)

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
