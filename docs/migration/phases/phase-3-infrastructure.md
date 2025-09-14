# Phase 3: Infrastructure Layer (Adapters)

## 🎯 Objective

Implement infrastructure adapters that provide concrete implementations of the ports defined in the application layer, connecting the application to external systems.

## ⏱️ Duration: 4-5 days | Effort: 32-40h

## 📋 Prerequisites

- ✅ Phase 2 completed (Application layer implemented)
- ✅ Use cases and ports defined
- ✅ DTOs and application services ready

## 🎯 Goals

1. **Git Adapter**: Implement Git operations using go-git library
2. **Configuration Adapter**: JSON-based configuration management
3. **Docker Adapter**: Docker operations using Docker SDK
4. **Validation Adapter**: System command and environment validation
5. **Database Adapter**: Repository persistence (SQLite/PostgreSQL)
6. **Dependency Injection**: Complete DI container setup

## 📁 Directory Structure

Create the following structure in `internal/infrastructure/`:

```text
internal/infrastructure/
├── adapters/              # External system adapters
│   ├── git/
│   │   ├── git_adapter.go
│   │   └── git_adapter_test.go
│   ├── config/
│   │   ├── json_adapter.go
│   │   └── env_adapter.go
│   ├── docker/
│   │   ├── docker_adapter.go
│   │   └── docker_client.go
│   ├── validation/
│   │   ├── system_validator.go
│   │   └── command_checker.go
│   ├── persistence/
│   │   ├── sqlite_repository.go
│   │   └── memory_repository.go
│   └── notification/
│       ├── logger_adapter.go
│       └── console_adapter.go
├── config/                # Configuration management
│   ├── container.go       # DI container
│   ├── config.go          # Configuration loading
│   └── environment.go     # Environment variables
└── server/                # Server setup (future CLI/API)
    ├── handlers/          # Request handlers
    └── middleware/        # Middleware components
```

## 🚀 Implementation Steps

### Step 1: Implement Git Adapter

#### Create `internal/infrastructure/adapters/git/git_adapter.go`

