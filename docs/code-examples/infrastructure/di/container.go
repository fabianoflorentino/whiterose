package di

import (
	"path/filepath"

	"github.com/fabianoflorentino/whiterose/docs/code-examples/application/usecases"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/repositories"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/infrastructure/adapters"
)

// Container holds all application dependencies
type Container struct {
	// Configuration
	ConfigPath string
	WorkingDir string

	// Repositories (Infrastructure adapters)
	GitRepo        repositories.GitRepository
	ConfigRepo     repositories.ConfigurationRepository
	RepositoryRepo repositories.RepositoryRepository
	ValidationRepo repositories.ValidationRepository

	// Use Cases (Application layer)
	SetupRepositoriesUC *usecases.SetupRepositoriesUseCase
	// Add other use cases here as they're implemented
}

// NewContainer creates and configures a new dependency injection container
func NewContainer(configPath, workingDir string) (*Container, error) {
	container := &Container{
		ConfigPath: configPath,
		WorkingDir: workingDir,
	}

	// Initialize infrastructure adapters
	if err := container.initializeAdapters(); err != nil {
		return nil, err
	}

	// Initialize use cases
	if err := container.initializeUseCases(); err != nil {
		return nil, err
	}

	return container, nil
}

// initializeAdapters creates and configures all infrastructure adapters
func (c *Container) initializeAdapters() error {
	// Git adapter
	c.GitRepo = adapters.NewGitAdapter()

	// Configuration adapter
	c.ConfigRepo = adapters.NewConfigAdapter(c.ConfigPath)

	// Repository adapter (would typically be a database implementation)
	// For this example, we'll use an in-memory implementation
	c.RepositoryRepo = adapters.NewInMemoryRepositoryAdapter()

	// Validation adapter
	c.ValidationRepo = adapters.NewSystemValidationAdapter()

	return nil
}

// initializeUseCases creates and configures all use cases with their dependencies
func (c *Container) initializeUseCases() error {
	// Setup Repositories Use Case
	c.SetupRepositoriesUC = usecases.NewSetupRepositoriesUseCase(
		c.RepositoryRepo,
		c.GitRepo,
		c.ConfigRepo,
		c.WorkingDir,
	)

	return nil
}

// GetSetupRepositoriesUseCase returns the setup repositories use case
func (c *Container) GetSetupRepositoriesUseCase() *usecases.SetupRepositoriesUseCase {
	return c.SetupRepositoriesUC
}

// GetGitRepository returns the git repository adapter
func (c *Container) GetGitRepository() repositories.GitRepository {
	return c.GitRepo
}

// GetConfigRepository returns the configuration repository adapter
func (c *Container) GetConfigRepository() repositories.ConfigurationRepository {
	return c.ConfigRepo
}

// GetRepositoryRepository returns the repository repository adapter
func (c *Container) GetRepositoryRepository() repositories.RepositoryRepository {
	return c.RepositoryRepo
}

// GetValidationRepository returns the validation repository adapter
func (c *Container) GetValidationRepository() repositories.ValidationRepository {
	return c.ValidationRepo
}

// ContainerConfig holds configuration for dependency injection
type ContainerConfig struct {
	ConfigPath string
	WorkingDir string
	// Add other configuration options as needed
}

// NewContainerFromConfig creates a container from configuration
func NewContainerFromConfig(config ContainerConfig) (*Container, error) {
	// Set defaults if not provided
	if config.ConfigPath == "" {
		config.ConfigPath = filepath.Join(".", "config", "whiterose.json")
	}

	if config.WorkingDir == "" {
		config.WorkingDir = filepath.Join(".", "repositories")
	}

	return NewContainer(config.ConfigPath, config.WorkingDir)
}

// Example of how to use environment variables for configuration
func NewContainerFromEnv() (*Container, error) {
	config := ContainerConfig{
		ConfigPath: getEnvOrDefault("WHITEROSE_CONFIG_PATH", "./config/whiterose.json"),
		WorkingDir: getEnvOrDefault("WHITEROSE_WORKING_DIR", "./repositories"),
	}

	return NewContainerFromConfig(config)
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	// This would typically import from utils package
	// For this example, we'll inline a simple implementation
	return defaultValue
}
