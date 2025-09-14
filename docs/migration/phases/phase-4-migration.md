# Phase 4: Migration & Integration

## 🎯 Objective

Integrate all layers (domain, application, infrastructure) and implement a complete migration strategy from the current procedural approach to the new hexagonal architecture.

## ⏱️ Duration: 3-4 days | Effort: 24-32h

## 📋 Prerequisites

- ✅ Phase 3 completed (Infrastructure layer implemented)
- ✅ All adapters and ports working
- ✅ Dependency injection container configured
- ✅ Unit and integration tests passing

## 🎯 Goals

1. **Integration Layer**: Complete CLI integration with new architecture
2. **Migration Strategy**: Gradual migration of existing commands
3. **Configuration Migration**: Convert existing configuration to new format
4. **Command Migration**: Migrate existing commands to use new architecture
5. **Backward Compatibility**: Ensure smooth transition without breaking changes
6. **Performance Optimization**: Optimize for production use

## 📁 Directory Structure

Update the main project structure:

```text
cmd/
├── cli/                   # New CLI structure
│   ├── main.go           # Entry point
│   ├── root.go           # Root command (migrated)
│   ├── setup.go          # Setup command (migrated)
│   ├── docker.go         # Docker command (migrated)
│   └── prereq.go         # Prerequisites command (migrated)
├── api/                  # Future API server
│   └── main.go           # API entry point
└── migration/            # Migration utilities
    ├── config_migrator.go
    └── command_migrator.go

internal/
├── cli/                  # CLI-specific code
│   ├── commands/         # Command implementations
│   │   ├── setup.go
│   │   ├── docker.go
│   │   ├── prereq.go
│   │   └── version.go
│   ├── handlers/         # CLI handlers
│   └── middleware/       # CLI middleware
└── migration/           # Migration utilities
    ├── legacy/          # Legacy code wrappers
    └── compatibility/   # Backward compatibility
```

## 🚀 Implementation Steps

### Step 1: Create New CLI Entry Point

#### Update `cmd/cli/main.go`

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/fabianoflorentino/whiterose/internal/infrastructure/config"
    "github.com/fabianoflorentino/whiterose/internal/cli/commands"
)

var (
    cfgFile     string
    workingDir  string
    verbose     bool
    container   *config.Container
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "whiterose",
    Short: "Development environment manager for multiple repositories",
    Long: `WhiteRose is a CLI tool that helps manage multiple Git repositories,
Docker containers, and development environments in a structured way.

This tool follows hexagonal architecture principles and provides:
- Multi-repository management
- Docker environment orchestration
- Prerequisites validation
- Configuration management`,
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        initializeContainer()
    },
    PersistentPostRun: func(cmd *cobra.Command, args []string) {
        if container != nil {
            container.Close()
        }
    },
}

func main() {
    // Setup graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle signals
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("\nShutting down gracefully...")
        cancel()
        if container != nil {
            container.Close()
        }
        os.Exit(0)
    }()

    // Execute root command
    if err := rootCmd.ExecuteContext(ctx); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    // Global flags
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.whiterose.json)")
    rootCmd.PersistentFlags().StringVar(&workingDir, "working-dir", "./repositories", "working directory for repositories")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

    // Add commands
    rootCmd.AddCommand(commands.NewSetupCommand())
    rootCmd.AddCommand(commands.NewDockerCommand())
    rootCmd.AddCommand(commands.NewPrereqCommand())
    rootCmd.AddCommand(commands.NewVersionCommand())
}

// initConfig reads in config file and ENV variables
func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)

        viper.AddConfigPath(home)
        viper.SetConfigType("json")
        viper.SetConfigName(".whiterose")
    }

    // Environment variables
    viper.SetEnvPrefix("WHITEROSE")
    viper.AutomaticEnv()

    // Read config file
    if err := viper.ReadInConfig(); err == nil {
        if verbose {
            fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
        }
    }
}

// initializeContainer creates and configures the dependency injection container
func initializeContainer() {
    // Get config file path
    configPath := cfgFile
    if configPath == "" {
        home, err := os.UserHomeDir()
        if err != nil {
            log.Fatalf("Failed to get home directory: %v", err)
        }
        configPath = filepath.Join(home, ".whiterose.json")
    }

    // Prepare configuration
    appConfig := &config.AppConfig{
        ConfigPath:          configPath,
        WorkingDir:          workingDir,
        DatabasePath:        "", // Use in-memory for CLI
        GitTimeout:          30 * time.Second,
        DockerTimeout:       60 * time.Second,
        LogLevel:            getLogLevel(),
        EnableNotifications: true,
    }

    // Create container
    var err error
    container, err = config.NewContainer(appConfig)
    if err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
    }
}

