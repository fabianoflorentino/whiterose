package entities

import (
	"fmt"
	"strings"
	"time"
)

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

type OperatingSystem string

const (
	OSLinux   OperatingSystem = "linux"
	OSDarwin  OperatingSystem = "darwin"
	OSWindows OperatingSystem = "windows"
)

type ApplicationStatus struct {
	Application      *Application
	IsInstalled      bool
	InstalledVersion string
	IsUpToDate       bool
	ErrorMessage     string
	CheckedAt        time.Time
}

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

func (a *Application) AddInstallInstruction(os OperatingSystem, instruction string) error {
	if instruction == "" {
		return fmt.Errorf("instruction cannot be empty")
	}

	a.InstallInstructions[string(os)] = instruction
	a.UpdatedAt = time.Now()

	return nil
}

// Helper functions
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

func generateApplicationID(name string) string {
	normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return fmt.Sprintf("app_%s_%d", normalizedName, time.Now().Unix()%10000)
}
