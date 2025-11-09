package utils

import (
	"os"
	"path/filepath"
)

var (
	extensionsYAML = []string{".yml", ".yaml"}
)

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

func IsFileJSON(file string) bool {
	return len(file) >= 5 && filepath.Ext(file) == ".json"
}

func IsFileYAML(file string) bool {
	ext := filepath.Ext(file)
	return ext == ".yaml" || ext == ".yml"
}
