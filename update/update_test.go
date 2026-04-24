package update

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fabianoflorentino/whiterose/internal/domain/entities"
	"github.com/fabianoflorentino/whiterose/mocks"
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

func TestVersionChecker_WithHTTP(t *testing.T) {
	vc := NewVersionChecker()
	mock := &mocks.MockHTTPClient{}
	vc2 := vc.WithHTTP(mock)
	if vc2 == nil {
		t.Error("WithHTTP returned nil")
	}
}

func TestVersionChecker_WithExecutor(t *testing.T) {
	vc := NewVersionChecker()
	mock := &mocks.MockCommandExecutor{}
	vc2 := vc.WithExecutor(mock)
	if vc2 == nil {
		t.Error("WithExecutor returned nil")
	}
}

func TestUpdateService_GetCurrentGoVersion(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n\ngo 1.24.0"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	s := &UpdateService{}
	version := s.getCurrentGoVersion(filepath.Join(tmpDir, "go.mod"))
	if version != "1.24.0" {
		t.Errorf("getCurrentGoVersion() = %v, want 1.24.0", version)
	}
}

func TestUpdateService_GetCurrentGoVersion_Empty(t *testing.T) {
	s := &UpdateService{}
	version := s.getCurrentGoVersion("/nonexistent/go.mod")
	if version != "" {
		t.Errorf("getCurrentGoVersion() = %v, want empty", version)
	}
}

func TestVersionChecker_UpdatePackages_Execute(t *testing.T) {
	calls := 0
	vc := NewVersionChecker().WithExecutor(&mocks.MockCommandExecutor{
		RunFunc: func(cmd string, args ...string) (string, error) {
			calls++
			return "go get: upgrading example.com/pkg v1.0.0 => v1.1.0", nil
		},
	})

	err := vc.UpdatePackages("/tmp", "minor", false)
	if err != nil {
		t.Errorf("UpdatePackages() error = %v", err)
	}
	if calls < 2 {
		t.Errorf("expected at least 2 calls (get + tidy), got %d", calls)
	}
}

func TestVersionChecker_UpdatePackages_GetError(t *testing.T) {
	vc := NewVersionChecker().WithExecutor(&mocks.MockCommandExecutor{
		RunFunc: func(cmd string, args ...string) (string, error) {
			return "", errors.New("go get failed")
		},
	})

	err := vc.UpdatePackages("/tmp", "patch", false)
	if err == nil {
		t.Error("expected error from UpdatePackages")
	}
}

func TestVersionChecker_fetchDockerHubTags_InvalidJSON(t *testing.T) {
	vc := NewVersionChecker().WithHTTP(&mocks.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("not json")),
			}, nil
		},
	})

	tags, err := vc.fetchDockerHubTags("golang")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if tags != nil {
		t.Errorf("tags = %v, want nil", tags)
	}
}

