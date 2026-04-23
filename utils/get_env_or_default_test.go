package utils

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key         string
		value       string
		defaultVal  string
		want        string
	}{
		{"env set", "TEST_KEY", "test_value", "default", "test_value"},
		{"env empty", "TEST_KEY", "", "default", "default"},
		{"env not set", "NONEXISTENT_KEY", "", "default", "default"},
		{"env set to empty string", "TEST_KEY_EMPTY", "", "fallback", "fallback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			got := GetEnvOrDefault(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("GetEnvOrDefault(%q, %q) = %v, want %v", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}