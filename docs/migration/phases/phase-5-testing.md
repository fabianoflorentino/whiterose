# Phase 5: Testing & Documentation

## 🎯 Objective

Implement comprehensive testing strategy and complete documentation for the hexagonal architecture migration, ensuring production readiness and maintainability.

## ⏱️ Duration: 3-4 days | Effort: 24-32h

## 📋 Prerequisites

- ✅ Phase 4 completed (Migration & Integration implemented)
- ✅ All layers integrated and functional
- ✅ CLI commands migrated and working
- ✅ Configuration migration tested

## 🎯 Goals

1. **Comprehensive Testing**: Unit, integration, and end-to-end tests
2. **Test Coverage**: Achieve >90% code coverage
3. **Performance Testing**: Load and stress testing
4. **Documentation**: Complete API documentation and guides
5. **CI/CD Pipeline**: Automated testing and deployment
6. **Monitoring**: Health checks and observability
7. **Production Readiness**: Security, logging, and error handling

## 📁 Directory Structure

Create comprehensive testing and documentation structure:

```text
test/
├── unit/                     # Unit tests
│   ├── domain/
│   │   ├── entities_test.go
│   │   └── repositories_test.go
│   ├── application/
│   │   ├── usecases_test.go
│   │   └── services_test.go
│   └── infrastructure/
│       ├── adapters_test.go
│       └── config_test.go
├── integration/              # Integration tests
│   ├── git_integration_test.go
│   ├── docker_integration_test.go
│   └── cli_integration_test.go
├── e2e/                     # End-to-end tests
│   ├── scenarios/
│   │   ├── setup_scenario_test.go
│   │   ├── migration_scenario_test.go
│   │   └── workflow_scenario_test.go
│   └── fixtures/
│       ├── configs/
│       └── repositories/
├── performance/             # Performance tests
│   ├── load_test.go
│   ├── stress_test.go
│   └── benchmarks/
├── testdata/               # Test data and fixtures
│   ├── configs/
│   ├── repositories/
│   └── docker/
└── helpers/                # Test utilities
    ├── mocks/
    ├── fixtures/
    └── assertions/

docs/
├── api/                    # API documentation
│   ├── openapi.yaml
│   └── endpoints/
├── guides/                 # User guides
│   ├── quick-start.md
│   ├── configuration.md
│   ├── troubleshooting.md
│   └── best-practices.md
├── architecture/           # Architecture documentation
│   ├── overview.md
│   ├── decisions/          # ADRs
│   └── diagrams/
└── development/           # Development documentation
    ├── contributing.md
    ├── testing.md
    └── deployment.md

.github/
├── workflows/             # GitHub Actions
│   ├── ci.yml
│   ├── release.yml
│   └── security.yml
└── ISSUE_TEMPLATE/
```

## 🚀 Implementation Steps

### Step 1: Comprehensive Unit Tests

#### Create `test/unit/domain/entities_test.go`

