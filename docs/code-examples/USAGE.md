# Usage Examples

This directory contains practical examples of how to use the hexagonal architecture implementation.

## Quick Start Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/fabianoflorentino/whiterose/docs/code-examples/infrastructure/di"
    "github.com/fabianoflorentino/whiterose/docs/code-examples/application/usecases"
)

func main() {
    // Initialize dependency injection container
    container, err := di.NewContainerFromEnv()
    if err != nil {
        log.Fatal("Failed to initialize container:", err)
    }

    // Example 1: Setup repositories
    if err := setupRepositoriesExample(container); err != nil {
        log.Fatal("Repository setup failed:", err)
    }

    // Example 2: Validate system prerequisites
    if err := validateSystemExample(container); err != nil {
        log.Fatal("System validation failed:", err)
    }
}

func setupRepositoriesExample(container *di.Container) error {
    ctx := context.Background()
    
    // Get the use case from container
    setupUC := container.GetSetupRepositoriesUseCase()
    
    // Prepare request
    request := usecases.SetupRepositoriesRequest{
        Repositories: []usecases.RepositorySetupData{
            {
                Name:   "my-project",
                URL:    "https://github.com/user/my-project.git",
                Branch: "main",
            },
            {
                Name:   "another-project",
                URL:    "https://github.com/user/another-project.git",
                Branch: "develop",
            },
        },
        ForceClone: false,
    }
    
    // Execute use case
    response, err := setupUC.Execute(ctx, request)
    if err != nil {
        return fmt.Errorf("setup failed: %w", err)
    }
    
    // Display results
    fmt.Printf("Setup completed: %d success, %d failed\n", 
        response.SuccessCount, response.FailureCount)
    
    for _, result := range response.SetupResults {
        fmt.Printf("- %s: %s (%s)\n", result.Name, result.Status, result.Message)
    }
    
    return nil
}

func validateSystemExample(container *di.Container) error {
    ctx := context.Background()
    
    // Get validation repository
    validationRepo := container.GetValidationRepository()
    
    // Validate environment
    results, err := validationRepo.ValidateEnvironment(ctx)
    if err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Display results
    fmt.Println("System Validation Results:")
    for _, result := range results {
        status := "✅"
        if result.Status == "fail" {
            status = "❌"
        } else if result.Status == "warning" {
            status = "⚠️"
        }
        
        fmt.Printf("%s %s: %s\n", status, result.Check, result.Message)
    }
    
    return nil
}
```

## Testing Example

```go
package main

import (
    "context"
    "testing"

    "github.com/fabianoflorentino/whiterose/docs/code-examples/domain/entities"
    "github.com/fabianoflorentino/whiterose/docs/code-examples/infrastructure/adapters"
)

func TestRepositoryEntity(t *testing.T) {
    // Test entity creation
    repo, err := entities.NewRepository(
        "test-repo", 
        "https://github.com/user/test.git", 
        "main",
    )
    if err != nil {
        t.Fatalf("Failed to create repository: %v", err)
    }

    // Test entity behavior
    if repo.Name() != "test-repo" {
        t.Errorf("Expected name 'test-repo', got '%s'", repo.Name())
    }

    if repo.IsCloned() {
        t.Error("Repository should not be marked as cloned initially")
    }

    // Test setting local path
    err = repo.SetLocalPath("/tmp/test-repo")
    if err != nil {
        t.Fatalf("Failed to set local path: %v", err)
    }

    // Test marking as cloned
    repo.MarkAsCloned()
    if !repo.IsCloned() {
        t.Error("Repository should be marked as cloned")
    }
}

