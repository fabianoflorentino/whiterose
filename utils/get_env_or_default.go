// GetEnvOrDefault returns the value of the environment variable specified by 'key'.
// If the environment variable is not set or is empty, it returns 'defaultValue' instead.
// This function is useful for providing fallback values when environment variables are missing.
package utils

import "os"

// GetEnvOrDefault returns the value of the environment variable specified by 'key'.
// If the environment variable is not set or is empty, it returns 'defaultValue' instead.
// This function is useful for providing fallback values when environment variables are missing.
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