```go
package domain

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
    "github.com/fabianoflorentino/whiterose/internal/domain/errors"
)

func TestRepository_Creation(t *testing.T) {
    tests := []struct {
        name        string
        repoName    string
        repoURL     string
        branch      string
        expectError bool
        errorType   error
    }{
        {
            name:        "valid repository",
            repoName:    "test-repo",
            repoURL:     "https://github.com/user/repo.git",
            branch:      "main",
            expectError: false,
        },
        {
            name:        "empty name",
            repoName:    "",
            repoURL:     "https://github.com/user/repo.git",
            branch:      "main",
            expectError: true,
            errorType:   &errors.ValidationError{},
        },
        {
            name:        "invalid URL",
            repoName:    "test-repo",
            repoURL:     "invalid-url",
            branch:      "main",
            expectError: true,
            errorType:   &errors.ValidationError{},
        },
        {
            name:        "empty branch",
            repoName:    "test-repo",
            repoURL:     "https://github.com/user/repo.git",
            branch:      "",
            expectError: true,
            errorType:   &errors.ValidationError{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo, err := entities.NewRepository(tt.repoName, tt.repoURL, tt.branch)

            if tt.expectError {
                assert.Error(t, err)
                assert.IsType(t, tt.errorType, err)
                assert.Nil(t, repo)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, repo)
                assert.Equal(t, tt.repoName, repo.Name())
                assert.Equal(t, tt.repoURL, repo.URL().String())
                assert.Equal(t, tt.branch, repo.Branch())
            }
        })
    }
}

func TestRepository_Clone(t *testing.T) {
    repo, err := entities.NewRepository("test-repo", "https://github.com/user/repo.git", "main")
    require.NoError(t, err)

    localPath := "/tmp/test-repo"

    err = repo.Clone(localPath)
    assert.NoError(t, err)
    assert.Equal(t, localPath, repo.LocalPath())
    assert.True(t, repo.IsCloned())
}

func TestRepository_Tags(t *testing.T) {
    repo, err := entities.NewRepository("test-repo", "https://github.com/user/repo.git", "main")
    require.NoError(t, err)

    // Test adding tags
    repo.AddTag("frontend")
    repo.AddTag("javascript")

    assert.True(t, repo.HasTag("frontend"))
    assert.True(t, repo.HasTag("javascript"))
    assert.False(t, repo.HasTag("backend"))

    // Test removing tags
    repo.RemoveTag("frontend")
    assert.False(t, repo.HasTag("frontend"))
    assert.True(t, repo.HasTag("javascript"))

    // Test getting all tags
    tags := repo.Tags()
    assert.Len(t, tags, 1)
    assert.Contains(t, tags, "javascript")
}

func TestRepository_Equality(t *testing.T) {
    repo1, err := entities.NewRepository("test-repo", "https://github.com/user/repo.git", "main")
    require.NoError(t, err)

    repo2, err := entities.NewRepository("test-repo", "https://github.com/user/repo.git", "main")
    require.NoError(t, err)

    repo3, err := entities.NewRepository("different-repo", "https://github.com/user/repo.git", "main")
    require.NoError(t, err)

    assert.True(t, repo1.Equals(repo2))
    assert.False(t, repo1.Equals(repo3))
}

func TestRepository_Validation(t *testing.T) {
    repo, err := entities.NewRepository("test-repo", "https://github.com/user/repo.git", "main")
    require.NoError(t, err)

    // Test valid repository
    err = repo.Validate()
    assert.NoError(t, err)

    // Test repository with invalid state (this would require modifying internal state)
    // In a real implementation, you might have methods that could put the repository in an invalid state
}

// Benchmark tests
func BenchmarkRepository_Creation(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = entities.NewRepository("test-repo", "https://github.com/user/repo.git", "main")
    }
}

func BenchmarkRepository_TagOperations(b *testing.B) {
    repo, _ := entities.NewRepository("test-repo", "https://github.com/user/repo.git", "main")
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        repo.AddTag("tag1")
        repo.AddTag("tag2")
        repo.HasTag("tag1")
        repo.RemoveTag("tag1")
    }
}
```

#### Create `test/unit/application/usecases_test.go`

