package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGitRepository(t *testing.T) {
	g := NewGitRepository()
	if g == nil {
		t.Error("NewGitRepository() returned nil")
	}
}

func TestLoadRepositoriesFromFile_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	config := `{"repositories":[{"url":"https://github.com/test/repo","directory":"test"}]}`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	repos, err := LoadRepositoriesFromFile(configPath)
	if err != nil {
		t.Errorf("LoadRepositoriesFromFile() error = %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("len(repos) = %d, want 1", len(repos))
	}
}

func TestLoadRepositoriesFromFile_YAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	config := `repositories:
  - url: https://github.com/test/repo
    directory: test`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	repos, err := LoadRepositoriesFromFile(configPath)
	if err != nil {
		t.Errorf("LoadRepositoriesFromFile() error = %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("len(repos) = %d, want 1", len(repos))
	}
}

func TestLoadRepositoriesFromFile_NotFound(t *testing.T) {
	_, err := LoadRepositoriesFromFile("/nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoadRepositoriesFromFile_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(configPath, []byte("{invalid"), 0644); err != nil {
		t.Fatalf("failed to create invalid config: %v", err)
	}

	_, err := LoadRepositoriesFromFile(configPath)
	if err == nil {
		t.Error("expected error for invalid file")
	}
}

func TestGitCloneOptions_Setup(t *testing.T) {
	t.Skip("Setup requires config file")
}

func TestGitCloneOptions_fetchRepositories(t *testing.T) {
	t.Skip("fetchRepositories clones real repos")
}
