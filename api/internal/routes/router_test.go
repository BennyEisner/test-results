package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewRouter(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	// Set up mock expectations for any database queries that might be called
	// during router initialization
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	// Create router
	router := NewRouter(db)

	if router == nil {
		t.Fatal("expected router to be created, got nil")
	}

	// Test that the router is not nil
	if router == nil {
		t.Error("expected router to be created")
	}
}

func TestHealthEndpoints(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	// Set up mock expectations
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	// Create router
	router := NewRouter(db)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "health check",
			method:         "GET",
			path:           "/healthz",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK\n",
		},
		{
			name:           "readiness check",
			method:         "GET",
			path:           "/readyz",
			expectedStatus: http.StatusOK,
			expectedBody:   "Ready\n",
		},
		{
			name:           "metrics endpoint",
			method:         "GET",
			path:           "/metrics",
			expectedStatus: http.StatusOK,
			expectedBody:   "# HELP dummy_metric A dummy metric\n# TYPE dummy_metric counter\ndummy_metric 1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestHelloEndpoint(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	// Set up mock expectations
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	// Create router
	router := NewRouter(db)

	req := httptest.NewRequest("GET", "/api/hello", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()
	if body == "" {
		t.Error("expected non-empty response body")
	}

	// The response should contain "Hello from" and a timestamp
	if len(body) < 10 {
		t.Errorf("expected longer response body, got %q", body)
	}
}

func TestRouterWithNilDB(t *testing.T) {
	// Test that NewRouter handles nil database gracefully
	router := NewRouter(nil)

	if router == nil {
		t.Fatal("expected router to be created even with nil database")
	}

	// Test a simple endpoint that doesn't require database access
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Should still return OK for health check
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestRouterMiddleware(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	// Set up mock expectations
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	// Create router
	router := NewRouter(db)

	// Test that CORS middleware is applied
	req := httptest.NewRequest("OPTIONS", "/api/hello", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Check for CORS headers
	corsHeader := rr.Header().Get("Access-Control-Allow-Origin")
	if corsHeader == "" {
		t.Error("expected CORS headers to be present")
	}
}