```go
package application

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"

    "github.com/fabianoflorentino/whiterose/internal/application/dtos"
    "github.com/fabianoflorentino/whiterose/internal/application/usecases"
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
    "github.com/fabianoflorentino/whiterose/test/helpers/mocks"
)

func TestSetupRepositoriesUseCase_Execute(t *testing.T) {
    // Setup mocks
    mockRepo := mocks.NewMockRepositoryRepository(t)
    mockGit := mocks.NewMockGitRepositoryPort(t)
    mockConfig := mocks.NewMockConfigurationPort(t)
    mockNotification := mocks.NewMockNotificationPort(t)

    // Create use case
    useCase := usecases.NewSetupRepositoriesUseCase(
        mockRepo,
        mockGit,
        mockConfig,
        mockNotification,
    )

    t.Run("successful setup", func(t *testing.T) {
        // Arrange
        ctx := context.Background()
        request := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: "/tmp/repos",
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "test-repo",
                    URL:    "https://github.com/user/repo.git",
                    Branch: "main",
                },
            },
            Force:  false,
            DryRun: false,
        }

        // Setup expectations
        repo, _ := entities.NewRepository("test-repo", "https://github.com/user/repo.git", "main")
        
        mockRepo.On("FindByName", ctx, "test-repo").Return(nil, nil)
        mockRepo.On("Save", ctx, mock.AnythingOfType("*entities.Repository")).Return(nil)
        mockGit.On("Clone", ctx, repo, "/tmp/repos/test-repo").Return(nil)
        mockNotification.On("Notify", ctx, mock.AnythingOfType("string")).Return(nil)

        // Act
        response, err := useCase.Execute(ctx, request)

        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, response)
        assert.Equal(t, 1, response.TotalRepositories)
        assert.Equal(t, 1, response.SuccessfulSetups)
        assert.Equal(t, 0, response.FailedSetups)
        assert.Equal(t, 0, response.SkippedRepositories)

        // Verify all expectations were met
        mockRepo.AssertExpectations(t)
        mockGit.AssertExpectations(t)
        mockConfig.AssertExpectations(t)
        mockNotification.AssertExpectations(t)
    })

    t.Run("repository already exists", func(t *testing.T) {
        // Arrange
        ctx := context.Background()
        request := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: "/tmp/repos",
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "existing-repo",
                    URL:    "https://github.com/user/repo.git",
                    Branch: "main",
                },
            },
            Force:  false,
            DryRun: false,
        }

        existingRepo, _ := entities.NewRepository("existing-repo", "https://github.com/user/repo.git", "main")
        
        mockRepo.On("FindByName", ctx, "existing-repo").Return(existingRepo, nil)
        mockNotification.On("Notify", ctx, mock.AnythingOfType("string")).Return(nil)

        // Act
        response, err := useCase.Execute(ctx, request)

        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, response)
        assert.Equal(t, 1, response.TotalRepositories)
        assert.Equal(t, 0, response.SuccessfulSetups)
        assert.Equal(t, 0, response.FailedSetups)
        assert.Equal(t, 1, response.SkippedRepositories)

        mockRepo.AssertExpectations(t)
        mockNotification.AssertExpectations(t)
    })

    t.Run("git clone failure", func(t *testing.T) {
        // Arrange
        ctx := context.Background()
        request := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: "/tmp/repos",
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "fail-repo",
                    URL:    "https://github.com/user/nonexistent.git",
                    Branch: "main",
                },
            },
            Force:  false,
            DryRun: false,
        }

        repo, _ := entities.NewRepository("fail-repo", "https://github.com/user/nonexistent.git", "main")
        
        mockRepo.On("FindByName", ctx, "fail-repo").Return(nil, nil)
        mockGit.On("Clone", ctx, repo, "/tmp/repos/fail-repo").Return(errors.New("clone failed"))
        mockNotification.On("Notify", ctx, mock.AnythingOfType("string")).Return(nil)

        // Act
        response, err := useCase.Execute(ctx, request)

        // Assert
        assert.NoError(t, err) // Use case doesn't fail, but tracks errors
        assert.NotNil(t, response)
        assert.Equal(t, 1, response.TotalRepositories)
        assert.Equal(t, 0, response.SuccessfulSetups)
        assert.Equal(t, 1, response.FailedSetups)
        assert.Len(t, response.Errors, 1)

        mockRepo.AssertExpectations(t)
        mockGit.AssertExpectations(t)
        mockNotification.AssertExpectations(t)
    })

    t.Run("dry run mode", func(t *testing.T) {
        // Arrange
        ctx := context.Background()
        request := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: "/tmp/repos",
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "dry-run-repo",
                    URL:    "https://github.com/user/repo.git",
                    Branch: "main",
                },
            },
            Force:  false,
            DryRun: true,
        }

        mockRepo.On("FindByName", ctx, "dry-run-repo").Return(nil, nil)
        mockNotification.On("Notify", ctx, mock.AnythingOfType("string")).Return(nil)

        // Act
        response, err := useCase.Execute(ctx, request)

        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, response)
        assert.Equal(t, 1, response.TotalRepositories)
        
        // In dry run, nothing should be saved or cloned
        mockRepo.AssertNotCalled(t, "Save")
        mockGit.AssertNotCalled(t, "Clone")
        mockNotification.AssertExpectations(t)
    })
}

// Table-driven test for validation
func TestSetupRepositoriesUseCase_Validation(t *testing.T) {
    tests := []struct {
        name        string
        request     *dtos.SetupRepositoriesRequest
        expectError bool
        errorMsg    string
    }{
        {
            name: "valid request",
            request: &dtos.SetupRepositoriesRequest{
                WorkingDirectory: "/tmp/repos",
                Repositories: []*dtos.RepositorySetupInfo{
                    {Name: "test", URL: "https://github.com/user/repo.git", Branch: "main"},
                },
            },
            expectError: false,
        },
        {
            name: "empty working directory",
            request: &dtos.SetupRepositoriesRequest{
                WorkingDirectory: "",
                Repositories: []*dtos.RepositorySetupInfo{
                    {Name: "test", URL: "https://github.com/user/repo.git", Branch: "main"},
                },
            },
            expectError: true,
            errorMsg:    "working directory cannot be empty",
        },
        {
            name: "no repositories",
            request: &dtos.SetupRepositoriesRequest{
                WorkingDirectory: "/tmp/repos",
                Repositories:     []*dtos.RepositorySetupInfo{},
            },
            expectError: true,
            errorMsg:    "at least one repository must be specified",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            mockRepo := mocks.NewMockRepositoryRepository(t)
            mockGit := mocks.NewMockGitRepositoryPort(t)
            mockConfig := mocks.NewMockConfigurationPort(t)
            mockNotification := mocks.NewMockNotificationPort(t)

            useCase := usecases.NewSetupRepositoriesUseCase(
                mockRepo, mockGit, mockConfig, mockNotification,
            )

            _, err := useCase.Execute(context.Background(), tt.request)

            if tt.expectError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errorMsg)
            } else {
                // For valid requests, we need to setup mock expectations
                if !tt.expectError {
                    // Add minimal expectations for valid test
                    mockRepo.On("FindByName", mock.Anything, mock.Anything).Return(nil, nil)
                    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
                    mockGit.On("Clone", mock.Anything, mock.Anything, mock.Anything).Return(nil)
                    mockNotification.On("Notify", mock.Anything, mock.Anything).Return(nil)
                }
                assert.NoError(t, err)
            }
        })
    }
}

// Performance test
func BenchmarkSetupRepositoriesUseCase_Execute(b *testing.B) {
    // Setup mocks with relaxed expectations
    mockRepo := mocks.NewMockRepositoryRepository(b)
    mockGit := mocks.NewMockGitRepositoryPort(b)
    mockConfig := mocks.NewMockConfigurationPort(b)
    mockNotification := mocks.NewMockNotificationPort(b)

    mockRepo.On("FindByName", mock.Anything, mock.Anything).Return(nil, nil)
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    mockGit.On("Clone", mock.Anything, mock.Anything, mock.Anything).Return(nil)
    mockNotification.On("Notify", mock.Anything, mock.Anything).Return(nil)

    useCase := usecases.NewSetupRepositoriesUseCase(
        mockRepo, mockGit, mockConfig, mockNotification,
    )

    request := &dtos.SetupRepositoriesRequest{
        WorkingDirectory: "/tmp/repos",
        Repositories: []*dtos.RepositorySetupInfo{
            {Name: "test", URL: "https://github.com/user/repo.git", Branch: "main"},
        },
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = useCase.Execute(context.Background(), request)
    }
}
```