func TestInMemoryRepositoryAdapter(t *testing.T) {
    ctx := context.Background()
    adapter := adapters.NewInMemoryRepositoryAdapter()

    // Create test repository
    repo, err := entities.NewRepository(
        "test-repo", 
        "https://github.com/user/test.git", 
        "main",
    )
    if err != nil {
        t.Fatalf("Failed to create repository: %v", err)
    }

    // Test save
    err = adapter.Save(ctx, repo)
    if err != nil {
        t.Fatalf("Failed to save repository: %v", err)
    }

    // Test find by name
    found, err := adapter.FindByName(ctx, "test-repo")
    if err != nil {
        t.Fatalf("Failed to find repository: %v", err)
    }

    if found.Name() != repo.Name() {
        t.Errorf("Expected name '%s', got '%s'", repo.Name(), found.Name())
    }

    // Test exists
    exists, err := adapter.Exists(ctx, "test-repo")
    if err != nil {
        t.Fatalf("Failed to check existence: %v", err)
    }

    if !exists {
        t.Error("Repository should exist")
    }
}
```

## CLI Integration Example

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/spf13/cobra"
    
    "github.com/fabianoflorentino/whiterose/docs/code-examples/infrastructure/di"
    "github.com/fabianoflorentino/whiterose/docs/code-examples/application/usecases"
)

func main() {
    // Initialize container
    container, err := di.NewContainerFromEnv()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
        os.Exit(1)
    }

    // Create root command
    rootCmd := &cobra.Command{
        Use:   "whiterose",
        Short: "WhiteRose development environment manager",
    }

    // Add subcommands
    rootCmd.AddCommand(createSetupCommand(container))
    rootCmd.AddCommand(createValidateCommand(container))
    rootCmd.AddCommand(createStatusCommand(container))

    // Execute
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func createSetupCommand(container *di.Container) *cobra.Command {
    var forceClone bool

    cmd := &cobra.Command{
        Use:   "setup",
        Short: "Setup development repositories",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            setupUC := container.GetSetupRepositoriesUseCase()

            // Load repositories from configuration
            configRepo := container.GetConfigRepository()
            repoConfigs, err := configRepo.GetRepositories(ctx)
            if err != nil {
                return fmt.Errorf("failed to load repositories: %w", err)
            }

            // Convert to use case format
            repositories := make([]usecases.RepositorySetupData, len(repoConfigs))
            for i, config := range repoConfigs {
                repositories[i] = usecases.RepositorySetupData{
                    Name:   config.Name,
                    URL:    config.URL,
                    Branch: config.Branch,
                }
            }

            // Execute setup
            request := usecases.SetupRepositoriesRequest{
                Repositories: repositories,
                ForceClone:   forceClone,
            }

            response, err := setupUC.Execute(ctx, request)
            if err != nil {
                return fmt.Errorf("setup failed: %w", err)
            }

            // Display results
            fmt.Printf("Setup completed: %d success, %d failed\n", 
                response.SuccessCount, response.FailureCount)

            return nil
        },
    }

    cmd.Flags().BoolVar(&forceClone, "force", false, "Force clone even if repository exists")
    
    return cmd
}

func createValidateCommand(container *di.Container) *cobra.Command {
    return &cobra.Command{
        Use:   "validate",
        Short: "Validate system prerequisites",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            validationRepo := container.GetValidationRepository()

            results, err := validationRepo.ValidateEnvironment(ctx)
            if err != nil {
                return fmt.Errorf("validation failed: %w", err)
            }

            fmt.Println("System Validation Results:")
            allPassed := true
            for _, result := range results {
                status := "✅"
                if result.Status == "fail" {
                    status = "❌"
                    allPassed = false
                } else if result.Status == "warning" {
                    status = "⚠️"
                }

                fmt.Printf("%s %s: %s\n", status, result.Check, result.Message)
            }

            if !allPassed {
                os.Exit(1)
            }

            return nil
        },
    }
}

func createStatusCommand(container *di.Container) *cobra.Command {
    return &cobra.Command{
        Use:   "status",
        Short: "Show repository status",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            setupUC := container.GetSetupRepositoriesUseCase()

            statuses, err := setupUC.GetRepositoryStatus(ctx)
            if err != nil {
                return fmt.Errorf("failed to get status: %w", err)
            }

            fmt.Println("Repository Status:")
            for _, status := range statuses {
                clonedStatus := "❌"
                if status.IsCloned {
                    clonedStatus = "✅"
                }

                fmt.Printf("%s %s (%s) - %s\n", 
                    clonedStatus, status.Name, status.Branch, status.LocalPath)
            }

            return nil
        },
    }
}
```

## Configuration Example

```json
{
  "working_directory": "./repositories",
  "repositories": [
    {
      "name": "frontend",
      "url": "https://github.com/company/frontend.git",
      "branch": "main",
      "path": "./repositories/frontend"
    },
    {
      "name": "backend",
      "url": "https://github.com/company/backend.git",
      "branch": "develop",
      "path": "./repositories/backend"
    }
  ],
  "docker": {
    "registry": "docker.io",
    "images": [
      {
        "name": "frontend",
        "tag": "latest",
        "dockerfile": "./frontend/Dockerfile",
        "context": "./frontend"
      }
    ],
    "networks": ["development"],
    "volumes": {
      "node_modules": "/app/node_modules"
    }
  },
  "environment": {
    "NODE_ENV": "development",
    "API_URL": "http://localhost:8080"
  }
}
```

## Best Practices

1. **Always use dependency injection** - Don't instantiate dependencies directly in your use cases
2. **Handle errors at the right layer** - Domain errors in domain layer, infrastructure errors in infrastructure layer
3. **Use interfaces for all external dependencies** - This enables easy testing and swapping implementations
4. **Keep your domain pure** - No external dependencies in domain entities and services
5. **Test each layer independently** - Use mocks for dependencies when testing use cases
6. **Use context for cancellation** - Always pass context through your call chain
