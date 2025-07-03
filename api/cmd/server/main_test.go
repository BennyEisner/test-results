package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
	}{
		{
			name: "all environment variables set",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expectedConfig: &Config{
				DBHost:     "localhost",
				DBPort:     5432,
				DBUser:     "testuser",
				DBPassword: "testpass",
				DBName:     "testdb",
				ServerAddr: ":8080",
			},
		},
		{
			name: "invalid port defaults to 5432",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "invalid",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expectedConfig: &Config{
				DBHost:     "localhost",
				DBPort:     5432, // Default value
				DBUser:     "testuser",
				DBPassword: "testpass",
				DBName:     "testdb",
				ServerAddr: ":8080",
			},
		},
		{
			name:    "empty environment variables",
			envVars: map[string]string{},
			expectedConfig: &Config{
				DBHost:     "",
				DBPort:     5432, // Default value
				DBUser:     "",
				DBPassword: "",
				DBName:     "",
				ServerAddr: ":8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_NAME")

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Clean up after test
			defer func() {
				os.Unsetenv("DB_HOST")
				os.Unsetenv("DB_PORT")
				os.Unsetenv("DB_USER")
				os.Unsetenv("DB_PASSWORD")
				os.Unsetenv("DB_NAME")
			}()

			config := loadConfig()

			if config.DBHost != tt.expectedConfig.DBHost {
				t.Errorf("expected DBHost %s, got %s", tt.expectedConfig.DBHost, config.DBHost)
			}
			if config.DBPort != tt.expectedConfig.DBPort {
				t.Errorf("expected DBPort %d, got %d", tt.expectedConfig.DBPort, config.DBPort)
			}
			if config.DBUser != tt.expectedConfig.DBUser {
				t.Errorf("expected DBUser %s, got %s", tt.expectedConfig.DBUser, config.DBUser)
			}
			if config.DBPassword != tt.expectedConfig.DBPassword {
				t.Errorf("expected DBPassword %s, got %s", tt.expectedConfig.DBPassword, config.DBPassword)
			}
			if config.DBName != tt.expectedConfig.DBName {
				t.Errorf("expected DBName %s, got %s", tt.expectedConfig.DBName, config.DBName)
			}
			if config.ServerAddr != tt.expectedConfig.ServerAddr {
				t.Errorf("expected ServerAddr %s, got %s", tt.expectedConfig.ServerAddr, config.ServerAddr)
			}
		})
	}
}

func TestCreateServer(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	server := createServer(db)

	if server == nil {
		t.Fatal("expected server to be created, got nil")
	}

	// Test that the server can handle requests
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	// Should return 200 OK for health check
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Verify no unexpected database calls were made
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestRunServer(t *testing.T) {
	// Create a simple test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("test")); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	// Test with invalid address (should fail quickly)
	err := runServer("invalid-address", handler)
	if err == nil {
		t.Error("expected error with invalid address, got none")
	}
}