### Step 2: Integration Tests

#### Create `test/integration/git_integration_test.go`

```go
package integration

import (
    "context"
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
    gitAdapter "github.com/fabianoflorentino/whiterose/internal/infrastructure/adapters/git"
)

func TestGitAdapter_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup
    adapter := gitAdapter.NewGitAdapter(gitAdapter.GitConfig{
        Timeout: 30 * time.Second,
    })

    t.Run("clone public repository", func(t *testing.T) {
        repo, err := entities.NewRepository(
            "git-test",
            "https://github.com/octocat/Hello-World.git",
            "master",
        )
        require.NoError(t, err)

        tmpDir := t.TempDir()
        localPath := filepath.Join(tmpDir, "hello-world")

        // Test clone
        err = adapter.Clone(context.Background(), repo, localPath)
        assert.NoError(t, err)

        // Verify repository structure
        assert.DirExists(t, localPath)
        assert.DirExists(t, filepath.Join(localPath, ".git"))
        assert.FileExists(t, filepath.Join(localPath, "README"))

        // Test get current branch
        branch, err := adapter.GetCurrentBranch(context.Background(), localPath)
        assert.NoError(t, err)
        assert.Equal(t, "master", branch)

        // Test repository is clean
        isClean, err := adapter.IsClean(context.Background(), localPath)
        assert.NoError(t, err)
        assert.True(t, isClean)
    })

    t.Run("clone non-existent repository", func(t *testing.T) {
        repo, err := entities.NewRepository(
            "non-existent",
            "https://github.com/nonexistent/nonexistent.git",
            "main",
        )
        require.NoError(t, err)

        tmpDir := t.TempDir()
        localPath := filepath.Join(tmpDir, "non-existent")

        // Test clone failure
        err = adapter.Clone(context.Background(), repo, localPath)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "failed to clone repository")

        // Verify no directory was created
        assert.NoDirExists(t, localPath)
    })

    t.Run("clone with timeout", func(t *testing.T) {
        // Create adapter with very short timeout
        shortTimeoutAdapter := gitAdapter.NewGitAdapter(gitAdapter.GitConfig{
            Timeout: 1 * time.Millisecond,
        })

        repo, err := entities.NewRepository(
            "timeout-test",
            "https://github.com/octocat/Hello-World.git",
            "master",
        )
        require.NoError(t, err)

        tmpDir := t.TempDir()
        localPath := filepath.Join(tmpDir, "timeout-test")

        // Test timeout
        err = shortTimeoutAdapter.Clone(context.Background(), repo, localPath)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "context deadline exceeded")
    })
}

func TestGitAdapter_PullAndCheckout(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    adapter := gitAdapter.NewGitAdapter(gitAdapter.GitConfig{
        Timeout: 30 * time.Second,
    })

    // Clone repository first
    repo, err := entities.NewRepository(
        "pull-test",
        "https://github.com/octocat/Hello-World.git",
        "master",
    )
    require.NoError(t, err)

    tmpDir := t.TempDir()
    localPath := filepath.Join(tmpDir, "pull-test")

    err = adapter.Clone(context.Background(), repo, localPath)
    require.NoError(t, err)

    t.Run("pull latest changes", func(t *testing.T) {
        err := adapter.Pull(context.Background(), localPath)
        assert.NoError(t, err)
    })

    t.Run("checkout existing branch", func(t *testing.T) {
        err := adapter.Checkout(context.Background(), localPath, "master")
        assert.NoError(t, err)

        branch, err := adapter.GetCurrentBranch(context.Background(), localPath)
        assert.NoError(t, err)
        assert.Equal(t, "master", branch)
    })
}

// Parallel integration tests
func TestGitAdapter_Parallel(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    adapter := gitAdapter.NewGitAdapter(gitAdapter.GitConfig{
        Timeout: 30 * time.Second,
    })

    repos := []struct {
        name string
        url  string
    }{
        {"repo1", "https://github.com/octocat/Hello-World.git"},
        {"repo2", "https://github.com/octocat/Hello-World.git"},
        {"repo3", "https://github.com/octocat/Hello-World.git"},
    }

    t.Run("parallel clones", func(t *testing.T) {
        for i, repoInfo := range repos {
            i, repoInfo := i, repoInfo // Capture loop variables
            t.Run(repoInfo.name, func(t *testing.T) {
                t.Parallel()

                repo, err := entities.NewRepository(
                    repoInfo.name,
                    repoInfo.url,
                    "master",
                )
                require.NoError(t, err)

                tmpDir := t.TempDir()
                localPath := filepath.Join(tmpDir, repoInfo.name)

                err = adapter.Clone(context.Background(), repo, localPath)
                assert.NoError(t, err)
                assert.DirExists(t, localPath)
            })
        }
    })
}
```

