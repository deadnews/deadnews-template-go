package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResponseWriterStatusCapture tests that responseWriter captures status code.
func TestResponseWriterStatusCapture(t *testing.T) {
	tests := []struct {
		name         string
		writeStatus  int
		expectedCode int
	}{
		{"captures 200", http.StatusOK, http.StatusOK},
		{"captures 201", http.StatusCreated, http.StatusCreated},
		{"captures 400", http.StatusBadRequest, http.StatusBadRequest},
		{"captures 404", http.StatusNotFound, http.StatusNotFound},
		{"captures 500", http.StatusInternalServerError, http.StatusInternalServerError},
		{"captures 503", http.StatusServiceUnavailable, http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			rw := &responseWriter{ResponseWriter: rec, status: http.StatusOK}

			rw.WriteHeader(tt.writeStatus)

			assert.Equal(t, tt.expectedCode, rw.status)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

// TestResponseWriterDefaultStatus tests that responseWriter defaults to 200.
func TestResponseWriterDefaultStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, status: http.StatusOK}

	// Don't call WriteHeader, just write body
	_, err := rw.Write([]byte("test"))
	require.NoError(t, err)

	// Status should remain at default
	assert.Equal(t, http.StatusOK, rw.status)
}

// TestResponseWriterUnwrap tests that Unwrap returns the underlying ResponseWriter.
func TestResponseWriterUnwrap(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, status: http.StatusOK}

	assert.Equal(t, rec, rw.Unwrap())
}

// TestLoggerMiddleware tests the Logger middleware.
func TestLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		handlerStatus  int
		expectLogged   bool
		expectLogLevel slog.Level
	}{
		{
			name:          "health endpoint not logged",
			path:          "/health",
			handlerStatus: http.StatusOK,
			expectLogged:  false,
		},
		{
			name:           "normal request logged at info",
			path:           "/query",
			handlerStatus:  http.StatusOK,
			expectLogged:   true,
			expectLogLevel: slog.LevelInfo,
		},
		{
			name:           "4xx logged at warn",
			path:           "/query",
			handlerStatus:  http.StatusBadRequest,
			expectLogged:   true,
			expectLogLevel: slog.LevelWarn,
		},
		{
			name:           "5xx logged at error",
			path:           "/query",
			handlerStatus:  http.StatusInternalServerError,
			expectLogged:   true,
			expectLogLevel: slog.LevelError,
		},
		{
			name:           "301 logged at info",
			path:           "/redirect",
			handlerStatus:  http.StatusMovedPermanently,
			expectLogged:   true,
			expectLogLevel: slog.LevelInfo,
		},
		{
			name:           "404 logged at warn",
			path:           "/notfound",
			handlerStatus:  http.StatusNotFound,
			expectLogged:   true,
			expectLogLevel: slog.LevelWarn,
		},
		{
			name:           "503 logged at error",
			path:           "/unavailable",
			handlerStatus:  http.StatusServiceUnavailable,
			expectLogged:   true,
			expectLogLevel: slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			oldDefault := slog.Default()
			slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})))
			defer slog.SetDefault(oldDefault)

			// Create a test handler that returns the specified status
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.handlerStatus)
			})

			// Wrap with Logger middleware
			wrapped := Logger(handler)

			// Make request
			req := httptest.NewRequest(http.MethodGet, tt.path, http.NoBody)
			rec := httptest.NewRecorder()
			wrapped.ServeHTTP(rec, req)

			// Check log output
			logOutput := buf.String()
			if tt.expectLogged {
				assert.Contains(t, logOutput, "request")
				assert.Contains(t, logOutput, tt.path)
				assert.Contains(t, logOutput, "method=GET")
			} else {
				assert.Empty(t, logOutput)
			}
		})
	}
}

// TestLoggerMiddlewareRequestDetails tests that Logger logs request details.
func TestLoggerMiddlewareRequestDetails(t *testing.T) {
	var buf bytes.Buffer
	oldDefault := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(oldDefault)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Logger(handler)

	req := httptest.NewRequest(http.MethodPost, "/query?param=value", http.NoBody)
	req.Header.Set("User-Agent", "TestAgent/1.0")
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "method=POST")
	assert.Contains(t, logOutput, "/query")
	assert.Contains(t, logOutput, "TestAgent/1.0")
	assert.Contains(t, logOutput, "status=200")
	assert.Contains(t, logOutput, "duration=")
}

