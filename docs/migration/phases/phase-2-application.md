# Phase 2: Application Layer (Orchestration)

## 🎯 Objective

Implement use cases and define application contracts that orchestrate domain entities and coordinate business flows.

## ⏱️ Duration: 3-4 days | Effort: 24-30h

## 📋 Prerequisites

- ✅ Phase 1 completed (Domain layer implemented)
- ✅ Domain entities and repositories defined
- ✅ Domain services and errors implemented

## 🎯 Goals

1. **Create Use Cases**: Implement business use cases that orchestrate domain logic
2. **Define DTOs**: Create data transfer objects for clean data boundaries
3. **Establish Ports**: Define primary and secondary ports for external communication
4. **Input Validation**: Implement validation for use case inputs
5. **Error Handling**: Proper error management at application layer

## 📁 Directory Structure

Create the following structure in `internal/application/`:

```text
internal/application/
├── usecases/              # Business use cases
│   ├── setup_repositories.go
│   ├── validate_prerequisites.go
│   └── manage_docker_images.go
├── dtos/                  # Data Transfer Objects
│   ├── requests.go        # Input DTOs
│   └── responses.go       # Output DTOs
├── ports/                 # Application ports
│   ├── primary.go         # Primary ports (driving adapters)
│   └── secondary.go       # Secondary ports (driven adapters)
└── services/              # Application services
    └── orchestration.go   # Cross-cutting orchestration logic
```

## 🚀 Implementation Steps

### Step 1: Create Data Transfer Objects (DTOs)

First, create the DTOs that will be used for data transfer between layers.

#### Create `internal/application/dtos/requests.go`

```go
package dtos

// SetupRepositoriesRequest represents the input for repository setup
type SetupRepositoriesRequest struct {
    Repositories []RepositorySetupData `json:"repositories" validate:"required,dive"`
    ForceClone   bool                  `json:"force_clone"`
    WorkingDir   string                `json:"working_dir" validate:"required"`
}

// RepositorySetupData represents data for setting up a repository
type RepositorySetupData struct {
    Name   string `json:"name" validate:"required,min=1,max=255"`
    URL    string `json:"url" validate:"required,url"`
    Branch string `json:"branch" validate:"required,min=1"`
}

// ValidatePrerequisitesRequest represents input for prerequisite validation
type ValidatePrerequisitesRequest struct {
    Commands      []string `json:"commands"`
    CheckVersions bool     `json:"check_versions"`
    FailFast      bool     `json:"fail_fast"`
}

// ManageDockerImagesRequest represents input for Docker image management
type ManageDockerImagesRequest struct {
    Action string                `json:"action" validate:"required,oneof=build list remove"`
    Images []DockerImageRequest  `json:"images,omitempty"`
}

// DockerImageRequest represents Docker image data
type DockerImageRequest struct {
    Name       string            `json:"name" validate:"required"`
    Tag        string            `json:"tag"`
    Dockerfile string            `json:"dockerfile"`
    Context    string            `json:"context"`
    Args       map[string]string `json:"args"`
}
```

#### Create `internal/application/dtos/responses.go`