func TestVersionChecker_fetchGHCRTags_Success(t *testing.T) {
	vc := NewVersionChecker().WithHTTP(&mocks.MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"tags": ["1.0", "1.1"]}`)),
			}, nil
		},
	})

	tags, err := vc.fetchGHCRTags("test/image")
	if err != nil {
		t.Errorf("fetchGHCRTags() error = %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("len(tags) = %d, want 2", len(tags))
	}
}

func TestUpdateService_updateGoModFile(t *testing.T) {
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")

	if err := os.WriteFile(goModPath, []byte("module test\n\ngo 1.24.0"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	s := &UpdateService{}
	err := s.updateGoModFile(goModPath, "1.25.0")
	if err != nil {
		t.Errorf("updateGoModFile() error = %v", err)
	}

	content, _ := os.ReadFile(goModPath)
	if !strings.Contains(string(content), "go 1.25.0") {
		t.Error("go.mod was not updated correctly")
	}
}

func TestUpdateService_updateGoModFile_MultipleGoLines(t *testing.T) {
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")

	content := `module test

go 1.24.0

require (
	example.com/pkg v1.0.0
)
`
	if err := os.WriteFile(goModPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	s := &UpdateService{}
	err := s.updateGoModFile(goModPath, "1.26.0")
	if err != nil {
		t.Errorf("updateGoModFile() error = %v", err)
	}

	data, _ := os.ReadFile(goModPath)
	lines := strings.Split(string(data), "\n")
	goLines := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "go ") {
			goLines++
		}
	}
	if goLines != 1 {
		t.Errorf("expected exactly one go line, got %d", goLines)
	}
}

func TestUpdateService_LoadUpdateConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	config := `projects:
  - name: test-project
    path: ./test
    type: go-mod
`
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	s := New()
	projects, err := s.LoadUpdateConfig(configPath)
	if err != nil {
		t.Errorf("LoadUpdateConfig() error = %v", err)
	}
	if len(projects) != 1 {
		t.Errorf("len(projects) = %d, want 1", len(projects))
	}
}

func TestVersionChecker_UpdatePackages_NoUpdates(t *testing.T) {
	vc := NewVersionChecker().WithExecutor(&mocks.MockCommandExecutor{
		RunFunc: func(cmd string, args ...string) (string, error) {
			return "", nil
		},
	})

	err := vc.UpdatePackages("/tmp", "patch", false)
	if err != nil {
		t.Errorf("UpdatePackages() error = %v", err)
	}
}

func TestUpdateService_UpdateGoMod_Success(t *testing.T) {
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")

	if err := os.WriteFile(goModPath, []byte("module test\n\ngo 1.24.0\n\nrequire ()\n"), 0644); err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	s := &UpdateService{
		executor: &mocks.MockCommandExecutor{
			RunFunc: func(cmd string, args ...string) (string, error) {
				return "", nil
			},
		},
	}

	err := s.UpdateGoMod(entities.UpdateProject{Name: "test", Path: tmpDir}, entities.StrategyMinor, false)
	if err != nil {
		t.Errorf("UpdateGoMod() error = %v", err)
	}

	content, _ := os.ReadFile(goModPath)
	if !strings.Contains(string(content), "go 1.25.0") {
		t.Errorf("go.mod not updated: %s", string(content))
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
	s := &UpdateService{
		executor: &mocks.MockCommandExecutor{
			RunFunc: func(cmd string, args ...string) (string, error) {
				return "", errors.New("go mod tidy failed")
			},
		},
	}
	err := s.runGoModTidy("/tmp")
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

func TestVersionChecker_ListGoLibUpdates_Success(t *testing.T) {
	vc := NewVersionChecker().WithExecutor(&mocks.MockCommandExecutor{
		RunFunc: func(cmd string, args ...string) (string, error) {
			if cmd == "go" && args[0] == "list" {
				return `github.com/foo v1.0.0
github.com/bar v2.0.0 [v2.1.0]`, nil
			}
			return "", nil
		},
	})

	err := vc.ListGoLibUpdates("/tmp")
	if err != nil {
		t.Errorf("ListGoLibUpdates() error = %v", err)
	}
}

func TestVersionChecker_ListGoLibUpdates_Error(t *testing.T) {
	vc := NewVersionChecker().WithExecutor(&mocks.MockCommandExecutor{
		RunFunc: func(cmd string, args ...string) (string, error) {
			if cmd == "go" && args[0] == "list" {
				return "", errors.New("go not installed")
			}
			return "", nil
		},
	})

	err := vc.ListGoLibUpdates("/tmp")
	if err == nil {
		t.Error("expected error from ListGoLibUpdates")
	}
}

func TestVersionChecker_ListGoLibUpdates_NoUpdates(t *testing.T) {
	vc := NewVersionChecker().WithExecutor(&mocks.MockCommandExecutor{
		RunFunc: func(cmd string, args ...string) (string, error) {
			if cmd == "go" && args[0] == "list" {
				return `github.com/foo v1.0.0
github.com/bar v2.0.0`, nil
			}
			return "", nil
		},
	})

	err := vc.ListGoLibUpdates("/tmp")
	if err != nil {
		t.Errorf("ListGoLibUpdates() error = %v", err)
	}
}

func TestVersionChecker_UpdatePackages_DryRun(t *testing.T) {
	vc := NewVersionChecker().WithExecutor(&mocks.MockCommandExecutor{
		RunFunc: func(cmd string, args ...string) (string, error) {
			if cmd == "go" && args[0] == "get" {
				return `go get: upgrading golang.org/x/net v0.21.0 => v0.22.0`, nil
			}
			return "", nil
		},
	})

	err := vc.UpdatePackages("/tmp", "major", true)
	if err != nil {
		t.Errorf("UpdatePackages() error = %v", err)
	}
}

func TestVersionChecker_UpdatePackages_DryRunError(t *testing.T) {
	vc := NewVersionChecker().WithExecutor(&mocks.MockCommandExecutor{
		RunFunc: func(cmd string, args ...string) (string, error) {
			if cmd == "go" && args[0] == "get" {
				return "", errors.New("dry-run failed")
			}
			return "", nil
		},
	})

	err := vc.UpdatePackages("/tmp", "major", true)
	if err == nil {
		t.Error("expected error from UpdatePackages with dry-run")
	}
}

func TestVersionChecker_ListDockerUpdates_InvalidFormat(t *testing.T) {
	vc := NewVersionChecker()
	err := vc.ListDockerUpdates("invalid-image")
	if err == nil {
		t.Error("expected error for invalid image format")
	}
}

func TestVersionChecker_ListDockerUpdates_Success(t *testing.T) {
	vc := NewVersionChecker().WithHTTP(&mocks.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"results": [{"name": "1.25.0"}, {"name": "1.25.1"}]}`)),
			}, nil
		},
	})

	err := vc.ListDockerUpdates("golang:1.25")
	if err != nil {
		t.Errorf("ListDockerUpdates() error = %v", err)
	}
}

func TestVersionChecker_ListDockerUpdates_NotFound(t *testing.T) {
	vc := NewVersionChecker().WithHTTP(&mocks.MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"tags": ["1.25.0"]}`)),
			}, nil
		},
	})

	err := vc.ListDockerUpdates("golang:1.25")
	if err != nil {
		t.Errorf("ListDockerUpdates() error = %v", err)
	}
}

func TestVersionChecker_ListGoVersions_Success(t *testing.T) {
	vc := NewVersionChecker().WithHTTP(&mocks.MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(`[
					{"version": "go1.25.0", "stable": true},
					{"version": "go1.26.0beta1", "stable": false}
				]`)),
			}, nil
		},
	})

	err := vc.ListGoVersions()
	if err != nil {
		t.Errorf("ListGoVersions() error = %v", err)
	}
}

func TestVersionChecker_ListGoVersions_HTTPError(t *testing.T) {
	vc := NewVersionChecker().WithHTTP(&mocks.MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	})

	err := vc.ListGoVersions()
	if err == nil {
		t.Error("expected error from ListGoVersions")
	}
}

func TestVersionChecker_ListGoVersions_InvalidJSON(t *testing.T) {
	vc := NewVersionChecker().WithHTTP(&mocks.MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("invalid json")),
			}, nil
		},
	})

	err := vc.ListGoVersions()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
