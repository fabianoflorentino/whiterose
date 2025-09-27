package entities

import (
	"fmt"
	"strings"
	"time"
)

// Application represents a software application with its metadata and installation instructions.
type Application struct {
	ID                    string
	Name                  string
	Command               string
	VersionFlag           string
	RecommendationVersion string
	InstallInstructions   map[string]string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// OperatingSystem defines supported operating systems for installation instructions.
type OperatingSystem string

const (
	OSLinux   OperatingSystem = "linux"
	OSDarwin  OperatingSystem = "darwin"
	OSWindows OperatingSystem = "windows"
)

// ApplicationStatus represents the installation status of an application on a system.
type ApplicationStatus struct {
	Application      *Application
	IsInstalled      bool
	InstalledVersion string
	IsUpToDate       bool
	ErrorMessage     string
	CheckedAt        time.Time
}

// NewApplication creates a new Application instance after validating the input data.
func NewApplication(name, command, versionFlag, recommendedVersion string) (*Application, error) {
	if err := validateApplicationData(name, command, versionFlag); err != nil {
		return nil, err
	}

	return &Application{
		ID:                    generateApplicationID(name),
		Name:                  name,
		Command:               command,
		RecommendationVersion: recommendedVersion,
		InstallInstructions:   make(map[string]string),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}, nil
}

// AddInstallInstruction adds or updates installation instructions for a specific operating system.
func (a *Application) AddInstallInstruction(os OperatingSystem, instruction string) error {
	if instruction == "" {
		return fmt.Errorf("instruction cannot be empty")
	}

	a.InstallInstructions[string(os)] = instruction
	a.UpdatedAt = time.Now()

	return nil
}

// Validate checks if the Application instance has valid data.
func (a *Application) Validate() error {
	return validateApplicationData(a.Name, a.Command, a.VersionFlag)
}

// NewApplicationStatus creates a new ApplicationStatus instance for the given application.
func NewApplicationStatus(a *Application) *ApplicationStatus {
	return &ApplicationStatus{
		Application: a,
		IsInstalled: false,
		IsUpToDate:  false,
		CheckedAt:   time.Now(),
	}
}

// SetInstalled updates the status to reflect that the application is installed with the given version.
func (as *ApplicationStatus) SetInstalled(version string) {
	as.IsInstalled = true
	as.InstalledVersion = version
	as.IsUpToDate = as.checkVersionCompatibility()
	as.ErrorMessage = ""
	as.CheckedAt = time.Now()
}

// SetNotInstalled updates the status to reflect that the application is not installed, with an optional error message.
func (as *ApplicationStatus) SetNotInstalled(errMsg string) {
	as.IsInstalled = false
	as.InstalledVersion = ""
	as.IsUpToDate = false
	as.ErrorMessage = errMsg
	as.CheckedAt = time.Now()
}

// Private method to check version compatibility

// checkVersionCompatibility compares the installed version with the recommended version.
func (as *ApplicationStatus) checkVersionCompatibility() bool {
	if as.Application.RecommendationVersion == "" || as.InstalledVersion == "" {
		return false
	}
	return as.Application.RecommendationVersion == as.InstalledVersion
}

// Helper functions

// validateApplicationData checks if the provided application data is valid.
func validateApplicationData(name, command, versionFlag string) error {
	if name == "" {
		return fmt.Errorf("application name cannot be empty")
	}

	if command == "" {
		return fmt.Errorf("application command cannot be empty")
	}

	if versionFlag == "" {
		return fmt.Errorf("application version flag cannot be empty")
	}

	return nil
}

// generateApplicationID creates a unique application ID based on the name and current timestamp.
func generateApplicationID(name string) string {
	normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return fmt.Sprintf("app_%s_%d", normalizedName, time.Now().Unix()%10000)
}