```go
package git

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing"
    "github.com/go-git/go-git/v5/plumbing/transport/http"

    "github.com/fabianoflorentino/whiterose/internal/application/ports"
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
    "github.com/fabianoflorentino/whiterose/internal/domain/errors"
)

// GitAdapter implements the GitRepositoryPort interface
type GitAdapter struct {
    timeout    time.Duration
    auth       *http.BasicAuth
}

// GitConfig holds configuration for Git operations
type GitConfig struct {
    Timeout  time.Duration
    Username string
    Token    string
}

// NewGitAdapter creates a new Git adapter with configuration
func NewGitAdapter(config GitConfig) *GitAdapter {
    adapter := &GitAdapter{
        timeout: config.Timeout,
    }

    // Setup authentication if credentials provided
    if config.Username != "" && config.Token != "" {
        adapter.auth = &http.BasicAuth{
            Username: config.Username,
            Password: config.Token,
        }
    }

    return adapter
}

// Clone clones a repository to the specified local path
func (g *GitAdapter) Clone(ctx context.Context, repo *entities.Repository, localPath string) error {
    // Ensure parent directory exists
    if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
        return fmt.Errorf("failed to create parent directory: %w", err)
    }

    // Setup clone options
    cloneOptions := &git.CloneOptions{
        URL:           repo.URL().String(),
        ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", repo.Branch())),
        SingleBranch:  true,
        Depth:         1, // Shallow clone for faster operation
    }

    // Add authentication if available
    if g.auth != nil {
        cloneOptions.Auth = g.auth
    }

    // Create context with timeout
    timeoutCtx, cancel := context.WithTimeout(ctx, g.timeout)
    defer cancel()

    // Perform clone
    _, err := git.PlainCloneContext(timeoutCtx, localPath, false, cloneOptions)
    if err != nil {
        // Clean up on failure
        os.RemoveAll(localPath)
        return fmt.Errorf("failed to clone repository %s: %w", repo.Name(), err)
    }

    return nil
}

// Pull updates the local repository with remote changes
func (g *GitAdapter) Pull(ctx context.Context, localPath string) error {
    repo, err := git.PlainOpen(localPath)
    if err != nil {
        return fmt.Errorf("failed to open repository at %s: %w", localPath, err)
    }

    workTree, err := repo.Worktree()
    if err != nil {
        return fmt.Errorf("failed to get worktree: %w", err)
    }

    // Setup pull options
    pullOptions := &git.PullOptions{}
    if g.auth != nil {
        pullOptions.Auth = g.auth
    }

    // Create context with timeout
    timeoutCtx, cancel := context.WithTimeout(ctx, g.timeout)
    defer cancel()

    err = workTree.PullContext(timeoutCtx, pullOptions)
    if err != nil && err != git.NoErrAlreadyUpToDate {
        return fmt.Errorf("failed to pull changes: %w", err)
    }

    return nil
}

// Checkout switches to the specified branch
func (g *GitAdapter) Checkout(ctx context.Context, localPath, branch string) error {
    repo, err := git.PlainOpen(localPath)
    if err != nil {
        return fmt.Errorf("failed to open repository at %s: %w", localPath, err)
    }

    workTree, err := repo.Worktree()
    if err != nil {
        return fmt.Errorf("failed to get worktree: %w", err)
    }

    // Checkout the branch
    err = workTree.Checkout(&git.CheckoutOptions{
        Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
    })
    if err != nil {
        // Try to checkout as new branch from origin
        err = workTree.Checkout(&git.CheckoutOptions{
            Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
            Create: true,
        })
        if err != nil {
            return fmt.Errorf("failed to checkout branch %s: %w", branch, err)
        }
    }

    return nil
}

// GetCurrentBranch returns the current branch name
func (g *GitAdapter) GetCurrentBranch(ctx context.Context, localPath string) (string, error) {
    repo, err := git.PlainOpen(localPath)
    if err != nil {
        return "", fmt.Errorf("failed to open repository at %s: %w", localPath, err)
    }

    head, err := repo.Head()
    if err != nil {
        return "", fmt.Errorf("failed to get HEAD reference: %w", err)
    }

    if head.Name().IsBranch() {
        return head.Name().Short(), nil
    }

    return "", errors.NewBusinessRuleError("HEAD is not pointing to a branch")
}

// IsClean checks if the repository has uncommitted changes
func (g *GitAdapter) IsClean(ctx context.Context, localPath string) (bool, error) {
    repo, err := git.PlainOpen(localPath)
    if err != nil {
        return false, fmt.Errorf("failed to open repository at %s: %w", localPath, err)
    }

    workTree, err := repo.Worktree()
    if err != nil {
        return false, fmt.Errorf("failed to get worktree: %w", err)
    }

    status, err := workTree.Status()
    if err != nil {
        return false, fmt.Errorf("failed to get status: %w", err)
    }

    return status.IsClean(), nil
}

// Compile-time check
var _ ports.GitRepositoryPort = (*GitAdapter)(nil)
```

### Step 2: Implement Configuration Adapter

#### Create `internal/infrastructure/adapters/config/json_adapter.go`

