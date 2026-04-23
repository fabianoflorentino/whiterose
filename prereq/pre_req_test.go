package prereq

import (
	"runtime"
	"testing"

	"github.com/fabianoflorentino/whiterose/utils"
)

func TestAppValidator_New(t *testing.T) {
	av := NewAppValidator()
	if av == nil {
		t.Error("NewAppValidator() returned nil")
	}
}

func TestAppValidator_AddApp(t *testing.T) {
	av := &AppValidator{
		apps: []utils.AppInfo{},
		os:   runtime.GOOS,
	}

	av.AddApp(utils.AppInfo{
		Name:        "Go",
		Command:     "go",
		VersionFlag: "version",
	})

	if len(av.apps) != 1 {
		t.Errorf("len(apps) = %d, want 1", len(av.apps))
	}

	if av.apps[0].Name != "Go" {
		t.Errorf("apps[0].Name = %v, want Go", av.apps[0].Name)
	}
}

func TestAppValidator_getOSName(t *testing.T) {
	tests := []struct {
		name string
		os   string
		want string
	}{
		{"darwin", "darwin", "macOS"},
		{"linux", "linux", "Linux"},
		{"windows", "windows", "Windows"},
		{"unknown", "freebsd", "freebsd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := &AppValidator{os: tt.os}
			if got := av.getOSName(); got != tt.want {
				t.Errorf("getOSName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppValidator_ValidateApps_Empty(t *testing.T) {
	av := &AppValidator{
		apps: []utils.AppInfo{},
		os:   runtime.GOOS,
	}

	av.ValidateApps()
}

func TestAppValidator_ValidateSpecificApps_NoApps(t *testing.T) {
	av := &AppValidator{
		apps: []utils.AppInfo{},
		os:   runtime.GOOS,
	}

	av.ValidateSpecificApps([]string{"nonexistent"})
}

func TestAppValidator_ListAvailableApps_Empty(t *testing.T) {
	av := &AppValidator{
		apps: []utils.AppInfo{},
		os:   runtime.GOOS,
	}

	av.ListAvailableApps()
}

func TestAppValidator_ListAvailableApps_WithApps(t *testing.T) {
	av := &AppValidator{
		apps: []utils.AppInfo{
			{Name: "Go", Command: "go"},
			{Name: "Git", Command: "git"},
		},
		os: runtime.GOOS,
	}

	av.ListAvailableApps()
}

func TestCheckAppInstalled_NotFound(t *testing.T) {
	av := &AppValidator{
		apps: []utils.AppInfo{},
		os:   runtime.GOOS,
	}

	_, _, err := av.checkAppInstalled(utils.AppInfo{Command: "nonexistent-cmd"})
	if err == nil {
		t.Error("expected error for nonexistent command")
	}
}

func TestCheckAppInstalled_Found(t *testing.T) {
	if testing.Short() {
		t.Skip("requires system commands")
	}

	av := &AppValidator{
		apps: []utils.AppInfo{},
		os:   runtime.GOOS,
	}

	_, _, err := av.checkAppInstalled(utils.AppInfo{Command: "go", VersionFlag: "version"})
	if err != nil {
		t.Logf("go may not be installed: %v", err)
	}
}