### Step 3: End-to-End Tests

#### Create `test/e2e/scenarios/setup_scenario_test.go`

```go
package scenarios

import (
    "context"
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/fabianoflorentino/whiterose/internal/application/dtos"
    "github.com/fabianoflorentino/whiterose/internal/infrastructure/config"
)

func TestCompleteSetupScenario(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    // Test setup
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, "config.json")
    workingDir := filepath.Join(tmpDir, "repositories")

    // Initialize container
    appConfig := &config.AppConfig{
        ConfigPath:          configPath,
        WorkingDir:          workingDir,
        DatabasePath:        filepath.Join(tmpDir, "test.db"),
        GitTimeout:          30 * time.Second,
        DockerTimeout:       60 * time.Second,
        LogLevel:            "debug",
        EnableNotifications: true,
    }

    container, err := config.NewContainer(appConfig)
    require.NoError(t, err)
    defer container.Close()

    ctx := context.Background()

    t.Run("complete setup workflow", func(t *testing.T) {
        // Step 1: Validate Prerequisites
        prereqUC := container.GetValidatePrerequisitesUseCase()
        prereqRequest := &dtos.ValidatePrerequisitesRequest{
            SkipOptional: true, // Skip optional for CI environments
        }

        prereqResponse, err := prereqUC.Execute(ctx, prereqRequest)
        require.NoError(t, err)
        assert.True(t, prereqResponse.Valid)

        // Step 2: Setup Repositories
        setupUC := container.GetSetupRepositoriesUseCase()
        setupRequest := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: workingDir,
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "hello-world",
                    URL:    "https://github.com/octocat/Hello-World.git",
                    Branch: "master",
                },
            },
            Force:  false,
            DryRun: false,
        }

        setupResponse, err := setupUC.Execute(ctx, setupRequest)
        require.NoError(t, err)
        assert.Equal(t, 1, setupResponse.TotalRepositories)
        assert.Equal(t, 1, setupResponse.SuccessfulSetups)
        assert.Equal(t, 0, setupResponse.FailedSetups)

        // Step 3: Verify repository was cloned
        repoPath := filepath.Join(workingDir, "hello-world")
        assert.DirExists(t, repoPath)
        assert.FileExists(t, filepath.Join(repoPath, "README"))

        // Step 4: Verify configuration was created
        assert.FileExists(t, configPath)

        // Step 5: Test idempotency - run setup again
        setupResponse2, err := setupUC.Execute(ctx, setupRequest)
        require.NoError(t, err)
        assert.Equal(t, 1, setupResponse2.TotalRepositories)
        assert.Equal(t, 0, setupResponse2.SuccessfulSetups) // Should be skipped
        assert.Equal(t, 1, setupResponse2.SkippedRepositories)
    })

    t.Run("setup with multiple repositories", func(t *testing.T) {
        setupUC := container.GetSetupRepositoriesUseCase()
        setupRequest := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: workingDir,
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "spoon-knife",
                    URL:    "https://github.com/octocat/Spoon-Knife.git",
                    Branch: "main",
                },
                {
                    Name:   "hello-world-2",
                    URL:    "https://github.com/octocat/Hello-World.git",
                    Branch: "master",
                },
            },
            Force:  false,
            DryRun: false,
        }

        setupResponse, err := setupUC.Execute(ctx, setupRequest)
        require.NoError(t, err)
        assert.Equal(t, 2, setupResponse.TotalRepositories)
        assert.Equal(t, 2, setupResponse.SuccessfulSetups)

        // Verify both repositories exist
        assert.DirExists(t, filepath.Join(workingDir, "spoon-knife"))
        assert.DirExists(t, filepath.Join(workingDir, "hello-world-2"))
    })

    t.Run("setup with force flag", func(t *testing.T) {
        setupUC := container.GetSetupRepositoriesUseCase()
        setupRequest := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: workingDir,
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "hello-world", // Already exists
                    URL:    "https://github.com/octocat/Hello-World.git",
                    Branch: "master",
                },
            },
            Force:  true, // Force re-setup
            DryRun: false,
        }

        setupResponse, err := setupUC.Execute(ctx, setupRequest)
        require.NoError(t, err)
        assert.Equal(t, 1, setupResponse.TotalRepositories)
        assert.Equal(t, 1, setupResponse.SuccessfulSetups)
        assert.Equal(t, 0, setupResponse.SkippedRepositories)
    })
}

func TestErrorHandlingScenario(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, "config.json")
    workingDir := filepath.Join(tmpDir, "repositories")

    appConfig := &config.AppConfig{
        ConfigPath:          configPath,
        WorkingDir:          workingDir,
        DatabasePath:        "",
        GitTimeout:          30 * time.Second,
        DockerTimeout:       60 * time.Second,
        LogLevel:            "debug",
        EnableNotifications: true,
    }

    container, err := config.NewContainer(appConfig)
    require.NoError(t, err)
    defer container.Close()

    ctx := context.Background()

    t.Run("handle invalid repository URL", func(t *testing.T) {
        setupUC := container.GetSetupRepositoriesUseCase()
        setupRequest := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: workingDir,
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "invalid-repo",
                    URL:    "https://github.com/nonexistent/nonexistent.git",
                    Branch: "main",
                },
            },
            Force:  false,
            DryRun: false,
        }

        setupResponse, err := setupUC.Execute(ctx, setupRequest)
        require.NoError(t, err) // Use case handles errors gracefully
        assert.Equal(t, 1, setupResponse.TotalRepositories)
        assert.Equal(t, 0, setupResponse.SuccessfulSetups)
        assert.Equal(t, 1, setupResponse.FailedSetups)
        assert.Len(t, setupResponse.Errors, 1)
    })

    t.Run("handle permission denied", func(t *testing.T) {
        // Create a directory with no write permissions
        restrictedDir := filepath.Join(tmpDir, "restricted")
        err := os.MkdirAll(restrictedDir, 0444) // Read-only
        require.NoError(t, err)

        setupUC := container.GetSetupRepositoriesUseCase()
        setupRequest := &dtos.SetupRepositoriesRequest{
            WorkingDirectory: restrictedDir,
            Repositories: []*dtos.RepositorySetupInfo{
                {
                    Name:   "permission-test",
                    URL:    "https://github.com/octocat/Hello-World.git",
                    Branch: "master",
                },
            },
            Force:  false,
            DryRun: false,
        }

        setupResponse, err := setupUC.Execute(ctx, setupRequest)
        require.NoError(t, err)
        assert.Equal(t, 1, setupResponse.FailedSetups)
        assert.Len(t, setupResponse.Errors, 1)
    })
}
```