```go
package config

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"

    "github.com/fabianoflorentino/whiterose/internal/application/ports"
    "github.com/fabianoflorentino/whiterose/internal/domain/errors"
)

// JSONConfigAdapter implements ConfigurationPort using JSON files
type JSONConfigAdapter struct {
    configPath string
}

// NewJSONConfigAdapter creates a new JSON configuration adapter
func NewJSONConfigAdapter(configPath string) *JSONConfigAdapter {
    return &JSONConfigAdapter{
        configPath: configPath,
    }
}

// Configuration represents the application configuration
type Configuration struct {
    WorkingDirectory string                    `json:"working_directory"`
    Repositories     []*RepositoryConfig       `json:"repositories"`
    Docker           *DockerConfig             `json:"docker"`
    Git              *GitConfig                `json:"git"`
    Environment      map[string]string         `json:"environment"`
    Validation       *ValidationConfig         `json:"validation"`
}

// RepositoryConfig represents repository configuration
type RepositoryConfig struct {
    Name        string            `json:"name"`
    URL         string            `json:"url"`
    Branch      string            `json:"branch"`
    Path        string            `json:"path"`
    Tags        []string          `json:"tags,omitempty"`
    Environment map[string]string `json:"environment,omitempty"`
}

// DockerConfig represents Docker configuration
type DockerConfig struct {
    Registry    string            `json:"registry"`
    Images      []*DockerImage    `json:"images"`
    Networks    []string          `json:"networks"`
    Volumes     map[string]string `json:"volumes"`
    ComposeFile string            `json:"compose_file,omitempty"`
}

// DockerImage represents Docker image configuration
type DockerImage struct {
    Name       string            `json:"name"`
    Tag        string            `json:"tag"`
    Dockerfile string            `json:"dockerfile"`
    Context    string            `json:"context"`
    Args       map[string]string `json:"args"`
    Ports      []string          `json:"ports,omitempty"`
}

// GitConfig represents Git configuration
type GitConfig struct {
    Username      string        `json:"username,omitempty"`
    Email         string        `json:"email,omitempty"`
    DefaultBranch string        `json:"default_branch"`
    Timeout       string        `json:"timeout"`
}

// ValidationConfig represents validation configuration
type ValidationConfig struct {
    RequiredCommands []CommandRequirement `json:"required_commands"`
    OptionalCommands []CommandRequirement `json:"optional_commands"`
    SkipValidation   bool                 `json:"skip_validation"`
}

// CommandRequirement represents a command requirement
type CommandRequirement struct {
    Command    string `json:"command"`
    MinVersion string `json:"min_version,omitempty"`
    Required   bool   `json:"required"`
}

// LoadConfig loads configuration from JSON file
func (c *JSONConfigAdapter) LoadConfig(ctx context.Context) (*Configuration, error) {
    // Check if config file exists
    if _, err := os.Stat(c.configPath); os.IsNotExist(err) {
        // Create default configuration if file doesn't exist
        config := c.getDefaultConfig()
        if err := c.SaveConfig(ctx, config); err != nil {
            return nil, fmt.Errorf("failed to create default config: %w", err)
        }
        return config, nil
    }

    // Read config file
    data, err := os.ReadFile(c.configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    // Parse JSON
    var config Configuration
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, errors.NewValidationError("invalid configuration format", err)
    }

    // Validate configuration
    if err := c.validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

// SaveConfig saves configuration to JSON file
func (c *JSONConfigAdapter) SaveConfig(ctx context.Context, config *Configuration) error {
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

    // Write to file atomically
    tmpFile := c.configPath + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        return fmt.Errorf("failed to write temp config file: %w", err)
    }

    if err := os.Rename(tmpFile, c.configPath); err != nil {
        os.Remove(tmpFile) // Cleanup temp file
        return fmt.Errorf("failed to replace config file: %w", err)
    }

    return nil
}

// GetRepositories returns configured repositories
func (c *JSONConfigAdapter) GetRepositories(ctx context.Context) ([]*RepositoryConfig, error) {
    config, err := c.LoadConfig(ctx)
    if err != nil {
        return nil, err
    }

    return config.Repositories, nil
}

// validateConfig validates the configuration structure
func (c *JSONConfigAdapter) validateConfig(config *Configuration) error {
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
func (c *JSONConfigAdapter) getDefaultConfig() *Configuration {
    return &Configuration{
        WorkingDirectory: "./repositories",
        Repositories:     []*RepositoryConfig{},
        Docker: &DockerConfig{
            Registry: "docker.io",
            Images:   []*DockerImage{},
            Networks: []string{"default"},
            Volumes:  make(map[string]string),
        },
        Git: &GitConfig{
            DefaultBranch: "main",
            Timeout:       "30s",
        },
        Environment: make(map[string]string),
        Validation: &ValidationConfig{
            RequiredCommands: []CommandRequirement{
                {Command: "git", MinVersion: "2.0.0", Required: true},
                {Command: "docker", MinVersion: "20.0.0", Required: true},
            },
            OptionalCommands: []CommandRequirement{
                {Command: "docker-compose", MinVersion: "1.29.0", Required: false},
                {Command: "go", MinVersion: "1.20.0", Required: false},
            },
        },
    }
}

// Compile-time check
var _ ports.ConfigurationPort = (*JSONConfigAdapter)(nil)
```