func getLogLevel() string {
    if verbose {
        return "debug"
    }
    return "info"
}
```

### Step 2: Migrate Setup Command

#### Create `internal/cli/commands/setup.go`

```go
package commands

import (
    "context"
    "fmt"

    "github.com/spf13/cobra"

    "github.com/fabianoflorentino/whiterose/internal/application/dtos"
    "github.com/fabianoflorentino/whiterose/internal/infrastructure/config"
)

// SetupOptions holds setup command options
type SetupOptions struct {
    ConfigFile      string
    WorkingDir      string
    RepositoryURLs  []string
    SkipValidation  bool
    Force           bool
    DryRun          bool
}

// NewSetupCommand creates a new setup command
func NewSetupCommand() *cobra.Command {
    opts := &SetupOptions{}

    cmd := &cobra.Command{
        Use:   "setup",
        Short: "Setup development environment with repositories",
        Long: `Setup command initializes the development environment by:
- Validating system prerequisites
- Cloning or updating configured repositories
- Setting up Docker environment
- Creating necessary configuration files`,
        RunE: func(cmd *cobra.Command, args []string) error {
            return runSetup(cmd.Context(), opts)
        },
    }

    // Add flags
    cmd.Flags().StringSliceVar(&opts.RepositoryURLs, "repo", []string{}, "Repository URLs to setup (can be used multiple times)")
    cmd.Flags().BoolVar(&opts.SkipValidation, "skip-validation", false, "Skip prerequisites validation")
    cmd.Flags().BoolVar(&opts.Force, "force", false, "Force setup even if repositories exist")
    cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false, "Show what would be done without executing")

    return cmd
}

// runSetup executes the setup command
func runSetup(ctx context.Context, opts *SetupOptions) error {
    // Get container from context (injected by root command)
    container, err := getContainerFromContext(ctx)
    if err != nil {
        return fmt.Errorf("failed to get application container: %w", err)
    }

    // 1. Validate prerequisites (unless skipped)
    if !opts.SkipValidation {
        fmt.Println("🔍 Validating prerequisites...")
        
        prereqUC := container.GetValidatePrerequisitesUseCase()
        prereqRequest := &dtos.ValidatePrerequisitesRequest{
            SkipOptional: false,
        }

        prereqResponse, err := prereqUC.Execute(ctx, prereqRequest)
        if err != nil {
            return fmt.Errorf("prerequisites validation failed: %w", err)
        }

        if !prereqResponse.Valid {
            fmt.Println("❌ Prerequisites validation failed:")
            for _, issue := range prereqResponse.Issues {
                fmt.Printf("  - %s: %s\n", issue.Component, issue.Message)
            }
            return fmt.Errorf("please fix prerequisites before continuing")
        }

        fmt.Println("✅ Prerequisites validation passed")
    }

    // 2. Setup repositories
    fmt.Println("📦 Setting up repositories...")
    
    setupUC := container.GetSetupRepositoriesUseCase()
    
    // Prepare repository configurations
    repositories := make([]*dtos.RepositorySetupInfo, 0)
    
    // Add repositories from command line arguments
    for _, url := range opts.RepositoryURLs {
        repositories = append(repositories, &dtos.RepositorySetupInfo{
            URL:    url,
            Branch: "main", // Default branch
        })
    }

    setupRequest := &dtos.SetupRepositoriesRequest{
        WorkingDirectory: opts.WorkingDir,
        Repositories:     repositories,
        Force:            opts.Force,
        DryRun:           opts.DryRun,
    }

    setupResponse, err := setupUC.Execute(ctx, setupRequest)
    if err != nil {
        return fmt.Errorf("repository setup failed: %w", err)
    }

    // Display results
    fmt.Printf("✅ Setup completed successfully!\n")
    fmt.Printf("📊 Summary:\n")
    fmt.Printf("  - Total repositories: %d\n", setupResponse.TotalRepositories)
    fmt.Printf("  - Successfully setup: %d\n", setupResponse.SuccessfulSetups)
    fmt.Printf("  - Failed setups: %d\n", setupResponse.FailedSetups)
    fmt.Printf("  - Skipped (already exists): %d\n", setupResponse.SkippedRepositories)

    if len(setupResponse.Errors) > 0 {
        fmt.Println("\n⚠️  Errors encountered:")
        for repo, err := range setupResponse.Errors {
            fmt.Printf("  - %s: %v\n", repo, err)
        }
    }

    return nil
}

