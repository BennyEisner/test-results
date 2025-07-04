package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		expected    *Config
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"POSTGRES_DSN": "host=localhost port=5432 user=test password=test dbname=test",
				"PORT":         "9090",
				"ENVIRONMENT":  "production",
			},
			expectError: false,
			expected: &Config{
				Port:        "9090",
				PostgresDSN: "host=localhost port=5432 user=test password=test dbname=test",
				Environment: "production",
			},
		},
		{
			name: "missing POSTGRES_DSN",
			envVars: map[string]string{
				"PORT":        "9090",
				"ENVIRONMENT": "production",
			},
			expectError: true,
			expected:    nil,
		},
		{
			name: "default values",
			envVars: map[string]string{
				"POSTGRES_DSN": "host=localhost port=5432 user=test password=test dbname=test",
			},
			expectError: false,
			expected: &Config{
				Port:        "8080",
				PostgresDSN: "host=localhost port=5432 user=test password=test dbname=test",
				Environment: "development",
			},
		},
		{
			name: "empty POSTGRES_DSN",
			envVars: map[string]string{
				"POSTGRES_DSN": "",
				"PORT":         "9090",
				"ENVIRONMENT":  "production",
			},
			expectError: true,
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Unsetenv("POSTGRES_DSN")
			os.Unsetenv("PORT")
			os.Unsetenv("ENVIRONMENT")

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Clean up after test
			defer func() {
				os.Unsetenv("POSTGRES_DSN")
				os.Unsetenv("PORT")
				os.Unsetenv("ENVIRONMENT")
			}()

			config, err := Load()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if config.Port != tt.expected.Port {
				t.Errorf("expected Port %s, got %s", tt.expected.Port, config.Port)
			}

			if config.PostgresDSN != tt.expected.PostgresDSN {
				t.Errorf("expected PostgresDSN %s, got %s", tt.expected.PostgresDSN, config.PostgresDSN)
			}

			if config.Environment != tt.expected.Environment {
				t.Errorf("expected Environment %s, got %s", tt.expected.Environment, config.Environment)
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		fallback string
		envValue string
		expected string
	}{
		{
			name:     "environment variable set",
			key:      "TEST_KEY",
			fallback: "default",
			envValue: "custom",
			expected: "custom",
		},
		{
			name:     "environment variable not set",
			key:      "TEST_KEY",
			fallback: "default",
			envValue: "",
			expected: "default",
		},
		{
			name:     "environment variable empty",
			key:      "TEST_KEY",
			fallback: "default",
			envValue: "",
			expected: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variable
			os.Unsetenv(tt.key)

			// Set test environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			// Clean up after test
			defer os.Unsetenv(tt.key)

			result := getEnvOrDefault(tt.key, tt.fallback)

			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
