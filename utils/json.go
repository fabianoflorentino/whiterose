// Package utils provides utility functions for handling JSON configuration files.
//
// This file defines structures and functions for reading repository information from a JSON file.
//
// Constants:
//   - repoFile: Contains a reference URL for usage instructions.
//
// Types:
//   - RepoInfo: Represents a repository with its URL and local directory.
//   - ConfigFile: Represents the configuration file structure containing a list of repositories.
//
// Functions:
//   - FetchReposFromJSON(file string) ([]RepoInfo, error):
//     Reads a JSON file specified by 'file', decodes its contents into a ConfigFile struct,
//     and returns the list of repositories. Logs fatal errors if the file cannot be opened or decoded.
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

// FetchReposFromJSON reads a JSON file specified by 'file', decodes its contents into a ConfigFile struct,
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
