package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const envVars string = `
Configure a .env file in your home directory with the following variables:

If you not use a custom settings, create the file with the variables empty,
the program using the default values.

GIT_USERNAME=
GIT_TOKEN=
SSH_KEY_PATH=
SSH_KEY_NAME=
`

func LoadDotEnv() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get home directory")
	}

	dotEnvFile := homeDir + "/.env"

	if err := godotenv.Load(dotEnvFile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env file: %v\n", err)
		fmt.Printf("%s\n", envVars)
		os.Exit(125)
	}

	return nil
}
