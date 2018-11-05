package system

import "os"

// GetEnv - Get environment variable if it exists or else return fallback string
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
