package system

import (
	"os"

	"github.com/satori/go.uuid"
)

// GetEnv - Get environment variable if it exists or else return fallback string
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// CreateUUID - Create UUID
func CreateUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}
