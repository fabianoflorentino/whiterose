// Package prereq provides utilities for validating the presence and versions of essential
// command-line applications required for development environments. It supports checking
// installation status, displaying recommended versions, and providing OS-specific installation
// instructions for applications such as Go, Git, Docker, jq, and yq.
//
// Types:
//
//   - AppInfo: Represents metadata about an application, including its name, command,
//     version flag, recommended version, and installation instructions per OS.
//   - AppValidator: Manages a list of applications to validate, detects the current OS,
//     and provides methods for validating all or specific applications.
//
// Functions:
//
//   - NewAppValidator: Constructs a new AppValidator pre-populated with common development tools.
//   - (*AppValidator) AddApp: Adds a custom application to the validator.
//   - (*AppValidator) ValidateApps: Checks all registered applications for installation and version.
//   - (*AppValidator) ValidateSpecificApps: Validates only the specified applications by name or command.
//   - (*AppValidator) ListAvailableApps: Lists all applications available for validation.
//   - (*AppValidator) getOSName: Returns a human-readable name for the current OS.
//   - (*AppValidator) checkAppInstalled: Checks if an application is installed and retrieves its version.
//
// Usage:
//
//	validator := prereq.NewAppValidator()
//	validator.ValidateApps()
//	validator.ValidateSpecificApps([]string{"Go", "Git"})
//	validator.ListAvailableApps()
package prereq

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fabianoflorentino/whiterose/utils"
)

// AppInfo holds information about a command-line application.
type AppInfo struct {
	Name                string
	Command             string
	VersionFlag         string
	RecommendedVersion  string
	InstallInstructions map[string]string
}

// AppValidator manages a list of applications to validate.
type AppValidator struct {
	apps []utils.AppInfo
	os   string
}

// NewAppValidator constructs a new AppValidator pre-populated with common development tools.
func NewAppValidator() *AppValidator {
	config, err := utils.LoadDotConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}

	apps, err := utils.FetchAppsInfo(config)
	if err != nil {
		fmt.Printf("Error fetching applications: %v\n", err)
		return &AppValidator{os: runtime.GOOS, apps: apps}
	}

	return &AppValidator{os: runtime.GOOS, apps: apps}
}

// AddApp adds a custom application to the validator.
func (av *AppValidator) AddApp(app utils.AppInfo) {
	av.apps = append(av.apps, app)
}

func (av *AppValidator) ValidateApps() {
	installedCount := 0

	for _, app := range av.apps {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		installed, version, err := av.checkAppInstalled(app)

		if installed && err == nil {
			fmt.Printf("💾 %s\n", app.Name)
			fmt.Printf("✅ Status: INSTALLED\n")
			fmt.Printf("📦 Version: %s\n", version)
			fmt.Printf("🎯 Recommended: %s\n", app.RecommendedVersion)
			installedCount++
		} else {
			fmt.Printf("❌ Status: NOT INSTALLED\n")
			fmt.Printf("🎯 Recommended Version: %s\n", app.RecommendedVersion)
			fmt.Printf("📥 Installation Instructions:\n")

			if instruction, exists := app.InstallInstructions[av.os]; exists {
				fmt.Printf("   %s\n", instruction)
			} else {
				fmt.Printf("   Instructions not available for %s\n", av.getOSName())
			}
		}

		fmt.Printf("\n")
	}
}

// ValidateSpecificApps validates only the specified applications by name or command.
// It accepts a slice of application names or commands to validate.
// If an application is not found in the predefined list, it is skipped with a message.
func (av *AppValidator) ValidateSpecificApps(appNames []string) {
	var appsToValidate []utils.AppInfo

	for _, name := range appNames {
		for _, app := range av.apps {
			if strings.EqualFold(app.Name, name) || strings.EqualFold(app.Command, name) {
				appsToValidate = append(appsToValidate, app)
				break
			}
		}
	}

	if len(appsToValidate) == 0 {
		fmt.Println("❌ No applications found in the list to validate.")
		return
	}

	// Temporarily replace the app list
	originalApps := av.apps
	av.apps = appsToValidate
	av.ValidateApps()
	av.apps = originalApps
}

// ListAvailableApps lists all applications available for validation.
func (av *AppValidator) ListAvailableApps() {
	fmt.Printf("📋 Available applications for validation:\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	for i, app := range av.apps {
		fmt.Printf("%d. %s (command: %s)\n", i+1, app.Name, app.Command)
	}
	fmt.Printf("\n")
}

func (av *AppValidator) getOSName() string {
	switch av.os {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return av.os
	}
}

// checkAppInstalled checks if a command-line application is installed and retrieves its version.
func (av *AppValidator) checkAppInstalled(app utils.AppInfo) (bool, string, error) {
	cmd := exec.Command(app.Command, app.VersionFlag)
	output, err := cmd.Output()
	if err != nil {
		return false, "", err
	}

	version := strings.TrimSpace(string(output))

	return true, version, nil
}
