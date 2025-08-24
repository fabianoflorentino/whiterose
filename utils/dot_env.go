package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const envVars string = `
https://github.com/fabianoflorentino/whiterose/blob/main/README.md#environment-variables
`

func LoadDotEnv() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get home directory")
	}

	dotEnvFile := homeDir + "/.env"

	if err := godotenv.Load(dotEnvFile); err != nil {
		log.Fatalf("Failed to load .env file: %v, %s", err, envVars)
		os.Exit(125)
	}

	return nil
}
