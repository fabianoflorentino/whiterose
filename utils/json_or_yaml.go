package utils

import (
	"os"
	"path/filepath"
)

var (
	extensionsYAML = []string{".yml", ".yaml"}
)

// YmlOrYamlExistsInHomeDir checks if a .config.yml or .config.yaml file exists in the user's home directory.
func YmlOrYamlExistsInHomeDir() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	for _, e := range extensionsYAML {
		if _, err := os.Stat(homeDir + "/.config" + e); err == nil {
			return true
		}
	}

	return false
}

// IsFileJSON checks if the given file has a .json extension.
func IsFileJSON(file string) bool {
	return len(file) >= 5 && filepath.Ext(file) == ".json"
}

// IsFileYAML checks if the given file has a .yaml or .yml extension.
func IsFileYAML(file string) bool {
	ext := filepath.Ext(file)
	return ext == ".yaml" || ext == ".yml"
}
