package services

import (
	"fmt"
	"testing"

	"github.com/fabianoflorentino/whiterose/internal/interfaces"
)

func TestConfigService_New(t *testing.T) {
	svc := NewConfigService("")
	if svc == nil {
		t.Error("NewConfigService() returned nil")
	}
}

func TestConfigService_New_WithPath(t *testing.T) {
	svc := NewConfigService("/fake/path")
	if svc.configPath != "/fake/path" {
		t.Errorf("configPath = %v, want /fake/path", svc.configPath)
	}
}

func TestConfigService_LoadRepositories_NotFound(t *testing.T) {
	svc := NewConfigService("/nonexistent")
	_, err := svc.LoadRepositories()
	if err == nil {
		t.Error("expected error for nonexistent config")
	}
}

func TestConfigService_LoadApps(t *testing.T) {
	svc := NewConfigService("")
	apps, err := svc.LoadApps()
	if err != nil {
		t.Errorf("LoadApps() error = %v", err)
	}
	if len(apps) == 0 {
		t.Error("expected default apps")
	}
}

func TestExecutorService_New(t *testing.T) {
	svc := NewExecutorService()
	if svc == nil {
		t.Error("NewExecutorService() returned nil")
	}
}

func TestExecutorService_Run(t *testing.T) {
	svc := NewExecutorService()
	out, err := svc.Run("echo", "hello")
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}
	if out != "hello" {
		t.Errorf("out = %v, want hello", out)
	}
}

func TestExecutorService_Which_NotFound(t *testing.T) {
	svc := NewExecutorService()
	_, err := svc.Which("nonexistent-command-xyz")
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
}

func TestAppValidatorService_New(t *testing.T) {
	svc := NewAppValidatorService()
	if svc == nil {
		t.Error("NewAppValidatorService() returned nil")
	}
}

func TestAppValidatorService_WithExecutor(t *testing.T) {
	svc := NewAppValidatorService()
	mockExec := &mockExecutor{}
	result := svc.WithExecutor(mockExec)
	if result != svc {
		t.Error("WithExecutor should return same service")
	}
}

func TestAppValidatorService_ValidateOne_NotInstalled(t *testing.T) {
	svc := NewAppValidatorService()
	svc.WithExecutor(&mockExecutor{found: false, err: fmt.Errorf("not found")})
	result := svc.ValidateOne(interfaces.AppInfo{
		Name:        "fake-app",
		Command:     "fake-cmd",
		VersionFlag: "--version",
	})
	if result.IsInstalled {
		t.Error("expected installed = false")
	}
}

func TestAppValidatorService_ValidateOne_Installed(t *testing.T) {
	svc := NewAppValidatorService()
	svc.WithExecutor(&mockExecutor{found: true, version: "1.0.0"})
	result := svc.ValidateOne(interfaces.AppInfo{
		Name:                 "fake-app",
		Command:             "fake-cmd",
		VersionFlag:          "--version",
		RecommendedVersion:  "1.0.0",
	})
	if !result.IsInstalled {
		t.Error("expected installed = true")
	}
	if result.CurrentVersion != "1.0.0" {
		t.Errorf("version = %v, want 1.0.0", result.CurrentVersion)
	}
}

func TestAppValidatorService_List(t *testing.T) {
	svc := NewAppValidatorService()
	apps := svc.List()
	if len(apps) == 0 {
		t.Error("expected default apps")
	}
}

func TestFileSystemService_New(t *testing.T) {
	svc := NewFileSystemService()
	if svc == nil {
		t.Error("NewFileSystemService() returned nil")
	}
}

func TestFileSystemService_Exists(t *testing.T) {
	svc := NewFileSystemService()
	if !svc.Exists("/tmp") {
		t.Error("/tmp should exist")
	}
}

func TestFileSystemService_Read(t *testing.T) {
	svc := NewFileSystemService()
	data, err := svc.Read("/etc/hostname")
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty data")
	}
}

func TestFileSystemService_Write(t *testing.T) {
	svc := NewFileSystemService()
	tmp := t.TempDir() + "/test.txt"
	err := svc.Write(tmp, []byte("test"))
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}
	data, _ := svc.Read(tmp)
	if string(data) != "test" {
		t.Errorf("data = %v, want test", string(data))
	}
}

func TestFileSystemService_MkdirAll(t *testing.T) {
	svc := NewFileSystemService()
	tmp := t.TempDir() + "/testdir/subdir"
	err := svc.MkdirAll(tmp)
	if err != nil {
		t.Errorf("MkdirAll() error = %v", err)
	}
	if !svc.Exists(tmp) {
		t.Error("directory should exist")
	}
}

func TestAppValidatorService_Validate(t *testing.T) {
	svc := NewAppValidatorService()
	svc.WithExecutor(&mockExecutor{found: true, version: "1.0.0"})
	apps := svc.List()
	result := svc.Validate(apps)
	if len(result) == 0 {
		t.Error("expected validation results")
	}
}

type mockExecutor struct {
	found   bool
	version string
	err    error
}

func (m *mockExecutor) Run(cmd string, args ...string) (string, error) {
	if !m.found {
		return "", m.err
	}
	return m.version, nil
}

func (m *mockExecutor) Which(cmd string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return "/usr/bin/" + cmd, nil
}