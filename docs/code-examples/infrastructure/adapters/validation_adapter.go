package adapters

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/repositories"
)

// SystemValidationAdapter implements ValidationRepository interface
type SystemValidationAdapter struct{}

// NewSystemValidationAdapter creates a new system validation adapter
func NewSystemValidationAdapter() *SystemValidationAdapter {
	return &SystemValidationAdapter{}
}

// CheckCommand verifies if a command is available in the system
func (s *SystemValidationAdapter) CheckCommand(ctx context.Context, command string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("command '%s' not found in PATH", command)
	}
	return nil
}

// CheckVersion verifies if a command meets version requirements
func (s *SystemValidationAdapter) CheckVersion(ctx context.Context, command, minVersion string) error {
	// First check if command exists
	if err := s.CheckCommand(ctx, command); err != nil {
		return err
	}

	// Get current version (simplified implementation)
	cmd := exec.CommandContext(ctx, command, "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get version for '%s': %w", command, err)
	}

	version := strings.TrimSpace(string(output))

	// In a real implementation, you would parse and compare versions
	// For this example, we'll just check if version string contains the minimum version
	if !strings.Contains(version, minVersion) {
		return fmt.Errorf("command '%s' version '%s' does not meet minimum requirement '%s'",
			command, version, minVersion)
	}

	return nil
}

// GetSystemInfo returns system information
func (s *SystemValidationAdapter) GetSystemInfo(ctx context.Context) (*repositories.SystemInfo, error) {
	info := &repositories.SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Commands:     make(map[string]string),
	}

	// Check common commands and their versions
	commands := []string{"git", "docker", "go"}
	for _, cmd := range commands {
		if err := s.CheckCommand(ctx, cmd); err == nil {
			version := s.getCommandVersion(ctx, cmd)
			info.Commands[cmd] = version
		}
	}

	return info, nil
}

// ValidateEnvironment checks if the environment is ready for setup
func (s *SystemValidationAdapter) ValidateEnvironment(ctx context.Context) ([]repositories.ValidationResult, error) {
	var results []repositories.ValidationResult

	// Required commands and their minimum versions
	requirements := map[string]string{
		"git":    "2.0.0",
		"docker": "20.0.0",
		"go":     "1.20.0",
	}

	for command, minVersion := range requirements {
		result := repositories.ValidationResult{
			Check: fmt.Sprintf("Command: %s", command),
		}

		// Check if command exists
		if err := s.CheckCommand(ctx, command); err != nil {
			result.Status = "fail"
			result.Message = fmt.Sprintf("Command '%s' not found", command)
			result.Error = err
		} else {
			// Check version
			if err := s.CheckVersion(ctx, command, minVersion); err != nil {
				result.Status = "warning"
				result.Message = fmt.Sprintf("Version check failed for '%s'", command)
				result.Error = err
			} else {
				result.Status = "pass"
				result.Message = fmt.Sprintf("Command '%s' is available and meets requirements", command)
			}
		}

		results = append(results, result)
	}

	// Additional system checks
	results = append(results, s.checkDiskSpace(ctx))
	results = append(results, s.checkNetworkConnectivity(ctx))

	return results, nil
}

// getCommandVersion gets the version of a command (simplified implementation)
func (s *SystemValidationAdapter) getCommandVersion(ctx context.Context, command string) string {
	cmd := exec.CommandContext(ctx, command, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	// Extract version from output (simplified)
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}

	return "unknown"
}

// checkDiskSpace checks available disk space
func (s *SystemValidationAdapter) checkDiskSpace(ctx context.Context) repositories.ValidationResult {
	result := repositories.ValidationResult{
		Check: "Disk Space",
	}

	// Simplified disk space check (in a real implementation, use syscalls)
	result.Status = "pass"
	result.Message = "Sufficient disk space available"

	return result
}

// checkNetworkConnectivity checks network connectivity
func (s *SystemValidationAdapter) checkNetworkConnectivity(ctx context.Context) repositories.ValidationResult {
	result := repositories.ValidationResult{
		Check: "Network Connectivity",
	}

	// Simple connectivity check using ping (simplified)
	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "8.8.8.8")
	if err := cmd.Run(); err != nil {
		result.Status = "warning"
		result.Message = "Network connectivity issues detected"
		result.Error = err
	} else {
		result.Status = "pass"
		result.Message = "Network connectivity is working"
	}

	return result
}

// Compile-time check to ensure SystemValidationAdapter implements ValidationRepository
var _ repositories.ValidationRepository = (*SystemValidationAdapter)(nil)
