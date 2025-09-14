package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/errors"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/repositories"
)

// ConfigAdapter implements the ConfigurationRepository interface
type ConfigAdapter struct {
	configPath string
}

// NewConfigAdapter creates a new configuration adapter
func NewConfigAdapter(configPath string) *ConfigAdapter {
	return &ConfigAdapter{
		configPath: configPath,
	}
}

// LoadConfig loads configuration from storage
func (c *ConfigAdapter) LoadConfig(ctx context.Context) (*repositories.Configuration, error) {
	// Check if config file exists
	if _, err := os.Stat(c.configPath); os.IsNotExist(err) {
		// Return default configuration if file doesn't exist
		return c.getDefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config repositories.Configuration
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, errors.NewValidationError("invalid configuration format", err)
	}

	// Validate configuration
	if err := c.validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves configuration to storage
func (c *ConfigAdapter) SaveConfig(ctx context.Context, config *repositories.Configuration) error {
	// Validate configuration
	if err := c.validateConfig(config); err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file
	if err := os.WriteFile(c.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetRepositories returns configured repositories
func (c *ConfigAdapter) GetRepositories(ctx context.Context) ([]*repositories.RepositoryConfig, error) {
	config, err := c.LoadConfig(ctx)
	if err != nil {
		return nil, err
	}

	return config.Repositories, nil
}

// AddRepository adds a new repository configuration
func (c *ConfigAdapter) AddRepository(ctx context.Context, repoConfig *repositories.RepositoryConfig) error {
	// Load current configuration
	config, err := c.LoadConfig(ctx)
	if err != nil {
		return err
	}

	// Check if repository already exists
	for _, existing := range config.Repositories {
		if existing.Name == repoConfig.Name {
			return errors.NewConflictError(fmt.Sprintf("repository %s already exists", repoConfig.Name))
		}
	}

	// Add new repository
	config.Repositories = append(config.Repositories, repoConfig)

	// Save updated configuration
	return c.SaveConfig(ctx, config)
}

// RemoveRepository removes a repository configuration
func (c *ConfigAdapter) RemoveRepository(ctx context.Context, name string) error {
	// Load current configuration
	config, err := c.LoadConfig(ctx)
	if err != nil {
		return err
	}

	// Find and remove repository
	found := false
	newRepos := make([]*repositories.RepositoryConfig, 0, len(config.Repositories))
	for _, repo := range config.Repositories {
		if repo.Name != name {
			newRepos = append(newRepos, repo)
		} else {
			found = true
		}
	}

	if !found {
		return errors.NewNotFoundError(fmt.Sprintf("repository %s not found", name))
	}

	// Update configuration
	config.Repositories = newRepos

	// Save updated configuration
	return c.SaveConfig(ctx, config)
}

// validateConfig validates the configuration structure
func (c *ConfigAdapter) validateConfig(config *repositories.Configuration) error {
	if config.WorkingDirectory == "" {
		return errors.NewValidationError("working directory cannot be empty", nil)
	}

	// Validate repositories
	repoNames := make(map[string]bool)
	for _, repo := range config.Repositories {
		if repo.Name == "" {
			return errors.NewValidationError("repository name cannot be empty", nil)
		}
		if repo.URL == "" {
			return errors.NewValidationError("repository URL cannot be empty", nil)
		}
		if repo.Branch == "" {
			return errors.NewValidationError("repository branch cannot be empty", nil)
		}

		// Check for duplicate names
		if repoNames[repo.Name] {
			return errors.NewValidationError(fmt.Sprintf("duplicate repository name: %s", repo.Name), nil)
		}
		repoNames[repo.Name] = true
	}

	return nil
}

// getDefaultConfig returns a default configuration
func (c *ConfigAdapter) getDefaultConfig() *repositories.Configuration {
	return &repositories.Configuration{
		WorkingDirectory: "./repositories",
		Repositories:     []*repositories.RepositoryConfig{},
		Docker: &repositories.DockerConfig{
			Registry: "docker.io",
			Images:   []*repositories.DockerImage{},
			Networks: []string{"default"},
			Volumes:  make(map[string]string),
		},
		Environment: make(map[string]string),
	}
}

// Compile-time check to ensure ConfigAdapter implements ConfigurationRepository
var _ repositories.ConfigurationRepository = (*ConfigAdapter)(nil)
