package utils

func IsFileJSON(file string) bool {
	return len(file) >= 5 && file[len(file)-5:] == ".json"
}

func IsFileYAML(file string) bool {
	return (len(file) >= 5 && file[len(file)-5:] == ".yaml") || (len(file) >= 4 && file[len(file)-4:] == ".yml")
}
