package update

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fabianoflorentino/whiterose/internal/domain/entities"
)

func TestUpdateService_New(t *testing.T) {
	if New() == nil {
		t.Error("New() returned nil")
	}
}

func TestUpdateService_SetPRBase(t *testing.T) {
	s := New()
	s.SetPRBase("develop")
	if s.prBase != "develop" {
		t.Errorf("prBase = %v, want develop", s.prBase)
	}
}

func TestCalculateNewVersion(t *testing.T) {
	s := &UpdateService{}

	tests := []struct {
		name     string
		current  string
		strategy entities.UpdateStrategy
		major    bool
		want     string
	}{
		{"patch", "1.24.0", entities.StrategyPatch, false, "1.24.1"},
		{"minor", "1.24.0", entities.StrategyMinor, false, "1.25.0"},
		{"major", "1.24.0", entities.StrategyMajor, false, "2.0.0-rc.1"},
		{"major flag", "1.24.0", entities.StrategyPatch, true, "2.0.0-rc.1"},
		{"empty", "", entities.StrategyPatch, false, "1.25.0"},
		{"invalid", "1.2", entities.StrategyPatch, false, "1.25.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.calculateNewVersion(tt.current, tt.strategy, tt.major); got != tt.want {
				t.Errorf("calculateNewVersion(%q, %s, %v) = %v, want %v", tt.current, tt.strategy, tt.major, got, tt.want)
			}
		})
	}
}

func TestExtractDockerImage(t *testing.T) {
	s := &UpdateService{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "FROM golang:1.20", "golang:1.20"},
		{"alpine", "FROM golang:1.20-alpine", "golang:1.20-alpine"},
		{"no FROM", "RUN echo hello", ""},
		{"lowercase", "from golang:1.20", ""},
		{"multi", "RUN ls\nFROM alpine:latest", "alpine:latest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.extractDockerImage(tt.input); got != tt.want {
				t.Errorf("extractDockerImage(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCalculateNewDockerImage(t *testing.T) {
	s := &UpdateService{}

	tests := []struct {
		name     string
		current  string
		strategy entities.UpdateStrategy
		major    bool
		want     string
	}{
		{"invalid", "invalid", entities.StrategyPatch, false, "invalid"},
		{"minor", "golang:1.20", entities.StrategyMinor, false, "golang:1.21"},
		{"major", "golang:1.20", entities.StrategyMajor, false, "golang:2.0"},
		{"alpine", "golang:1.20-alpine", entities.StrategyMinor, false, "golang:1.21-alpine"},
		{"no version", "golang:latest", entities.StrategyMinor, false, "golang:latest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.calculateNewDockerImage(tt.current, tt.strategy, tt.major); got != tt.want {
				t.Errorf("calculateNewDockerImage(%q, %s, %v) = %v, want %v", tt.current, tt.strategy, tt.major, got, tt.want)
			}
		})
	}
}

func TestUpdateGoMod_FileNotFound(t *testing.T) {
	err := New().UpdateGoMod(entities.UpdateProject{Name: "test", Path: "/nonexistent"}, entities.StrategyPatch, false)
	if err == nil {
		t.Error("expected error")
	}
}

func TestUpdateGoVersion(t *testing.T) {
	err := New().UpdateGoVersion(entities.UpdateProject{Name: "test", Path: "/nonexistent"}, entities.StrategyPatch, false)
	if err == nil {
		t.Error("expected error")
	}
}

func TestUpdateDockerImage(t *testing.T) {
	err := New().UpdateDockerImage(entities.UpdateProject{Name: "test", Path: "/nonexistent"}, entities.StrategyPatch, false)
	if err == nil {
		t.Error("expected error")
	}
}

func TestCreateBranchAndCommit_NotFound(t *testing.T) {
	project := entities.UpdateProject{Name: "test", Path: "/nonexistent"}
	_, err := New().CreateBranchAndCommit(project, []string{})
	if err == nil {
		t.Error("expected error")
	}
}

func TestCreateBranchAndCommit_Success(t *testing.T) {
	t.Skip("requires git repo")
}

func TestCreatePR(t *testing.T) {
	s := New()
	s.SetPRBase("main")
	project := entities.UpdateProject{Name: "test", Path: "/tmp"}
	err := s.CreatePR(project, "branch", []string{})
	if err == nil {
		t.Error("expected error (gh cli not installed)")
	}
}

func TestCreatePRWithBase(t *testing.T) {
	s := New()
	project := entities.UpdateProject{Name: "test", Path: "/tmp"}
	err := s.CreatePRWithBase(project, "branch", []string{}, "develop")
	if err == nil {
		t.Error("expected error (gh cli not installed)")
	}
}

func TestLoadUpdateConfig_Valid(t *testing.T) {
	validConfig := `projects:
  - name: test
    path: ./test
`
	projects, err := entities.ParseUpdateConfig([]byte(validConfig))
	if err != nil {
		t.Errorf("ParseUpdateConfig() error = %v", err)
	}
	if len(projects) != 1 {
		t.Errorf("len = %d, want 1", len(projects))
	}

	if projects[0].Name != "test" {
		t.Errorf("Name = %v, want test", projects[0].Name)
	}
}

