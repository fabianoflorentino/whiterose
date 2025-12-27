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
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	repoFile string = `
https://github.com/fabianoflorentino/whiterose/blob/main/README.md#usage
`
)

type RepoInfo struct {
	URL       string `json:"url" yaml:"url"`
	Directory string `json:"directory" yaml:"directory"`
}

type AppInfo struct {
	Name                string            `json:"name" yaml:"name"`
	Command             string            `json:"command" yaml:"command"`
	VersionFlag         string            `json:"versionFlag" yaml:"versionFlag"`
	RecommendedVersion  string            `json:"recommendedVersion" yaml:"recommendedVersion"`
	InstallInstructions map[string]string `json:"installInstructions" yaml:"installInstructions"`
}

type ConfigFile struct {
	Repositories []RepoInfo `json:"repositories" yaml:"repositories"`
	Applications []AppInfo  `json:"applications" yaml:"applications"`
}

// FetchRepositories reads a JSON file specified by 'file', decodes its contents into a ConfigFile struct,
// and returns the list of repositories.
func FetchRepositories(file string) ([]RepoInfo, error) {
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

	if err := configDecode(file, &rf); err != nil {
		return nil, err
	}

	return rf.Repositories, nil
}

// FetchAppsInfo reads a JSON or YAML file specified by 'file', decodes its contents into a ConfigFile struct,
// and returns the list of applications.
func FetchAppsInfo(file string) ([]AppInfo, error) {
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

	if err := configDecode(file, &cfg); err != nil {
		return nil, err
	}

	return cfg.Applications, nil
}

// configDecode decodes the configuration file (JSON or YAML) into the provided ConfigFile struct.
func configDecode(file string, cfg *ConfigFile) error {
	fileHandle, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v, %s", err, repoFile)
	}
	defer func() {
		if err := fileHandle.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing %s: %v\n", file, err)
		}
	}()

	var f io.Reader = fileHandle

	if IsFileJSON(file) {
		if err := json.NewDecoder(f).Decode(&cfg); err != nil {
			return fmt.Errorf("failed to decode JSON: %v", err)
		}
	}

	if IsFileYAML(file) {
		if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
			return fmt.Errorf("failed to decode YAML: %v", err)
		}
	}

	return nil
}
