package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// Helper to create a temporary JSON config file for testing.
func createTempJSONConfig(t *testing.T, repos []RepoInfo) string {
	t.Helper()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.json")

	cfg := ConfigFile{
		Repositories: repos,
	}

	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(cfg); err != nil {
		t.Fatalf("failed to encode JSON: %v", err)
	}

	return filePath
}

func createTempYAMLConfig(t *testing.T, repos []RepoInfo) string {
	t.Helper()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.yaml")

	cfg := ConfigFile{
		Repositories: repos,
	}

	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal YAML: %v", err)
	}

	if _, err := f.Write(data); err != nil {
		t.Fatalf("failed to write YAML to file: %v", err)
	}

	return filePath
}

func createTempYMLConfig(t *testing.T, repos []RepoInfo) string {
	t.Helper()
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.yml")

	cfg := ConfigFile{
		Repositories: repos,
	}

	f, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()

	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal YAML: %v", err)
	}

	if _, err := f.Write(data); err != nil {
		t.Fatalf("failed to write YAML to file: %v", err)
	}

	return filePath
}

func TestFetchRepositories_JSONSuccess(t *testing.T) {
	repos := []RepoInfo{
		{URL: "https://github.com/example/repo1", Directory: "repo1"},
		{URL: "https://github.com/example/repo2", Directory: "repo2"},
	}
	file := createTempJSONConfig(t, repos)

	got, err := FetchRepositories(file)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != len(repos) {
		t.Fatalf("expected %d repos, got %d", len(repos), len(got))
	}
	for i, repo := range repos {
		if got[i] != repo {
			t.Errorf("expected repo %v, got %v", repo, got[i])
		}
	}
}

func TestFetchRepositories_YAMLSuccess(t *testing.T) {
	repos := []RepoInfo{
		{URL: "https://github.com/example/repo1", Directory: "repo1"},
		{URL: "https://github.com/example/repo2", Directory: "repo2"},
	}
	file := createTempYAMLConfig(t, repos)

	got, err := FetchRepositories(file)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != len(repos) {
		t.Fatalf("expected %d repos, got %d", len(repos), len(got))
	}
	for i, repo := range repos {
		if got[i] != repo {
			t.Errorf("expected repo %v, got %v", repo, got[i])
		}
	}
}

func TestFetchRepositories_YMLSuccess(t *testing.T) {
	repos := []RepoInfo{
		{URL: "https://github.com/example/repo1", Directory: "repo1"},
		{URL: "https://github.com/example/repo2", Directory: "repo2"},
	}
	file := createTempYMLConfig(t, repos)

	got, err := FetchRepositories(file)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != len(repos) {
		t.Fatalf("expected %d repos, got %d", len(repos), len(got))
	}
	for i, repo := range repos {
		if got[i] != repo {
			t.Errorf("expected repo %v, got %v", repo, got[i])
		}
	}
}

func TestFetchRepositories_FileNotFound(t *testing.T) {
	_, err := FetchRepositories("nonexistent.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}

func TestFetchRepositories_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(filePath, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("failed to write invalid json: %v", err)
	}

	_, err := FetchRepositories(filePath)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestFetchRepositories_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.yaml")
	if err := os.WriteFile(filePath, []byte("invalid: [yaml"), 0644); err != nil {
		t.Fatalf("failed to write invalid yaml: %v", err)
	}

	_, err := FetchRepositories(filePath)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestFetchRepositories_InvalidYML(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.yml")
	if err := os.WriteFile(filePath, []byte("invalid: [yml"), 0644); err != nil {
		t.Fatalf("failed to write invalid yml: %v", err)
	}

	_, err := FetchRepositories(filePath)

	if err == nil {
		t.Fatal("expected error for invalid YML, got nil")
	}
}

func TestFetchAppsInfo(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	
	config := `{"applications":[{"name":"Go","command":"go","versionFlag":"version","recommendedVersion":"1.20"}]}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	
	apps, err := FetchAppsInfo(configPath)
	if err != nil {
		t.Errorf("FetchAppsInfo() error = %v", err)
	}
	if len(apps) != 1 {
		t.Errorf("len(apps) = %d, want 1", len(apps))
	}
	if apps[0].Name != "Go" {
		t.Errorf("Name = %v, want Go", apps[0].Name)
	}
}

func TestFetchAppsInfo_WithYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	
	config := `applications:
  - name: Go
    command: go
    versionFlag: version`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	
	apps, err := FetchAppsInfo(configPath)
	if err != nil {
		t.Errorf("FetchAppsInfo() error = %v", err)
	}
	if len(apps) != 1 {
		t.Errorf("len(apps) = %d, want 1", len(apps))
	}
}
