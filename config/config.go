package config

import (
	"log" // log: package for logging
	"os"  // os: package for operating system functionality

	"github.com/joho/godotenv" // godotenv: package for reading .env files
)

func LoadEnv() { // LoadEnv: function to load environment variables
	if err := godotenv.Load(); err != nil { // godotenv.Load: function to load .env file
		log.Fatalf("Error loading .env file") // log.Fatalf: function to log fatal error
	}
}

func GetEnv(key, fallback string) string { // GetEnv: function to get environment variable
	if value, exists := os.LookupEnv(key); exists { // os.LookupEnv: function to get environment variable
		return value // return value
	}
	return fallback // return fallback
}