```go
package dtos

import "time"

// SetupRepositoriesResponse represents the output of repository setup
type SetupRepositoriesResponse struct {
    SetupResults []RepositorySetupResult `json:"setup_results"`
    TotalCount   int                     `json:"total_count"`
    SuccessCount int                     `json:"success_count"`
    FailureCount int                     `json:"failure_count"`
    Duration     time.Duration           `json:"duration"`
}

// RepositorySetupResult represents the result of setting up a single repository
type RepositorySetupResult struct {
    Name      string `json:"name"`
    Status    string `json:"status"` // "success", "failed", "skipped"
    Message   string `json:"message"`
    LocalPath string `json:"local_path,omitempty"`
    Error     string `json:"error,omitempty"`
}

// ValidatePrerequisitesResponse represents output of prerequisite validation
type ValidatePrerequisitesResponse struct {
    Results    []ValidationResult `json:"results"`
    AllPassed  bool               `json:"all_passed"`
    Summary    ValidationSummary  `json:"summary"`
}

// ValidationResult represents a single validation check result
type ValidationResult struct {
    Check     string `json:"check"`
    Status    string `json:"status"` // "pass", "fail", "warning"
    Message   string `json:"message"`
    Details   string `json:"details,omitempty"`
    Required  bool   `json:"required"`
}

// ValidationSummary provides summary statistics
type ValidationSummary struct {
    Total    int `json:"total"`
    Passed   int `json:"passed"`
    Failed   int `json:"failed"`
    Warnings int `json:"warnings"`
}

// ManageDockerImagesResponse represents output of Docker image management
type ManageDockerImagesResponse struct {
    Action  string              `json:"action"`
    Results []DockerImageResult `json:"results"`
    Summary ImageSummary        `json:"summary"`
}

// DockerImageResult represents result of Docker image operation
type DockerImageResult struct {
    Name    string `json:"name"`
    Status  string `json:"status"`
    Message string `json:"message"`
    ImageID string `json:"image_id,omitempty"`
    Size    string `json:"size,omitempty"`
}

// ImageSummary provides Docker operation summary
type ImageSummary struct {
    Total     int `json:"total"`
    Succeeded int `json:"succeeded"`
    Failed    int `json:"failed"`
}
```

### Step 2: Define Application Ports

#### Create `internal/application/ports/primary.go`

```go
package ports

import (
    "context"
    
    "github.com/fabianoflorentino/whiterose/internal/application/dtos"
)

// RepositoryManagementPort defines the primary port for repository operations
type RepositoryManagementPort interface {
    SetupRepositories(ctx context.Context, request dtos.SetupRepositoriesRequest) (*dtos.SetupRepositoriesResponse, error)
    GetRepositoryStatus(ctx context.Context) ([]dtos.RepositoryStatus, error)
    UpdateRepository(ctx context.Context, name, newBranch string) error
    RemoveRepository(ctx context.Context, name string) error
}

// SystemValidationPort defines the primary port for system validation
type SystemValidationPort interface {
    ValidatePrerequisites(ctx context.Context, request dtos.ValidatePrerequisitesRequest) (*dtos.ValidatePrerequisitesResponse, error)
    GetSystemInfo(ctx context.Context) (*dtos.SystemInfoResponse, error)
    CheckCommand(ctx context.Context, command string) (*dtos.CommandCheckResponse, error)
}

// DockerManagementPort defines the primary port for Docker operations
type DockerManagementPort interface {
    ManageImages(ctx context.Context, request dtos.ManageDockerImagesRequest) (*dtos.ManageDockerImagesResponse, error)
    BuildImage(ctx context.Context, request dtos.DockerImageRequest) (*dtos.DockerImageResult, error)
    ListImages(ctx context.Context) (*dtos.ManageDockerImagesResponse, error)
}
```

#### Create `internal/application/ports/secondary.go`

```go
package ports

import (
    "context"
    
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
)

// GitRepositoryPort defines the secondary port for Git operations
type GitRepositoryPort interface {
    Clone(ctx context.Context, repo *entities.Repository, localPath string) error
    Pull(ctx context.Context, localPath string) error
    Checkout(ctx context.Context, localPath, branch string) error
    GetCurrentBranch(ctx context.Context, localPath string) (string, error)
    IsClean(ctx context.Context, localPath string) (bool, error)
}

// ConfigurationPort defines the secondary port for configuration management
type ConfigurationPort interface {
    LoadConfig(ctx context.Context) (*Configuration, error)
    SaveConfig(ctx context.Context, config *Configuration) error
    GetRepositories(ctx context.Context) ([]*RepositoryConfig, error)
}

// ValidationPort defines the secondary port for system validation
type ValidationPort interface {
    CheckCommand(ctx context.Context, command string) error
    CheckVersion(ctx context.Context, command, minVersion string) error
    ValidateEnvironment(ctx context.Context) ([]ValidationResult, error)
}

// DockerPort defines the secondary port for Docker operations
type DockerPort interface {
    BuildImage(ctx context.Context, imageName, dockerfile, context string, args map[string]string) (string, error)
    ListImages(ctx context.Context) ([]DockerImage, error)
    RemoveImage(ctx context.Context, imageID string) error
}

// NotificationPort defines the secondary port for notifications
type NotificationPort interface {
    SendNotification(ctx context.Context, level, message string) error
    LogOperation(ctx context.Context, operation, details string) error
}
```