// Helper function to get container from context
func getContainerFromContext(ctx context.Context) (*config.Container, error) {
    // This would be injected by the root command
    // For now, we'll create a new container (this should be improved)
    
    appConfig := &config.AppConfig{
        ConfigPath:          "", // Will use default
        WorkingDir:          "./repositories",
        DatabasePath:        "",
        GitTimeout:          30 * time.Second,
        DockerTimeout:       60 * time.Second,
        LogLevel:            "info",
        EnableNotifications: true,
    }

    return config.NewContainer(appConfig)
}
```

### Step 3: Migration Strategy Implementation

#### Create `internal/migration/config_migrator.go`

```go
package migration

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"

    configAdapter "github.com/fabianoflorentino/whiterose/internal/infrastructure/adapters/config"
)

// ConfigMigrator handles migration of configuration files
type ConfigMigrator struct {
    legacyConfigPath string
    newConfigPath    string
}

// NewConfigMigrator creates a new configuration migrator
func NewConfigMigrator(legacyPath, newPath string) *ConfigMigrator {
    return &ConfigMigrator{
        legacyConfigPath: legacyPath,
        newConfigPath:    newPath,
    }
}

// LegacyConfig represents the old configuration format
type LegacyConfig struct {
    Repositories []LegacyRepository `yaml:"repositories"`
    Docker       LegacyDocker       `yaml:"docker"`
    Settings     LegacySettings     `yaml:"settings"`
}

// LegacyRepository represents old repository configuration
type LegacyRepository struct {
    Name   string `yaml:"name"`
    URL    string `yaml:"url"`
    Branch string `yaml:"branch"`
    Path   string `yaml:"path"`
}

// LegacyDocker represents old Docker configuration
type LegacyDocker struct {
    ComposeFile string   `yaml:"compose_file"`
    Images      []string `yaml:"images"`
}

// LegacySettings represents old settings
type LegacySettings struct {
    WorkingDir string `yaml:"working_dir"`
    LogLevel   string `yaml:"log_level"`
}

// MigrateConfig migrates configuration from legacy format to new format
func (m *ConfigMigrator) MigrateConfig() error {
    // Check if legacy config exists
    if _, err := os.Stat(m.legacyConfigPath); os.IsNotExist(err) {
        return fmt.Errorf("legacy config file not found: %s", m.legacyConfigPath)
    }

    // Check if new config already exists
    if _, err := os.Stat(m.newConfigPath); err == nil {
        return fmt.Errorf("new config file already exists: %s", m.newConfigPath)
    }

    // Read legacy config
    legacyData, err := os.ReadFile(m.legacyConfigPath)
    if err != nil {
        return fmt.Errorf("failed to read legacy config: %w", err)
    }

    // Parse legacy config
    var legacyConfig LegacyConfig
    if err := yaml.Unmarshal(legacyData, &legacyConfig); err != nil {
        return fmt.Errorf("failed to parse legacy config: %w", err)
    }

    // Convert to new format
    newConfig := m.convertToNewFormat(&legacyConfig)

    // Write new config
    if err := m.writeNewConfig(newConfig); err != nil {
        return fmt.Errorf("failed to write new config: %w", err)
    }

    // Backup legacy config
    backupPath := m.legacyConfigPath + ".backup"
    if err := os.Rename(m.legacyConfigPath, backupPath); err != nil {
        return fmt.Errorf("failed to backup legacy config: %w", err)
    }

    fmt.Printf("✅ Configuration migrated successfully!\n")
    fmt.Printf("  - New config: %s\n", m.newConfigPath)
    fmt.Printf("  - Legacy backup: %s\n", backupPath)

    return nil
}

// convertToNewFormat converts legacy configuration to new format
func (m *ConfigMigrator) convertToNewFormat(legacy *LegacyConfig) *configAdapter.Configuration {
    newConfig := &configAdapter.Configuration{
        WorkingDirectory: legacy.Settings.WorkingDir,
        Repositories:     make([]*configAdapter.RepositoryConfig, len(legacy.Repositories)),
        Docker: &configAdapter.DockerConfig{
            Registry:    "docker.io",
            Images:      []*configAdapter.DockerImage{},
            Networks:    []string{"default"},
            Volumes:     make(map[string]string),
            ComposeFile: legacy.Docker.ComposeFile,
        },
        Git: &configAdapter.GitConfig{
            DefaultBranch: "main",
            Timeout:       "30s",
        },
        Environment: make(map[string]string),
        Validation: &configAdapter.ValidationConfig{
            RequiredCommands: []configAdapter.CommandRequirement{
                {Command: "git", MinVersion: "2.0.0", Required: true},
                {Command: "docker", MinVersion: "20.0.0", Required: true},
            },
            OptionalCommands: []configAdapter.CommandRequirement{
                {Command: "docker-compose", MinVersion: "1.29.0", Required: false},
            },
        },
    }

    // Convert repositories
    for i, repo := range legacy.Repositories {
        newConfig.Repositories[i] = &configAdapter.RepositoryConfig{
            Name:        repo.Name,
            URL:         repo.URL,
            Branch:      repo.Branch,
            Path:        repo.Path,
            Tags:        []string{},
            Environment: make(map[string]string),
        }
    }

    // Convert Docker images
    for _, imageName := range legacy.Docker.Images {
        newConfig.Docker.Images = append(newConfig.Docker.Images, &configAdapter.DockerImage{
            Name:       imageName,
            Tag:        "latest",
            Dockerfile: "Dockerfile",
            Context:    ".",
            Args:       make(map[string]string),
        })
    }

    return newConfig
}