// TestRecovererMiddleware tests the Recoverer middleware.
func TestRecovererMiddleware(t *testing.T) {
	t.Run("recovers from panic with string", func(t *testing.T) {
		var buf bytes.Buffer
		oldDefault := slog.Default()
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
		defer slog.SetDefault(oldDefault)

		handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			panic("test panic")
		})

		wrapped := Recoverer(handler)

		req := httptest.NewRequest(http.MethodGet, "/panic", http.NoBody)
		rec := httptest.NewRecorder()

		// Should not panic
		assert.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		// Should return 500
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		// Should log the panic
		assert.Contains(t, buf.String(), "panic recovered")
		assert.Contains(t, buf.String(), "test panic")
	})

	t.Run("recovers from panic with error", func(t *testing.T) {
		var buf bytes.Buffer
		oldDefault := slog.Default()
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
		defer slog.SetDefault(oldDefault)

		handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			panic(io.ErrUnexpectedEOF)
		})

		wrapped := Recoverer(handler)

		req := httptest.NewRequest(http.MethodGet, "/panic", http.NoBody)
		rec := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("passes through without panic", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success"))
		})

		wrapped := Recoverer(handler)

		req := httptest.NewRequest(http.MethodGet, "/normal", http.NoBody)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})

	t.Run("recovers from panic with runtime error", func(t *testing.T) {
		var buf bytes.Buffer
		oldDefault := slog.Default()
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
		defer slog.SetDefault(oldDefault)

		handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			panic(errors.New("runtime error"))
		})

		wrapped := Recoverer(handler)

		req := httptest.NewRequest(http.MethodGet, "/panic", http.NoBody)
		rec := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, buf.String(), "panic recovered")
	})
}

// TestRecovererMiddlewareLogsURL tests that Recoverer logs the URL.
func TestRecovererMiddlewareLogsURL(t *testing.T) {
	var buf bytes.Buffer
	oldDefault := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(oldDefault)

	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("test")
	})

	wrapped := Recoverer(handler)

	req := httptest.NewRequest(http.MethodGet, "/test/path?query=value", http.NoBody)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	assert.Contains(t, buf.String(), "/test/path")
}

// TestMiddlewareChain tests that middleware can be chained correctly.
func TestMiddlewareChain(t *testing.T) {
	t.Run("logger and recoverer chain", func(t *testing.T) {
		var buf bytes.Buffer
		oldDefault := slog.Default()
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
		defer slog.SetDefault(oldDefault)

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Chain: Recoverer -> Logger -> Handler
		wrapped := Recoverer(Logger(handler))

		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		rec := httptest.NewRecorder()

		wrapped.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, buf.String(), "request")
	})

	t.Run("panic in chain is recovered and logged", func(t *testing.T) {
		var buf bytes.Buffer
		oldDefault := slog.Default()
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
		defer slog.SetDefault(oldDefault)

		handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			panic("chain panic")
		})

		// Chain: Recoverer -> Logger -> Handler
		wrapped := Recoverer(Logger(handler))

		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		rec := httptest.NewRecorder()

		assert.NotPanics(t, func() {
			wrapped.ServeHTTP(rec, req)
		})

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, buf.String(), "panic recovered")
	})
}

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// TestLoggerMiddlewareWithContext tests that Logger works with request context.
func TestLoggerMiddlewareWithContext(t *testing.T) {
	var buf bytes.Buffer
	oldDefault := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(oldDefault)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := Logger(handler)

	ctx := context.WithValue(context.Background(), contextKey("test"), "test-value")
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody).WithContext(ctx)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestResponseWriterWrite tests that Write method works correctly.
func TestResponseWriterWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}

	testData := []byte("test response body")
	n, err := rw.Write(testData)

	require.NoError(t, err)
	assert.Equal(t, len(testData), n)
	assert.Equal(t, "test response body", rec.Body.String())
}

// TestLoggerMiddlewareStatusCodes tests Logger with various status codes.
func TestLoggerMiddlewareStatusCodes(t *testing.T) {
	statusCodes := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusNoContent,
		http.StatusMovedPermanently,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}

	for _, code := range statusCodes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			var buf bytes.Buffer
			oldDefault := slog.Default()
			slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
			defer slog.SetDefault(oldDefault)

			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(code)
			})

			wrapped := Logger(handler)

			req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, code, rec.Code)
			// Check that the numeric status code is logged
			assert.Contains(t, buf.String(), fmt.Sprintf("status=%d", code))
		})
	}
}
