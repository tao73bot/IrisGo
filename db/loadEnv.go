package db

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// General environment variables.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
