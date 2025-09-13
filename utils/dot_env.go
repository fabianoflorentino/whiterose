// LoadDotEnv loads environment variables from a .env file located in the user's home directory.
// It uses the github.com/joho/godotenv package to parse the file. If the .env file cannot be loaded,
// the function logs a fatal error message including a reference to the environment variables documentation
// and exits the program with status code 125. Returns an error if loading fails.
//
// See: https://github.com/fabianoflorentino/whiterose/blob/main/README.md#environment-variables
package utils

import (
	"fmt"
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
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	if err := godotenv.Load(homeDir + "/.env"); err != nil {
		log.Fatalf("failed to load .env file: %v, %s", err, envVars)
		os.Exit(125)
	}

	return nil
}
