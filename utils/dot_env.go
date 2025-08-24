package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadDotEnv() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get home directory")
	}

	dotEnvFile := homeDir + "/.env"

	if err := godotenv.Load(dotEnvFile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env file: %v\n", err)
		os.Exit(1)
	}

	return nil
}