### Step 3: Implement Use Cases

#### Create `internal/application/usecases/setup_repositories.go`

```go
package usecases

import (
    "context"
    "fmt"
    "path/filepath"
    "time"

    "github.com/fabianoflorentino/whiterose/internal/application/dtos"
    "github.com/fabianoflorentino/whiterose/internal/application/ports"
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
    "github.com/fabianoflorentino/whiterose/internal/domain/errors"
    "github.com/fabianoflorentino/whiterose/internal/domain/repositories"
)

// SetupRepositoriesUseCase handles repository setup operations
type SetupRepositoriesUseCase struct {
    repositoryRepo repositories.RepositoryRepository
    gitPort        ports.GitRepositoryPort
    configPort     ports.ConfigurationPort
    notifyPort     ports.NotificationPort
}

// NewSetupRepositoriesUseCase creates a new use case instance
func NewSetupRepositoriesUseCase(
    repositoryRepo repositories.RepositoryRepository,
    gitPort ports.GitRepositoryPort,
    configPort ports.ConfigurationPort,
    notifyPort ports.NotificationPort,
) *SetupRepositoriesUseCase {
    return &SetupRepositoriesUseCase{
        repositoryRepo: repositoryRepo,
        gitPort:        gitPort,
        configPort:     configPort,
        notifyPort:     notifyPort,
    }
}

// Execute performs the repository setup operation
func (uc *SetupRepositoriesUseCase) Execute(ctx context.Context, request dtos.SetupRepositoriesRequest) (*dtos.SetupRepositoriesResponse, error) {
    startTime := time.Now()
    
    // Validate request
    if err := uc.validateRequest(request); err != nil {
        return nil, err
    }

    // Log operation start
    uc.notifyPort.LogOperation(ctx, "setup_repositories", fmt.Sprintf("Starting setup of %d repositories", len(request.Repositories)))

    response := &dtos.SetupRepositoriesResponse{
        SetupResults: make([]dtos.RepositorySetupResult, 0, len(request.Repositories)),
        TotalCount:   len(request.Repositories),
    }

    // Process each repository
    for _, repoData := range request.Repositories {
        result := uc.setupSingleRepository(ctx, repoData, request.ForceClone, request.WorkingDir)
        response.SetupResults = append(response.SetupResults, result)

        if result.Status == "success" {
            response.SuccessCount++
        } else {
            response.FailureCount++
        }
    }

    response.Duration = time.Since(startTime)
    
    // Send completion notification
    uc.notifyPort.SendNotification(ctx, "info", 
        fmt.Sprintf("Repository setup completed: %d success, %d failed", 
            response.SuccessCount, response.FailureCount))

    return response, nil
}

// validateRequest validates the input request
func (uc *SetupRepositoriesUseCase) validateRequest(request dtos.SetupRepositoriesRequest) error {
    if len(request.Repositories) == 0 {
        return errors.NewValidationError("at least one repository must be specified", nil)
    }

    if request.WorkingDir == "" {
        return errors.NewValidationError("working directory must be specified", nil)
    }

    // Check for duplicate names
    names := make(map[string]bool)
    for _, repo := range request.Repositories {
        if names[repo.Name] {
            return errors.NewValidationError(fmt.Sprintf("duplicate repository name: %s", repo.Name), nil)
        }
        names[repo.Name] = true
    }

    return nil
}

// setupSingleRepository sets up a single repository
func (uc *SetupRepositoriesUseCase) setupSingleRepository(ctx context.Context, repoData dtos.RepositorySetupData, forceClone bool, workingDir string) dtos.RepositorySetupResult {
    // Create repository entity
    repo, err := entities.NewRepository(repoData.Name, repoData.URL, repoData.Branch)
    if err != nil {
        return dtos.RepositorySetupResult{
            Name:    repoData.Name,
            Status:  "failed",
            Message: "Failed to create repository entity",
            Error:   err.Error(),
        }
    }

    // Check if repository already exists
    existingRepo, err := uc.repositoryRepo.FindByName(ctx, repoData.Name)
    if err != nil && !errors.IsNotFoundError(err) {
        return dtos.RepositorySetupResult{
            Name:    repoData.Name,
            Status:  "failed",
            Message: "Failed to check existing repository",
            Error:   err.Error(),
        }
    }

    // Skip if repository exists and not forcing clone
    if existingRepo != nil && !forceClone {
        return dtos.RepositorySetupResult{
            Name:      repoData.Name,
            Status:    "skipped",
            Message:   "Repository already exists (use force_clone to override)",
            LocalPath: existingRepo.LocalPath(),
        }
    }

    // Set local path
    localPath := filepath.Join(workingDir, repoData.Name)
    if err := repo.SetLocalPath(localPath); err != nil {
        return dtos.RepositorySetupResult{
            Name:    repoData.Name,
            Status:  "failed",
            Message: "Failed to set local path",
            Error:   err.Error(),
        }
    }

    // Clone repository
    if err := uc.gitPort.Clone(ctx, repo, localPath); err != nil {
        return dtos.RepositorySetupResult{
            Name:    repoData.Name,
            Status:  "failed",
            Message: "Failed to clone repository",
            Error:   err.Error(),
        }
    }

    // Mark as cloned
    repo.MarkAsCloned()

    // Save repository
    if err := uc.repositoryRepo.Save(ctx, repo); err != nil {
        return dtos.RepositorySetupResult{
            Name:    repoData.Name,
            Status:  "failed",
            Message: "Failed to save repository",
            Error:   err.Error(),
        }
    }

    return dtos.RepositorySetupResult{
        Name:      repoData.Name,
        Status:    "success",
        Message:   "Repository successfully cloned and configured",
        LocalPath: localPath,
    }
}
```

