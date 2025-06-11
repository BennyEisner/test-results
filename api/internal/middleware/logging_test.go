package middleware_test

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	_ "time"

	"github.com/BennyEisner/test-results/middleware"
)

func TestLoggingMiddleware(t *testing.T) {
	type testCase struct {
		name       string
		method     string
		path       string
		statusCode int
		handler    http.HandlerFunc
	}

	tests := []testCase{
		{
			name:       "GET 200 OK",
			method:     http.MethodGet,
			path:       "/hello",
			statusCode: http.StatusOK,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "OK")
			},
		},
		{
			name:       "POST 201 Created",
			method:     http.MethodPost,
			path:       "/create",
			statusCode: http.StatusCreated,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			},
		},
		{
			name:       "GET 404 Not Found",
			method:     http.MethodGet,
			path:       "/not-found",
			statusCode: http.StatusNotFound,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo}))

			mw := middleware.LoggingMiddleware(logger)
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rr := httptest.NewRecorder()

			handler := mw(tc.handler)
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.statusCode {
				t.Errorf("expected status %d, got %d", tc.statusCode, rr.Code)
			}

			// Check log contains expected method and path
			logOutput := logBuf.String()
			if !strings.Contains(logOutput, tc.method) {
				t.Errorf("log missing method %q: %s", tc.method, logOutput)
			}
			if !strings.Contains(logOutput, tc.path) {
				t.Errorf("log missing path %q: %s", tc.path, logOutput)
			}
			if !strings.Contains(logOutput, http.StatusText(tc.statusCode)) &&
				!strings.Contains(logOutput, "status="+string(rune(tc.statusCode))) {
				t.Errorf("log missing status %d: %s", tc.statusCode, logOutput)
			}
		})
	}
}