### Step 3: Implement Docker Adapter

#### Create `internal/infrastructure/adapters/docker/docker_adapter.go`

```go
package docker

import (
    "context"
    "fmt"
    "io"
    "os"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/client"
    "github.com/docker/docker/pkg/archive"

    "github.com/fabianoflorentino/whiterose/internal/application/ports"
)

// DockerAdapter implements the DockerPort interface
type DockerAdapter struct {
    client *client.Client
}

// NewDockerAdapter creates a new Docker adapter
func NewDockerAdapter() (*DockerAdapter, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, fmt.Errorf("failed to create Docker client: %w", err)
    }

    return &DockerAdapter{
        client: cli,
    }, nil
}

// BuildImage builds a Docker image
func (d *DockerAdapter) BuildImage(ctx context.Context, imageName, dockerfile, buildContext string, args map[string]string) (string, error) {
    // Create build context
    buildCtx, err := archive.TarWithOptions(buildContext, &archive.TarOptions{})
    if err != nil {
        return "", fmt.Errorf("failed to create build context: %w", err)
    }
    defer buildCtx.Close()

    // Prepare build args
    buildArgs := make(map[string]*string)
    for key, value := range args {
        buildArgs[key] = &value
    }

    // Build options
    buildOptions := types.ImageBuildOptions{
        Tags:       []string{imageName},
        Dockerfile: dockerfile,
        BuildArgs:  buildArgs,
        Remove:     true, // Remove intermediate containers
    }

    // Build image
    response, err := d.client.ImageBuild(ctx, buildCtx, buildOptions)
    if err != nil {
        return "", fmt.Errorf("failed to build image: %w", err)
    }
    defer response.Body.Close()

    // Read build output (you might want to process this differently)
    _, err = io.Copy(os.Stdout, response.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read build output: %w", err)
    }

    // Get the image ID
    images, err := d.client.ImageList(ctx, types.ImageListOptions{})
    if err != nil {
        return "", fmt.Errorf("failed to list images: %w", err)
    }

    for _, image := range images {
        for _, tag := range image.RepoTags {
            if tag == imageName {
                return image.ID, nil
            }
        }
    }

    return "", fmt.Errorf("built image not found")
}

// ListImages returns a list of Docker images
func (d *DockerAdapter) ListImages(ctx context.Context) ([]DockerImage, error) {
    images, err := d.client.ImageList(ctx, types.ImageListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to list images: %w", err)
    }

    result := make([]DockerImage, len(images))
    for i, img := range images {
        result[i] = DockerImage{
            ID:      img.ID,
            Tags:    img.RepoTags,
            Size:    img.Size,
            Created: img.Created,
        }
    }

    return result, nil
}

// RemoveImage removes a Docker image
func (d *DockerAdapter) RemoveImage(ctx context.Context, imageID string) error {
    _, err := d.client.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
        Force:         true,
        PruneChildren: true,
    })
    if err != nil {
        return fmt.Errorf("failed to remove image %s: %w", imageID, err)
    }

    return nil
}

// DockerImage represents a Docker image
type DockerImage struct {
    ID      string   `json:"id"`
    Tags    []string `json:"tags"`
    Size    int64    `json:"size"`
    Created int64    `json:"created"`
}

// Close closes the Docker client
func (d *DockerAdapter) Close() error {
    return d.client.Close()
}

// Compile-time check
var _ ports.DockerPort = (*DockerAdapter)(nil)
```