// writeNewConfig writes the new configuration to file
func (m *ConfigMigrator) writeNewConfig(config *configAdapter.Configuration) error {
    // Ensure directory exists
    dir := filepath.Dir(m.newConfigPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }

    // Marshal to JSON
    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    // Write to file
    if err := os.WriteFile(m.newConfigPath, data, 0644); err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }

    return nil
}

// ValidateMigration validates the migrated configuration
func (m *ConfigMigrator) ValidateMigration() error {
    // Read and validate new config
    configAdapter := configAdapter.NewJSONConfigAdapter(m.newConfigPath)
    _, err := configAdapter.LoadConfig(context.Background())
    if err != nil {
        return fmt.Errorf("migrated configuration is invalid: %w", err)
    }

    fmt.Println("✅ Migrated configuration is valid")
    return nil
}
```

### Step 4: Backward Compatibility Layer

#### Create `internal/migration/compatibility/legacy_wrapper.go`

```go
package compatibility

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "github.com/fabianoflorentino/whiterose/internal/infrastructure/config"
    "github.com/fabianoflorentino/whiterose/internal/migration"
)

// LegacyWrapper provides backward compatibility for legacy commands
type LegacyWrapper struct {
    container     *config.Container
    migrator      *migration.ConfigMigrator
    legacyDetected bool
}

// NewLegacyWrapper creates a new backward compatibility wrapper
func NewLegacyWrapper() *LegacyWrapper {
    return &LegacyWrapper{}
}

// Initialize checks for legacy configuration and sets up compatibility
func (lw *LegacyWrapper) Initialize() error {
    // Check for legacy configuration files
    legacyPaths := []string{
        ".whiterose.yaml",
        ".whiterose.yml",
        "whiterose.yaml",
        "whiterose.yml",
    }

    var legacyPath string
    for _, path := range legacyPaths {
        if _, err := os.Stat(path); err == nil {
            legacyPath = path
            lw.legacyDetected = true
            break
        }
    }

    if !lw.legacyDetected {
        return nil // No legacy configuration found
    }

    // Setup migrator
    home, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }

    newConfigPath := filepath.Join(home, ".whiterose.json")
    lw.migrator = migration.NewConfigMigrator(legacyPath, newConfigPath)

    fmt.Println("🔄 Legacy configuration detected!")
    fmt.Printf("  Legacy file: %s\n", legacyPath)
    fmt.Printf("  New file: %s\n", newConfigPath)

    return nil
}

// MigrateIfNeeded performs migration if legacy configuration is detected
func (lw *LegacyWrapper) MigrateIfNeeded() error {
    if !lw.legacyDetected {
        return nil
    }

    fmt.Println("🚀 Starting configuration migration...")

    // Perform migration
    if err := lw.migrator.MigrateConfig(); err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }

    // Validate migration
    if err := lw.migrator.ValidateMigration(); err != nil {
        return fmt.Errorf("migration validation failed: %w", err)
    }

    fmt.Println("✅ Migration completed successfully!")
    fmt.Println("💡 Legacy commands will continue to work during the transition period.")

    return nil
}

// WrapLegacyCommand wraps a legacy command to use the new architecture
func (lw *LegacyWrapper) WrapLegacyCommand(commandName string, args []string) error {
    if lw.container == nil {
        return fmt.Errorf("container not initialized")
    }

    switch commandName {
    case "setup":
        return lw.wrapSetupCommand(args)
    case "docker":
        return lw.wrapDockerCommand(args)
    case "prereq":
        return lw.wrapPrereqCommand(args)
    default:
        return fmt.Errorf("unknown legacy command: %s", commandName)
    }
}

