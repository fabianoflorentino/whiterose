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
	"fmt"
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

type AppInfo struct {
	Name                string
	Command             string
	VersionFlag         string
	RecommendedVersion  string
	InstallInstructions map[string]string
}

type ConfigFile struct {
	Repositories []RepoInfo `json:"repositories"`
	Applications []AppInfo  `json:"applications"`
}

// FetchReposFromJSON reads a JSON file specified by 'file', decodes its contents into a ConfigFile struct,
func FetchReposFromJSON(file string) ([]RepoInfo, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v, %s", err, repoFile)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing %s: %v\n", file, err)
		}
	}()

	var rf ConfigFile
	if err := json.NewDecoder(f).Decode(&rf); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return rf.Repositories, nil
}

func FetchAppsInfoFromJSON(file string) ([]AppInfo, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v, %s", err, repoFile)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing %s: %v\n", file, err)
		}
	}()

	var cfg ConfigFile
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return cfg.Applications, nil
}