### Step 4: Implement Dependency Injection Container

#### Create `internal/infrastructure/config/container.go`

```go
package config

import (
    "fmt"
    "time"

    "github.com/fabianoflorentino/whiterose/internal/application/usecases"
    "github.com/fabianoflorentino/whiterose/internal/application/services"
    "github.com/fabianoflorentino/whiterose/internal/domain/repositories"
    gitAdapter "github.com/fabianoflorentino/whiterose/internal/infrastructure/adapters/git"
    configAdapter "github.com/fabianoflorentino/whiterose/internal/infrastructure/adapters/config"
    dockerAdapter "github.com/fabianoflorentino/whiterose/internal/infrastructure/adapters/docker"
    validationAdapter "github.com/fabianoflorentino/whiterose/internal/infrastructure/adapters/validation"
    persistenceAdapter "github.com/fabianoflorentino/whiterose/internal/infrastructure/adapters/persistence"
    notificationAdapter "github.com/fabianoflorentino/whiterose/internal/infrastructure/adapters/notification"
)

// Container holds all application dependencies
type Container struct {
    // Configuration
    config *AppConfig

    // Infrastructure Adapters
    gitAdapter          *gitAdapter.GitAdapter
    configAdapter       *configAdapter.JSONConfigAdapter
    dockerAdapter       *dockerAdapter.DockerAdapter
    validationAdapter   *validationAdapter.SystemValidator
    repositoryAdapter   repositories.RepositoryRepository
    notificationAdapter *notificationAdapter.LoggerAdapter

    // Application Services
    orchestrationService *services.OrchestrationService

    // Use Cases
    setupRepositoriesUC     *usecases.SetupRepositoriesUseCase
    validatePrerequisitesUC *usecases.ValidatePrerequisitesUseCase
    manageDockerImagesUC    *usecases.ManageDockerImagesUseCase
}

// AppConfig holds application configuration
type AppConfig struct {
    ConfigPath          string
    WorkingDir          string
    DatabasePath        string
    GitTimeout          time.Duration
    DockerTimeout       time.Duration
    LogLevel            string
    EnableNotifications bool
}

// NewContainer creates and configures a new dependency injection container
func NewContainer(config *AppConfig) (*Container, error) {
    container := &Container{
        config: config,
    }

    // Initialize infrastructure adapters
    if err := container.initializeAdapters(); err != nil {
        return nil, fmt.Errorf("failed to initialize adapters: %w", err)
    }

    // Initialize application services
    if err := container.initializeServices(); err != nil {
        return nil, fmt.Errorf("failed to initialize services: %w", err)
    }

    // Initialize use cases
    if err := container.initializeUseCases(); err != nil {
        return nil, fmt.Errorf("failed to initialize use cases: %w", err)
    }

    return container, nil
}

// initializeAdapters creates and configures all infrastructure adapters
func (c *Container) initializeAdapters() error {
    // Git adapter
    gitConfig := gitAdapter.GitConfig{
        Timeout: c.config.GitTimeout,
    }
    c.gitAdapter = gitAdapter.NewGitAdapter(gitConfig)

    // Configuration adapter
    c.configAdapter = configAdapter.NewJSONConfigAdapter(c.config.ConfigPath)

    // Docker adapter
    dockerAdp, err := dockerAdapter.NewDockerAdapter()
    if err != nil {
        return fmt.Errorf("failed to create Docker adapter: %w", err)
    }
    c.dockerAdapter = dockerAdp

    // Validation adapter
    c.validationAdapter = validationAdapter.NewSystemValidator()

    // Repository adapter (choose between SQLite and in-memory)
    if c.config.DatabasePath != "" {
        repoAdapter, err := persistenceAdapter.NewSQLiteRepository(c.config.DatabasePath)
        if err != nil {
            return fmt.Errorf("failed to create SQLite repository: %w", err)
        }
        c.repositoryAdapter = repoAdapter
    } else {
        c.repositoryAdapter = persistenceAdapter.NewInMemoryRepository()
    }

    // Notification adapter
    c.notificationAdapter = notificationAdapter.NewLoggerAdapter(c.config.LogLevel)

    return nil
}

// initializeServices creates and configures application services
func (c *Container) initializeServices() error {
    c.orchestrationService = services.NewOrchestrationService(
        c.repositoryAdapter,
        c.configAdapter,
        c.validationAdapter,
        c.notificationAdapter,
    )

    return nil
}

// initializeUseCases creates and configures all use cases
func (c *Container) initializeUseCases() error {
    // Setup Repositories Use Case
    c.setupRepositoriesUC = usecases.NewSetupRepositoriesUseCase(
        c.repositoryAdapter,
        c.gitAdapter,
        c.configAdapter,
        c.notificationAdapter,
    )

    // Validate Prerequisites Use Case
    c.validatePrerequisitesUC = usecases.NewValidatePrerequisitesUseCase(
        c.validationAdapter,
        c.notificationAdapter,
    )

    // Manage Docker Images Use Case
    c.manageDockerImagesUC = usecases.NewManageDockerImagesUseCase(
        c.dockerAdapter,
        c.configAdapter,
        c.notificationAdapter,
    )

    return nil
}

// Getter methods for accessing dependencies

func (c *Container) GetSetupRepositoriesUseCase() *usecases.SetupRepositoriesUseCase {
    return c.setupRepositoriesUC
}

func (c *Container) GetValidatePrerequisitesUseCase() *usecases.ValidatePrerequisitesUseCase {
    return c.validatePrerequisitesUC
}

func (c *Container) GetManageDockerImagesUseCase() *usecases.ManageDockerImagesUseCase {
    return c.manageDockerImagesUC
}

func (c *Container) GetOrchestrationService() *services.OrchestrationService {
    return c.orchestrationService
}

// Close cleans up resources
func (c *Container) Close() error {
    var errors []error

    if c.dockerAdapter != nil {
        if err := c.dockerAdapter.Close(); err != nil {
            errors = append(errors, err)
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("failed to close container: %v", errors)
    }

    return nil
}
```

