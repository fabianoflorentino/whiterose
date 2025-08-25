// LoadDotEnv loads environment variables from a .env file located in the user's home directory.
// It uses the github.com/joho/godotenv package to parse the file. If the .env file cannot be loaded,
// the function logs a fatal error message including a reference to the environment variables documentation
// and exits the program with status code 125. Returns an error if loading fails.
//
// See: https://github.com/fabianoflorentino/whiterose/blob/main/README.md#environment-variables
package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const envVars string = `
https://github.com/fabianoflorentino/whiterose/blob/main/README.md#environment-variables
`

// LoadDotEnv loads environment variables from a .env file located in the user's home directory.
// It uses the github.com/joho/godotenv package to parse the file. If the .env file cannot be loaded,
// the function logs a fatal error message including a reference to the environment variables documentation
// and exits the program with status code 125. Returns an error if loading fails.
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
