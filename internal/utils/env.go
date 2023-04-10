package utils

import "os"

// EnvOrDefault returns the value of the environment variable named by the key.
// If the variable is not present, it returns the default value.
func EnvOrDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