### Step 4: Create Application Services

#### Create `internal/application/services/orchestration.go`

```go
package services

import (
    "context"
    "fmt"

    "github.com/fabianoflorentino/whiterose/internal/application/ports"
    "github.com/fabianoflorentino/whiterose/internal/domain/repositories"
)

// OrchestrationService coordinates operations across multiple use cases
type OrchestrationService struct {
    repositoryRepo repositories.RepositoryRepository
    configPort     ports.ConfigurationPort
    validationPort ports.ValidationPort
    notifyPort     ports.NotificationPort
}

// NewOrchestrationService creates a new orchestration service
func NewOrchestrationService(
    repositoryRepo repositories.RepositoryRepository,
    configPort ports.ConfigurationPort,
    validationPort ports.ValidationPort,
    notifyPort ports.NotificationPort,
) *OrchestrationService {
    return &OrchestrationService{
        repositoryRepo: repositoryRepo,
        configPort:     configPort,
        validationPort: validationPort,
        notifyPort:     notifyPort,
    }
}

// ValidateAndSetup performs validation before repository setup
func (s *OrchestrationService) ValidateAndSetup(ctx context.Context) error {
    // First validate environment
    results, err := s.validationPort.ValidateEnvironment(ctx)
    if err != nil {
        return fmt.Errorf("environment validation failed: %w", err)
    }

    // Check if any critical validations failed
    for _, result := range results {
        if result.Status == "fail" && result.Required {
            s.notifyPort.SendNotification(ctx, "error", 
                fmt.Sprintf("Critical validation failed: %s", result.Message))
            return fmt.Errorf("critical validation failed: %s", result.Message)
        }
    }

    s.notifyPort.SendNotification(ctx, "info", "Environment validation passed")
    return nil
}

// CleanupOrphanedRepositories removes repositories that are no longer configured
func (s *OrchestrationService) CleanupOrphanedRepositories(ctx context.Context) error {
    // Get configured repositories
    configRepos, err := s.configPort.GetRepositories(ctx)
    if err != nil {
        return fmt.Errorf("failed to get configured repositories: %w", err)
    }

    // Get stored repositories
    storedRepos, err := s.repositoryRepo.FindAll(ctx)
    if err != nil {
        return fmt.Errorf("failed to get stored repositories: %w", err)
    }

    // Create map of configured repository names
    configuredNames := make(map[string]bool)
    for _, repo := range configRepos {
        configuredNames[repo.Name] = true
    }

    // Remove orphaned repositories
    for _, stored := range storedRepos {
        if !configuredNames[stored.Name()] {
            if err := s.repositoryRepo.Delete(ctx, stored.ID()); err != nil {
                s.notifyPort.SendNotification(ctx, "warning", 
                    fmt.Sprintf("Failed to remove orphaned repository %s: %v", stored.Name(), err))
            } else {
                s.notifyPort.SendNotification(ctx, "info", 
                    fmt.Sprintf("Removed orphaned repository: %s", stored.Name()))
            }
        }
    }

    return nil
}
```