// wrapSetupCommand wraps the legacy setup command
func (lw *LegacyWrapper) wrapSetupCommand(args []string) error {
    fmt.Println("🔄 Running legacy setup command with new architecture...")

    // Parse legacy arguments and convert to new format
    // This is a simplified example - you'd need to implement proper argument parsing

    setupUC := lw.container.GetSetupRepositoriesUseCase()
    
    // Execute setup with default parameters
    // In a real implementation, you'd parse the legacy arguments
    request := &dtos.SetupRepositoriesRequest{
        WorkingDirectory: "./repositories",
        Repositories:     []*dtos.RepositorySetupInfo{}, // Would be populated from config
        Force:            false,
        DryRun:           false,
    }

    _, err := setupUC.Execute(context.Background(), request)
    if err != nil {
        return fmt.Errorf("setup command failed: %w", err)
    }

    fmt.Println("✅ Legacy setup command completed")
    return nil
}

// wrapDockerCommand wraps the legacy docker command
func (lw *LegacyWrapper) wrapDockerCommand(args []string) error {
    fmt.Println("🔄 Running legacy docker command with new architecture...")

    dockerUC := lw.container.GetManageDockerImagesUseCase()
    
    // Execute docker management with default parameters
    request := &dtos.ManageDockerImagesRequest{
        Action:    "list", // Default action
        ImageName: "",
        Force:     false,
    }

    _, err := dockerUC.Execute(context.Background(), request)
    if err != nil {
        return fmt.Errorf("docker command failed: %w", err)
    }

    fmt.Println("✅ Legacy docker command completed")
    return nil
}

// wrapPrereqCommand wraps the legacy prereq command
func (lw *LegacyWrapper) wrapPrereqCommand(args []string) error {
    fmt.Println("🔄 Running legacy prereq command with new architecture...")

    prereqUC := lw.container.GetValidatePrerequisitesUseCase()
    
    request := &dtos.ValidatePrerequisitesRequest{
        SkipOptional: false,
    }

    response, err := prereqUC.Execute(context.Background(), request)
    if err != nil {
        return fmt.Errorf("prereq command failed: %w", err)
    }

    if !response.Valid {
        fmt.Println("❌ Prerequisites validation failed")
        return fmt.Errorf("prerequisites not met")
    }

    fmt.Println("✅ Legacy prereq command completed")
    return nil
}

// SetContainer sets the application container
func (lw *LegacyWrapper) SetContainer(container *config.Container) {
    lw.container = container
}
```

## ✅ Acceptance Criteria

- [ ] Complete CLI integration with new architecture
- [ ] All existing commands migrated and functional
- [ ] Configuration migration working smoothly
- [ ] Backward compatibility maintained
- [ ] Performance optimized for production use
- [ ] Integration tests covering migration scenarios
- [ ] Documentation updated for new commands
- [ ] Legacy commands still work during transition

## 🧪 Testing Strategy

### Migration Tests

```go
func TestConfigMigration(t *testing.T) {
    // Setup legacy config
    legacyConfig := `
repositories:
  - name: "test-repo"
    url: "https://github.com/test/repo.git"
    branch: "main"
    path: "./repos/test-repo"
docker:
  compose_file: "docker-compose.yml"
  images: ["nginx", "redis"]
settings:
  working_dir: "./repositories"
  log_level: "info"
`

    tmpDir := t.TempDir()
    legacyPath := filepath.Join(tmpDir, "legacy.yaml")
    newPath := filepath.Join(tmpDir, "new.json")

    err := os.WriteFile(legacyPath, []byte(legacyConfig), 0644)
    require.NoError(t, err)

    // Test migration
    migrator := migration.NewConfigMigrator(legacyPath, newPath)
    err = migrator.MigrateConfig()
    assert.NoError(t, err)

    // Verify new config exists and is valid
    assert.FileExists(t, newPath)
    assert.FileExists(t, legacyPath+".backup")

    // Validate new config
    err = migrator.ValidateMigration()
    assert.NoError(t, err)
}
```

## 📚 Next Steps

After completing Phase 4:

1. ✅ All layers integrated and working together
2. ✅ Migration strategy implemented and tested
3. ✅ Backward compatibility ensured
4. 🚀 Proceed to [Phase 5: Testing & Documentation](phase-5-testing.md)

## 🔗 Related Documentation

- [Infrastructure Layer Guide](phase-3-infrastructure.md)
- [Code Examples - Integration](../code-examples/USAGE.md)
- [Migration Strategy](../migration/migration-strategy.md)
- [Backward Compatibility Guide](../how-to/backward-compatibility.md)
