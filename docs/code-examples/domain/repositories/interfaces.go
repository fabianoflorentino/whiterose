package repositories

import (
	"context"

	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/entities"
)

// RepositoryRepository defines the interface for repository management
// This follows the Repository pattern as a port in hexagonal architecture
type RepositoryRepository interface {
	// Save persists a repository entity
	Save(ctx context.Context, repo *entities.Repository) error

	// FindByID retrieves a repository by its ID
	FindByID(ctx context.Context, id string) (*entities.Repository, error)

	// FindByName retrieves a repository by its name
	FindByName(ctx context.Context, name string) (*entities.Repository, error)

	// FindAll retrieves all repositories
	FindAll(ctx context.Context) ([]*entities.Repository, error)

	// Update updates an existing repository
	Update(ctx context.Context, repo *entities.Repository) error

	// Delete removes a repository by ID
	Delete(ctx context.Context, id string) error

	// Exists checks if a repository exists by name
	Exists(ctx context.Context, name string) (bool, error)
}

// GitRepository defines the interface for Git operations
// This is a secondary port for external Git systems
type GitRepository interface {
	// Clone clones a repository to the specified local path
	Clone(ctx context.Context, repo *entities.Repository, localPath string) error

	// Pull updates the local repository with remote changes
	Pull(ctx context.Context, localPath string) error

	// Checkout switches to the specified branch
	Checkout(ctx context.Context, localPath, branch string) error

	// GetCurrentBranch returns the current branch name
	GetCurrentBranch(ctx context.Context, localPath string) (string, error)

	// ListBranches returns all available branches
	ListBranches(ctx context.Context, localPath string) ([]string, error)

	// IsClean checks if the repository has uncommitted changes
	IsClean(ctx context.Context, localPath string) (bool, error)

	// GetLastCommit returns information about the last commit
	GetLastCommit(ctx context.Context, localPath string) (*CommitInfo, error)
}

// ConfigurationRepository defines the interface for configuration management
type ConfigurationRepository interface {
	// LoadConfig loads configuration from storage
	LoadConfig(ctx context.Context) (*Configuration, error)

	// SaveConfig saves configuration to storage
	SaveConfig(ctx context.Context, config *Configuration) error

	// GetRepositories returns configured repositories
	GetRepositories(ctx context.Context) ([]*RepositoryConfig, error)

	// AddRepository adds a new repository configuration
	AddRepository(ctx context.Context, config *RepositoryConfig) error

	// RemoveRepository removes a repository configuration
	RemoveRepository(ctx context.Context, name string) error
}

// ValidationRepository defines the interface for system validation
type ValidationRepository interface {
	// CheckCommand verifies if a command is available in the system
	CheckCommand(ctx context.Context, command string) error

	// CheckVersion verifies if a command meets version requirements
	CheckVersion(ctx context.Context, command, minVersion string) error

	// GetSystemInfo returns system information
	GetSystemInfo(ctx context.Context) (*SystemInfo, error)

	// ValidateEnvironment checks if the environment is ready for setup
	ValidateEnvironment(ctx context.Context) ([]ValidationResult, error)
}

// Supporting types

// CommitInfo represents information about a Git commit
type CommitInfo struct {
	Hash      string
	Message   string
	Author    string
	Email     string
	Timestamp string
}

// Configuration represents application configuration
type Configuration struct {
	WorkingDirectory string              `json:"working_directory"`
	Repositories     []*RepositoryConfig `json:"repositories"`
	Docker           *DockerConfig       `json:"docker"`
	Environment      map[string]string   `json:"environment"`
}

// RepositoryConfig represents repository configuration
type RepositoryConfig struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
	Path   string `json:"path"`
}

// DockerConfig represents Docker configuration
type DockerConfig struct {
	Registry string            `json:"registry"`
	Images   []*DockerImage    `json:"images"`
	Networks []string          `json:"networks"`
	Volumes  map[string]string `json:"volumes"`
}

// DockerImage represents Docker image configuration
type DockerImage struct {
	Name       string            `json:"name"`
	Tag        string            `json:"tag"`
	Dockerfile string            `json:"dockerfile"`
	Context    string            `json:"context"`
	Args       map[string]string `json:"args"`
}

// SystemInfo represents system information
type SystemInfo struct {
	OS           string            `json:"os"`
	Architecture string            `json:"architecture"`
	Commands     map[string]string `json:"commands"` // command -> version
}

// ValidationResult represents a validation check result
type ValidationResult struct {
	Check   string `json:"check"`
	Status  string `json:"status"` // "pass", "fail", "warning"
	Message string `json:"message"`
	Error   error  `json:"error,omitempty"`
}
