package database

import (
	"log"
	"os"
	"strings"
)

// GetDatabaseName extracts the database name from the MONGO_URI environment variable
func GetDatabaseName() string {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set in the environment")
	}

	// Parse the database name from the URI
	parts := strings.Split(mongoURI, "/")
	if len(parts) < 4 {
		log.Fatal("Invalid MONGO_URI format. Unable to parse database name.")
	}
	return strings.Split(parts[3], "?")[0] // Extract the database name before query parameters
}
