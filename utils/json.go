package utils

import (
	"encoding/json"
	"log"
	"os"
)

const (
	repoFile string = `
https://github.com/fabianoflorentino/whiterose/blob/main/README.md#usage
`
)

type RepoInfo struct {
	URL       string `json:"url"`
	Directory string `json:"directory"`
}

type ConfigFile struct {
	Repositories []RepoInfo `json:"repositories"`
}

func FetchReposFromJSON(file string) ([]RepoInfo, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open file: %v, %s", err, repoFile)
		return nil, err
	}
	defer f.Close()

	var rf ConfigFile
	if err := json.NewDecoder(f).Decode(&rf); err != nil {
		log.Fatalf("failed to decode JSON: %v", err)
		return nil, err
	}

	return rf.Repositories, nil
}
