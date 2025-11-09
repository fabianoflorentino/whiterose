package utils

import "os"

func YmlOrYamlExistsInHomeDir() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	ymlPath := homeDir + "/" + ".config.yml"
	yamlPath := homeDir + "/" + ".config.yaml"

	if _, err := os.Stat(ymlPath); err == nil {
		return true
	}

	if _, err := os.Stat(yamlPath); err == nil {
		return true
	}

	return false
}

func IsFileJSON(file string) bool {
	return len(file) >= 5 && file[len(file)-5:] == ".json"
}

func IsFileYAML(file string) bool {
	return (len(file) >= 5 && file[len(file)-5:] == ".yaml") || (len(file) >= 4 && file[len(file)-4:] == ".yml")
}
