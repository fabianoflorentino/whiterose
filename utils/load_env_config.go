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
	"strings"

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

// LoadDotConfigJSON loads a .config.json file from the user's home directory.
// It returns the file path as a string and an error if the file cannot be opened.
func LoadDotConfig() (string, error) {
	var jsonFile string = strings.TrimSpace(".config.json")
	var yamlFile string = strings.TrimSpace(".config.yaml")
	var ymlFile string = strings.TrimSpace(".config.yml")

	if YmlOrYamlExistsInHomeDir() {
		if fileExistsInHomeDir(ymlFile) {
			return getFilePathInHomeDir(ymlFile)
		} else if fileExistsInHomeDir(yamlFile) {
			return getFilePathInHomeDir(yamlFile)
		}
	} else if fileExistsInHomeDir(jsonFile) {
		return getFilePathInHomeDir(jsonFile)
	}

	return "", fmt.Errorf("no configuration file found in home directory")
}

// fileExistsInHomeDir checks if a file exists in the user's home directory.
func fileExistsInHomeDir(file string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	_, err = os.Stat(homeDir + "/" + file)
	return !os.IsNotExist(err)
}

// getFilePathInHomeDir returns the full path of a file in the user's home directory.
func getFilePathInHomeDir(file string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	f, err := os.Open(homeDir + "/" + file)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing %s: %v\n", file, err)
		}
	}()

	return homeDir + "/" + file, nil
}