### Step 4: Performance Testing

#### Create `test/performance/load_test.go`

```go
package performance

import (
    "context"
    "fmt"
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"

    "github.com/fabianoflorentino/whiterose/internal/application/dtos"
    "github.com/fabianoflorentino/whiterose/internal/infrastructure/config"
)

func TestConcurrentSetupPerformance(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping performance test in short mode")
    }

    const (
        concurrentUsers = 10
        reposPerUser    = 5
        maxDuration     = 2 * time.Minute
    )

    // Setup
    appConfig := &config.AppConfig{
        ConfigPath:          "",
        WorkingDir:          t.TempDir(),
        DatabasePath:        "",
        GitTimeout:          30 * time.Second,
        DockerTimeout:       60 * time.Second,
        LogLevel:            "warn", // Reduce logging for performance
        EnableNotifications: false,
    }

    container, err := config.NewContainer(appConfig)
    if err != nil {
        t.Fatalf("Failed to create container: %v", err)
    }
    defer container.Close()

    ctx, cancel := context.WithTimeout(context.Background(), maxDuration)
    defer cancel()

    // Performance test
    start := time.Now()
    var wg sync.WaitGroup
    results := make(chan result, concurrentUsers)

    for i := 0; i < concurrentUsers; i++ {
        wg.Add(1)
        go func(userID int) {
            defer wg.Done()

            userResult := result{
                UserID: userID,
                Start:  time.Now(),
            }

            // Create repositories for this user
            repositories := make([]*dtos.RepositorySetupInfo, reposPerUser)
            for j := 0; j < reposPerUser; j++ {
                repositories[j] = &dtos.RepositorySetupInfo{
                    Name:   fmt.Sprintf("user-%d-repo-%d", userID, j),
                    URL:    "https://github.com/octocat/Hello-World.git",
                    Branch: "master",
                }
            }

            setupUC := container.GetSetupRepositoriesUseCase()
            request := &dtos.SetupRepositoriesRequest{
                WorkingDirectory: fmt.Sprintf("%s/user-%d", appConfig.WorkingDir, userID),
                Repositories:     repositories,
                Force:            false,
                DryRun:           false,
            }

            response, err := setupUC.Execute(ctx, request)
            userResult.End = time.Now()
            userResult.Duration = userResult.End.Sub(userResult.Start)
            userResult.Error = err
            if response != nil {
                userResult.SuccessfulSetups = response.SuccessfulSetups
                userResult.FailedSetups = response.FailedSetups
            }

            results <- userResult
        }(i)
    }

    // Wait for all goroutines to complete
    wg.Wait()
    close(results)

    totalDuration := time.Since(start)

    // Collect and analyze results
    var (
        totalSuccessful = 0
        totalFailed     = 0
        totalDuration   = time.Duration(0)
        maxDuration     = time.Duration(0)
        errors          = 0
    )

    for result := range results {
        totalSuccessful += result.SuccessfulSetups
        totalFailed += result.FailedSetups
        totalDuration += result.Duration
        if result.Duration > maxDuration {
            maxDuration = result.Duration
        }
        if result.Error != nil {
            errors++
        }
    }

    avgDuration := totalDuration / time.Duration(concurrentUsers)

    // Performance assertions
    t.Logf("Performance Results:")
    t.Logf("  Concurrent Users: %d", concurrentUsers)
    t.Logf("  Repos per User: %d", reposPerUser)
    t.Logf("  Total Duration: %v", totalDuration)
    t.Logf("  Average Duration: %v", avgDuration)
    t.Logf("  Max Duration: %v", maxDuration)
    t.Logf("  Total Successful: %d", totalSuccessful)
    t.Logf("  Total Failed: %d", totalFailed)
    t.Logf("  Errors: %d", errors)

    // Assertions
    assert.Equal(t, 0, errors, "No errors should occur during load test")
    assert.Equal(t, concurrentUsers*reposPerUser, totalSuccessful, "All repositories should be set up successfully")
    assert.LessOrEqual(t, maxDuration, 30*time.Second, "Maximum setup time should be reasonable")
    assert.LessOrEqual(t, avgDuration, 15*time.Second, "Average setup time should be reasonable")
}

type result struct {
    UserID           int
    Start            time.Time
    End              time.Time
    Duration         time.Duration
    SuccessfulSetups int
    FailedSetups     int
    Error            error
}

func BenchmarkSetupSingleRepository(b *testing.B) {
    appConfig := &config.AppConfig{
        ConfigPath:          "",
        WorkingDir:          b.TempDir(),
        DatabasePath:        "",
        GitTimeout:          30 * time.Second,
        DockerTimeout:       60 * time.Second,
        LogLevel:            "error",
        EnableNotifications: false,
    }

    container, err := config.NewContainer(appConfig)
    if err != nil {
        b.Fatalf("Failed to create container: %v", err)
    }
    defer container.Close()

    setupUC := container.GetSetupRepositoriesUseCase()
    request := &dtos.SetupRepositoriesRequest{
        WorkingDirectory: appConfig.WorkingDir,
        Repositories: []*dtos.RepositorySetupInfo{
            {
                Name:   "benchmark-repo",
                URL:    "https://github.com/octocat/Hello-World.git",
                Branch: "master",
            },
        },
        Force:  true, // Always force for consistent benchmarking
        DryRun: false,
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := setupUC.Execute(context.Background(), request)
        if err != nil {
            b.Errorf("Setup failed: %v", err)
        }
    }
}

// Memory usage benchmark
func BenchmarkMemoryUsage(b *testing.B) {
    appConfig := &config.AppConfig{
        ConfigPath:          "",
        WorkingDir:          b.TempDir(),
        DatabasePath:        "",
        GitTimeout:          30 * time.Second,
        DockerTimeout:       60 * time.Second,
        LogLevel:            "error",
        EnableNotifications: false,
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        container, err := config.NewContainer(appConfig)
        if err != nil {
            b.Fatalf("Failed to create container: %v", err)
        }
        container.Close()
    }
}
```