func TestLoadUpdateConfig_Empty(t *testing.T) {
	projects, err := entities.ParseUpdateConfig([]byte("projects: []"))
	if err != nil {
		t.Errorf("ParseUpdateConfig() error = %v", err)
	}
	if len(projects) != 0 {
		t.Errorf("len = %d, want 0", len(projects))
	}
}

func TestLoadUpdateConfig_Invalid(t *testing.T) {
	_, err := entities.ParseUpdateConfig([]byte("invalid: ["))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadUpdateConfig(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"with goMod", "projects:\n  - name: app\n    path: ./app\n    goMod:\n      updateStrategy: minor"},
		{"with dockerImage", "projects:\n  - name: app\n    path: ./app\n    dockerImage:\n      base: golang:1.20"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projects, err := entities.ParseUpdateConfig([]byte(tt.data))
			if err != nil {
				t.Errorf("ParseUpdateConfig() error = %v", err)
			}
			if len(projects) != 1 {
				t.Errorf("len = %d, want 1", len(projects))
			}
		})
	}
}

func TestGetCurrentGoVersion(t *testing.T) {
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\ngo 1.24.0"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	if got := New().getCurrentGoVersion(goModPath); got != "1.24.0" {
		t.Errorf("getCurrentGoVersion() = %v, want 1.24.0", got)
	}
}

func TestGetCurrentGoVersion_NotFound(t *testing.T) {
	if got := New().getCurrentGoVersion("/nonexistent"); got != "" {
		t.Errorf("getCurrentGoVersion() = %v, want empty", got)
	}
}

func TestUpdateGoModFile(t *testing.T) {
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\ngo 1.24.0"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	if err := New().updateGoModFile(goModPath, "1.25.0"); err != nil {
		t.Errorf("updateGoModFile() error = %v", err)
	}
}

func TestRunGoModTidy(t *testing.T) {
	err := New().runGoModTidy("/nonexistent")
	if err == nil {
		t.Error("expected error")
	}
}

func TestVersionChecker_New(t *testing.T) {
	if NewVersionChecker() == nil {
		t.Error("NewVersionChecker() returned nil")
	}
}

func TestVersionChecker_ExtractMajorVersion(t *testing.T) {
	vc := NewVersionChecker()

	tests := []struct {
		name string
		tag  string
		want string
	}{
		{"simple", "1.20", "1"},
		{"alpine", "1.20-alpine", "1"},
		{"rc", "1.21rc1", "1"},
		{"no number", "latest", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vc.extractMajorVersion(tt.tag); got != tt.want {
				t.Errorf("extractMajorVersion(%q) = %v, want %v", tt.tag, got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct{ a, b, want int }{{1, 2, 1}, {5, 3, 3}, {4, 4, 4}}
	for _, tt := range tests {
		if got := min(tt.a, tt.b); got != tt.want {
			t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestLoadUpdateConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "update.yaml")

	config := `projects:
  - name: test-project
    path: /tmp/test
    goMod:
      updateStrategy: patch
`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	s := New()
	projects, err := s.LoadUpdateConfig(configPath)
	if err != nil {
		t.Fatalf("LoadUpdateConfig() error = %v", err)
	}
	if len(projects) != 1 {
		t.Errorf("len(projects) = %d, want 1", len(projects))
	}
	if projects[0].Name != "test-project" {
		t.Errorf("Name = %v, want test-project", projects[0].Name)
	}
}

func TestLoadUpdateConfig_NotFound(t *testing.T) {
	s := New()
	_, err := s.LoadUpdateConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestCreateBranchAndCommit(t *testing.T) {
	t.Skip("Requires git repo setup")
}

func TestVersionChecker_GetCurrentGoVersion(t *testing.T) {
	vc := NewVersionChecker()
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")

	if err := os.WriteFile(goModPath, []byte("module test\ngo 1.26.0"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	got := vc.GetCurrentGoVersion(tmpDir)
	if got != "1.26.0" {
		t.Errorf("GetCurrentGoVersion() = %v, want 1.26.0", got)
	}
}

func TestVersionChecker_GetCurrentGoVersion_NotFound(t *testing.T) {
	vc := NewVersionChecker()
	got := vc.GetCurrentGoVersion("/nonexistent")
	if got != "" {
		t.Errorf("GetCurrentGoVersion() = %v, want empty", got)
	}
}

func TestVersionChecker_ListGoLibUpdates(t *testing.T) {
	t.Skip("Requires network access")
}

func TestVersionChecker_UpdatePackages(t *testing.T) {
	t.Skip("Requires network and git setup")
}

func TestVersionChecker_ListDockerUpdates(t *testing.T) {
	t.Skip("Requires network access")
}

func TestVersionChecker_ListGoVersions(t *testing.T) {
	t.Skip("Requires network access")
}