## ✅ Acceptance Criteria

- [ ] All infrastructure adapters implemented
- [ ] Git operations working with go-git
- [ ] Configuration management with JSON
- [ ] Docker operations functional
- [ ] System validation implemented
- [ ] Database persistence working
- [ ] Dependency injection container complete
- [ ] Integration tests passing (>80% coverage)
- [ ] All external dependencies properly abstracted

## 🧪 Testing Strategy

### Integration Tests

```go
func TestGitAdapter_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    adapter := git.NewGitAdapter(git.GitConfig{
        Timeout: 30 * time.Second,
    })

    repo, err := entities.NewRepository("test-repo", "https://github.com/git/git.git", "master")
    require.NoError(t, err)

    tmpDir := t.TempDir()
    localPath := filepath.Join(tmpDir, "test-repo")

    err = adapter.Clone(context.Background(), repo, localPath)
    assert.NoError(t, err)

    // Verify repository was cloned
    assert.DirExists(t, localPath)
    assert.FileExists(t, filepath.Join(localPath, ".git"))
}
```

## 📚 Next Steps

After completing Phase 3:

1. ✅ All infrastructure adapters implemented and tested
2. ✅ External dependencies properly abstracted
3. ✅ Dependency injection container working
4. 🚀 Proceed to [Phase 4: Migration & Integration](phase-4-migration.md)

## 🔗 Related Documentation

- [Application Layer Guide](phase-2-application.md)
- [Code Examples - Infrastructure Layer](../code-examples/infrastructure/)
- [Architecture Overview](../architecture/current-vs-proposed.md)
- [Testing Guide](../how-to/testing-guide.md)