## ✅ Acceptance Criteria

- [ ] **Unit Tests**: >90% code coverage across all layers
- [ ] **Integration Tests**: All external dependencies tested
- [ ] **E2E Tests**: Complete user workflows validated
- [ ] **Performance Tests**: Load and stress testing completed
- [ ] **Documentation**: Comprehensive API and user documentation
- [ ] **CI/CD Pipeline**: Automated testing and deployment
- [ ] **Security**: Security scanning and vulnerability assessment
- [ ] **Monitoring**: Health checks and observability implemented

## 🧪 Testing Strategy Summary

### Test Pyramid Implementation

1. **Unit Tests (70%)**
   - Domain entities and value objects
   - Application use cases and services
   - Infrastructure adapters (with mocks)
   - Fast execution (<1s per test suite)

2. **Integration Tests (20%)**
   - Database persistence
   - External API integrations
   - File system operations
   - Moderate execution time (<30s per test suite)

3. **E2E Tests (10%)**
   - Complete user workflows
   - CLI command integration
   - Configuration migration
   - Slower execution (<5min per test suite)

### Coverage Goals

- **Overall Coverage**: >90%
- **Domain Layer**: >95% (critical business logic)
- **Application Layer**: >90% (use cases and services)
- **Infrastructure Layer**: >85% (adapters and external integrations)

### Test Environment

- **Local Development**: All tests runnable locally
- **CI/CD Pipeline**: Automated test execution on all PRs
- **Performance Testing**: Dedicated environment for load testing
- **Security Testing**: Automated security scanning

## 📚 Next Steps

After completing Phase 5:

1. ✅ Comprehensive testing strategy implemented
2. ✅ Documentation complete and up-to-date
3. ✅ CI/CD pipeline functional
4. ✅ Production readiness achieved
5. 🎉 **Migration Complete!** - Ready for production deployment

## 🔗 Related Documentation

- [Migration & Integration Guide](phase-4-migration.md)
- [Testing Best Practices](../guides/testing.md)
- [Performance Optimization](../guides/performance.md)
- [Production Deployment](../guides/deployment.md)