## ✅ Acceptance Criteria

- [ ] All use cases implemented and working
- [ ] DTOs properly structured and validated
- [ ] Primary and secondary ports defined
- [ ] Input validation implemented
- [ ] Error handling at application layer
- [ ] Unit tests with >90% coverage
- [ ] Integration tests for use case flows
- [ ] Documentation updated

## 🧪 Testing Strategy

### Unit Tests

```go
func TestSetupRepositoriesUseCase_Execute(t *testing.T) {
    // Mock dependencies
    mockRepoRepo := &mocks.RepositoryRepository{}
    mockGitPort := &mocks.GitRepositoryPort{}
    mockConfigPort := &mocks.ConfigurationPort{}
    mockNotifyPort := &mocks.NotificationPort{}

    // Create use case
    uc := NewSetupRepositoriesUseCase(mockRepoRepo, mockGitPort, mockConfigPort, mockNotifyPort)

    // Test successful execution
    request := dtos.SetupRepositoriesRequest{
        Repositories: []dtos.RepositorySetupData{
            {Name: "test", URL: "https://github.com/test/test.git", Branch: "main"},
        },
        ForceClone: false,
        WorkingDir: "/tmp",
    }

    // Setup mocks
    mockRepoRepo.On("FindByName", mock.Anything, "test").Return(nil, errors.NewNotFoundError("not found"))
    mockGitPort.On("Clone", mock.Anything, mock.Anything, mock.Anything).Return(nil)
    mockRepoRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    mockNotifyPort.On("LogOperation", mock.Anything, mock.Anything, mock.Anything).Return(nil)
    mockNotifyPort.On("SendNotification", mock.Anything, mock.Anything, mock.Anything).Return(nil)

    // Execute
    response, err := uc.Execute(context.Background(), request)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 1, response.SuccessCount)
    assert.Equal(t, 0, response.FailureCount)
}
```

## 📚 Next Steps

After completing Phase 2:

1. ✅ All use cases implemented and tested
2. ✅ Application layer provides clean interfaces
3. ✅ Ready to implement infrastructure adapters
4. 🚀 Proceed to [Phase 3: Infrastructure Layer](phase-3-infrastructure.md)

## 🔗 Related Documentation

- [Domain Layer Guide](phase-1-domain.md)
- [Code Examples - Application Layer](../code-examples/application/)
- [Architecture Overview](../architecture/current-vs-proposed.md)
- [Testing Guide](../how-to/testing-guide.md)
