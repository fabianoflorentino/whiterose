package entities

import (
	"testing"
)

func TestNewApplication(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		version   string
		wantErr   bool
	}{
		{"Go", "go", "version", false},
		{"Docker", "docker", "--version", false},
		{"", "go", "version", true},
		{"Go", "", "version", true},
		{"Go", "go", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApplication(tt.name, tt.command, tt.version, "1.0.0")
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApplication() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Name != tt.name {
				t.Errorf("Name = %v, want %v", got.Name, tt.name)
			}
		})
	}
}

func TestApplication_AddInstallInstruction(t *testing.T) {
	app := &Application{
		Name:                "Go",
		Command:             "go",
		VersionFlag:         "version",
		InstallInstructions: make(map[string]string),
	}

	tests := []struct {
		name        string
		os         OperatingSystem
		instruction string
		wantErr    bool
	}{
		{"linux", OSLinux, "sudo apt install golang", false},
		{"darwin", OSDarwin, "brew install go", false},
		{"windows", OSWindows, "choco install golang", false},
		{"empty instruction", OSLinux, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.AddInstallInstruction(tt.os, tt.instruction)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddInstallInstruction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_Validate(t *testing.T) {
	tests := []struct {
		name    string
		app    *Application
		wantErr bool
	}{
		{"valid", &Application{Name: "Go", Command: "go", VersionFlag: "version"}, false},
		{"invalid name", &Application{Name: "", Command: "go", VersionFlag: "version"}, true},
		{"invalid command", &Application{Name: "Go", Command: "", VersionFlag: "version"}, true},
		{"invalid version", &Application{Name: "Go", Command: "go", VersionFlag: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.app.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewApplicationStatus(t *testing.T) {
	app := &Application{Name: "Go", Command: "go", VersionFlag: "version"}
	status := NewApplicationStatus(app)

	if status.Application != app {
		t.Error("Application not set correctly")
	}
	if status.IsInstalled != false {
		t.Errorf("IsInstalled = %v, want false", status.IsInstalled)
	}
}

func TestApplicationStatus_SetInstalled(t *testing.T) {
	app := &Application{Name: "Go", Command: "go", VersionFlag: "version", RecommendationVersion: "1.20.0"}
	status := NewApplicationStatus(app)
	status.SetInstalled("1.20.0")

	if !status.IsInstalled {
		t.Error("IsInstalled should be true")
	}
	if status.InstalledVersion != "1.20.0" {
		t.Errorf("InstalledVersion = %v, want 1.20.0", status.InstalledVersion)
	}
	if !status.IsUpToDate {
		t.Error("IsUpToDate should be true")
	}
}

func TestApplicationStatus_SetNotInstalled(t *testing.T) {
	app := &Application{Name: "Go", Command: "go", VersionFlag: "version"}
	status := NewApplicationStatus(app)
	status.SetNotInstalled("Go not found")

	if status.IsInstalled {
		t.Error("IsInstalled should be false")
	}
	if status.ErrorMessage != "Go not found" {
		t.Errorf("ErrorMessage = %v, want 'Go not found'", status.ErrorMessage)
	}
}