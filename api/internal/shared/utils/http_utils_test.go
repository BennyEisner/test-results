package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		data       interface{}
		expectJSON string
	}{
		{
			name:   "simple object",
			status: http.StatusOK,
			data: map[string]string{
				"message": "success",
			},
			expectJSON: `{"message":"success"}`,
		},
		{
			name:       "array",
			status:     http.StatusOK,
			data:       []string{"item1", "item2"},
			expectJSON: `["item1","item2"]`,
		},
		{
			name:   "custom status",
			status: http.StatusCreated,
			data: map[string]int{
				"id": 123,
			},
			expectJSON: `{"id":123}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			RespondWithJSON(rr, tt.status, tt.data)

			if rr.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, rr.Code)
			}

			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			// Parse the response body to verify JSON
			var responseData interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &responseData); err != nil {
				t.Errorf("failed to parse JSON response: %v", err)
			}

			// Compare JSON strings (normalized)
			expectedBytes, _ := json.Marshal(tt.data)
			expected := string(expectedBytes)
			actual := strings.TrimSpace(rr.Body.String())

			if expected != actual {
				t.Errorf("expected JSON %q, got %q", expected, actual)
			}
		})
	}
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		message    string
		expectJSON string
	}{
		{
			name:       "bad request",
			status:     http.StatusBadRequest,
			message:    "invalid input",
			expectJSON: `{"error":"invalid input"}`,
		},
		{
			name:       "not found",
			status:     http.StatusNotFound,
			message:    "resource not found",
			expectJSON: `{"error":"resource not found"}`,
		},
		{
			name:       "internal server error",
			status:     http.StatusInternalServerError,
			message:    "something went wrong",
			expectJSON: `{"error":"something went wrong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			RespondWithError(rr, tt.status, tt.message)

			if rr.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, rr.Code)
			}

			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			if strings.TrimSpace(rr.Body.String()) != tt.expectJSON {
				t.Errorf("expected JSON %q, got %q", tt.expectJSON, strings.TrimSpace(rr.Body.String()))
			}
		})
	}
}

func TestNullStringToStringPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected *string
	}{
		{
			name: "valid string",
			input: sql.NullString{
				String: "test",
				Valid:  true,
			},
			expected: stringPtr("test"),
		},
		{
			name: "invalid string",
			input: sql.NullString{
				String: "",
				Valid:  false,
			},
			expected: nil,
		},
		{
			name: "empty valid string",
			input: sql.NullString{
				String: "",
				Valid:  true,
			},
			expected: stringPtr(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NullStringToStringPtr(tt.input)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", *result)
				}
			} else {
				if result == nil {
					t.Errorf("expected %s, got nil", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("expected %s, got %s", *tt.expected, *result)
				}
			}
		})
	}
}

func TestNullInt64ToIntPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullInt64
		expected *int64
	}{
		{
			name: "valid int64",
			input: sql.NullInt64{
				Int64: 123,
				Valid: true,
			},
			expected: int64Ptr(123),
		},
		{
			name: "invalid int64",
			input: sql.NullInt64{
				Int64: 0,
				Valid: false,
			},
			expected: nil,
		},
		{
			name: "zero valid int64",
			input: sql.NullInt64{
				Int64: 0,
				Valid: true,
			},
			expected: int64Ptr(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NullInt64ToIntPtr(tt.input)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", *result)
				}
			} else {
				if result == nil {
					t.Errorf("expected %d, got nil", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("expected %d, got %d", *tt.expected, *result)
				}
			}
		})
	}
}

func TestNullFloat64ToFloat64Ptr(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullFloat64
		expected *float64
	}{
		{
			name: "valid float64",
			input: sql.NullFloat64{
				Float64: 123.45,
				Valid:   true,
			},
			expected: float64Ptr(123.45),
		},
		{
			name: "invalid float64",
			input: sql.NullFloat64{
				Float64: 0.0,
				Valid:   false,
			},
			expected: nil,
		},
		{
			name: "zero valid float64",
			input: sql.NullFloat64{
				Float64: 0.0,
				Valid:   true,
			},
			expected: float64Ptr(0.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NullFloat64ToFloat64Ptr(tt.input)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", *result)
				}
			} else {
				if result == nil {
					t.Errorf("expected %f, got nil", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("expected %f, got %f", *tt.expected, *result)
				}
			}
		})
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
