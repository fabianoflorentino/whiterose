# 🎯 Phase 1: Domain Layer (Foundation)

## 📋 Objective

Create the solid foundation of the application with domain entities, business rules, and well-defined interfaces. This phase establishes the core of hexagonal architecture.

## ⏱️ Duration: 2-3 days | Effort: 16-20h

## 🎯 Deliverables

### 1. Folder Structure

```bash
internal/
└── domain/
    ├── entities/
    │   ├── repository.go
    │   ├── application.go
    │   ├── docker_image.go
    │   └── environment.go
    ├── repositories/
    │   ├── git_repository.go
    │   ├── config_repository.go
    │   ├── docker_repository.go
    │   └── environment_repository.go
    ├── services/
    │   ├── git_service.go
    │   ├── validation_service.go
    │   └── docker_service.go
    └── errors/
        └── domain_errors.go
```

### 2. Complete Implementations

## 📄 Complete File Code

### `internal/domain/entities/repository.go`

```go
package entities

import (
    "fmt"
    "net/url"
    "strings"
    "time"
)

// Repository represents a Git repository in the domain
type Repository struct {
    ID          string
    URL         string
    Name        string
    Directory   string
    Branch      string
    AuthMethod  AuthenticationMethod
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// AuthType defines supported authentication types
type AuthType string

const (
    AuthTypeSSH   AuthType = "ssh"
    AuthTypeHTTPS AuthType = "https"
)

// AuthenticationMethod contains authentication information
type AuthenticationMethod struct {
    Type     AuthType
    Username string
    Token    string
    SSHKey   SSHKeyConfig
}

// SSHKeyConfig contains SSH key configurations
type SSHKeyConfig struct {
    Path     string
    Name     string
    Passphrase string
}

// NewRepository creates a new Repository instance
func NewRepository(url, directory string) (*Repository, error) {
    if err := validateURL(url); err != nil {
        return nil, fmt.Errorf("invalid repository URL: %w", err)
    }
    
    if directory == "" {
        return nil, fmt.Errorf("directory cannot be empty")
    }
    
    name := extractRepositoryName(url)
    authMethod := determineAuthMethod(url)
    
    return &Repository{
        ID:         generateID(url),
        URL:        url,
        Name:       name,
        Directory:  directory,
        Branch:     "main", // default branch
        AuthMethod: authMethod,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}

// SetBranch sets the repository branch
func (r *Repository) SetBranch(branch string) error {
    if branch == "" {
        return fmt.Errorf("branch cannot be empty")
    }
    
    r.Branch = branch
    r.UpdatedAt = time.Now()
    return nil
}

// SetAuthentication configures the authentication method
func (r *Repository) SetAuthentication(auth AuthenticationMethod) error {
    if err := auth.Validate(); err != nil {
        return fmt.Errorf("invalid authentication: %w", err)
    }
    
    r.AuthMethod = auth
    r.UpdatedAt = time.Now()
    return nil
}

// IsSSH returns true if the repository uses SSH authentication
func (r *Repository) IsSSH() bool {
    return r.AuthMethod.Type == AuthTypeSSH
}

// IsHTTPS returns true if the repository uses HTTPS authentication
func (r *Repository) IsHTTPS() bool {
    return r.AuthMethod.Type == AuthTypeHTTPS
}

// Validate validates the Repository entity
func (r *Repository) Validate() error {
    if r.URL == "" {
        return fmt.Errorf("URL is required")
    }
    
    if r.Directory == "" {
        return fmt.Errorf("directory is required")
    }
    
    if r.Branch == "" {
        return fmt.Errorf("branch is required")
    }
    
    return r.AuthMethod.Validate()
}

// Validate validates the authentication method
func (a *AuthenticationMethod) Validate() error {
    switch a.Type {
    case AuthTypeSSH:
        if a.SSHKey.Path == "" && a.SSHKey.Name == "" {
            return fmt.Errorf("SSH key path or name is required")
        }
    case AuthTypeHTTPS:
        if a.Username == "" || a.Token == "" {
            return fmt.Errorf("username and token are required for HTTPS")
        }
    default:
        return fmt.Errorf("invalid auth type: %s", a.Type)
    }
    
    return nil
}

// Helper functions
func validateURL(rawURL string) error {
    if rawURL == "" {
        return fmt.Errorf("URL cannot be empty")
    }
    
    // Validate Git URLs (both HTTPS and SSH)
    if strings.HasPrefix(rawURL, "git@") {
        // SSH format: git@github.com:user/repo.git
        return nil
    }
    
    // Validate HTTPS URLs
    if strings.HasPrefix(rawURL, "https://") {
        _, err := url.Parse(rawURL)
        return err
    }
    
    return fmt.Errorf("unsupported URL format")
}

func extractRepositoryName(rawURL string) string {
    // Extract repository name from URL
    if strings.HasPrefix(rawURL, "git@") {
        // SSH format: git@github.com:user/repo.git
        parts := strings.Split(rawURL, ":")
        if len(parts) >= 2 {
            repoPath := parts[len(parts)-1]
            repoPath = strings.TrimSuffix(repoPath, ".git")
            pathParts := strings.Split(repoPath, "/")
            if len(pathParts) > 0 {
                return pathParts[len(pathParts)-1]
            }
        }
    } else if strings.HasPrefix(rawURL, "https://") {
        // HTTPS format: https://github.com/user/repo.git
        u, err := url.Parse(rawURL)
        if err == nil {
            path := strings.TrimPrefix(u.Path, "/")
            path = strings.TrimSuffix(path, ".git")
            pathParts := strings.Split(path, "/")
            if len(pathParts) > 0 {
                return pathParts[len(pathParts)-1]
            }
        }
    }
    
    return "unknown"
}

func determineAuthMethod(rawURL string) AuthenticationMethod {
    if strings.HasPrefix(rawURL, "git@") {
        return AuthenticationMethod{
            Type: AuthTypeSSH,
            SSHKey: SSHKeyConfig{
                Name: "id_rsa", // default
            },
        }
    }
    
    return AuthenticationMethod{
        Type: AuthTypeHTTPS,
    }
}

func generateID(url string) string {
    // Simple ID generation based on URL hash
    // In production, consider using a proper UUID
    return fmt.Sprintf("repo_%d", len(url)*int(time.Now().Unix())%10000)
}
```

### `internal/domain/entities/application.go`

```go
package entities

import (
    "fmt"
    "strings"
    "time"
)

// Application representa uma aplicação/ferramenta no sistema
type Application struct {
    ID                  string
    Name                string
    Command             string
    VersionFlag         string
    RecommendedVersion  string
    InstallInstructions map[string]string
    CreatedAt          time.Time
    UpdatedAt          time.Time
}

// OperatingSystem defines supported operating systems
type OperatingSystem string

const (
    OSLinux   OperatingSystem = "linux"
    OSDarwin  OperatingSystem = "darwin"
    OSWindows OperatingSystem = "windows"
)

// ApplicationStatus representa o status de uma aplicação
type ApplicationStatus struct {
    Application     *Application
    IsInstalled     bool
    InstalledVersion string
    IsUpToDate      bool
    ErrorMessage    string
    CheckedAt       time.Time
}

// NewApplication cria uma nova instância de Application
func NewApplication(name, command, versionFlag, recommendedVersion string) (*Application, error) {
    if err := validateApplicationData(name, command, versionFlag); err != nil {
        return nil, err
    }
    
    return &Application{
        ID:                 generateApplicationID(name),
        Name:               name,
        Command:            command,
        VersionFlag:        versionFlag,
        RecommendedVersion: recommendedVersion,
        InstallInstructions: make(map[string]string),
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }, nil
}

// AddInstallInstruction adiciona instrução de instalação para um OS
func (a *Application) AddInstallInstruction(os OperatingSystem, instruction string) error {
    if instruction == "" {
        return fmt.Errorf("instruction cannot be empty")
    }
    
    a.InstallInstructions[string(os)] = instruction
    a.UpdatedAt = time.Now()
    return nil
}

// GetInstallInstruction retorna a instrução de instalação para um OS
func (a *Application) GetInstallInstruction(os OperatingSystem) (string, bool) {
    instruction, exists := a.InstallInstructions[string(os)]
    return instruction, exists
}

// Validate valida a entidade Application
func (a *Application) Validate() error {
    return validateApplicationData(a.Name, a.Command, a.VersionFlag)
}

// NewApplicationStatus cria um novo status de aplicação
func NewApplicationStatus(app *Application) *ApplicationStatus {
    return &ApplicationStatus{
        Application: app,
        IsInstalled: false,
        IsUpToDate:  false,
        CheckedAt:   time.Now(),
    }
}

// SetInstalled marca a aplicação como instalada com versão
func (as *ApplicationStatus) SetInstalled(version string) {
    as.IsInstalled = true
    as.InstalledVersion = version
    as.IsUpToDate = as.checkVersionCompatibility()
    as.ErrorMessage = ""
    as.CheckedAt = time.Now()
}

// SetNotInstalled marca a aplicação como não instalada
func (as *ApplicationStatus) SetNotInstalled(errorMsg string) {
    as.IsInstalled = false
    as.InstalledVersion = ""
    as.IsUpToDate = false
    as.ErrorMessage = errorMsg
    as.CheckedAt = time.Now()
}

// checkVersionCompatibility verifica se a versão instalada é compatível
func (as *ApplicationStatus) checkVersionCompatibility() bool {
    if as.Application.RecommendedVersion == "" || as.InstalledVersion == "" {
        return true // Se não há versão recomendada, considera compatível
    }
    
    // Basic version comparison implementation
    // In production, use a library like github.com/Masterminds/semver
    return strings.Contains(as.InstalledVersion, as.Application.RecommendedVersion)
}

// Helper functions
func validateApplicationData(name, command, versionFlag string) error {
    if name == "" {
        return fmt.Errorf("application name cannot be empty")
    }
    
    if command == "" {
        return fmt.Errorf("application command cannot be empty")
    }
    
    if versionFlag == "" {
        return fmt.Errorf("version flag cannot be empty")
    }
    
    return nil
}

func generateApplicationID(name string) string {
    // Simple ID generation
    normalizedName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
    return fmt.Sprintf("app_%s_%d", normalizedName, time.Now().Unix()%10000)
}
```

### `internal/domain/entities/docker_image.go`

```go
package entities

import (
    "fmt"
    "strings"
    "time"
)

// DockerImage representa uma imagem Docker no domínio
type DockerImage struct {
    ID          string
    Name        string
    Tag         string
    FullName    string
    Size        int64
    Created     time.Time
    BuildArgs   map[string]string
    Target      string
    Context     string
    Dockerfile  string
}

// DockerBuildOptions contém opções para build de imagem
type DockerBuildOptions struct {
    ImageName    string
    Tag          string
    Dockerfile   string
    Context      string
    BuildArgs    map[string]string
    Target       string
    NoCache      bool
    Progress     string
}

// DockerImageStatus representa o status de uma imagem
type DockerImageStatus struct {
    Image       *DockerImage
    Exists      bool
    IsBuilding  bool
    BuildError  string
    LastChecked time.Time
}

// NewDockerImage cria uma nova instância de DockerImage
func NewDockerImage(name, tag string) (*DockerImage, error) {
    if err := validateDockerImageData(name, tag); err != nil {
        return nil, err
    }
    
    fullName := buildFullImageName(name, tag)
    
    return &DockerImage{
        ID:       generateDockerImageID(fullName),
        Name:     name,
        Tag:      tag,
        FullName: fullName,
        BuildArgs: make(map[string]string),
        Created:  time.Now(),
    }, nil
}

// AddBuildArg adiciona um argumento de build
func (di *DockerImage) AddBuildArg(key, value string) error {
    if key == "" {
        return fmt.Errorf("build arg key cannot be empty")
    }
    
    di.BuildArgs[key] = value
    return nil
}

// SetDockerfile define o caminho do Dockerfile
func (di *DockerImage) SetDockerfile(path string) error {
    if path == "" {
        return fmt.Errorf("dockerfile path cannot be empty")
    }
    
    di.Dockerfile = path
    return nil
}

// SetContext define o contexto de build
func (di *DockerImage) SetContext(context string) error {
    if context == "" {
        return fmt.Errorf("build context cannot be empty")
    }
    
    di.Context = context
    return nil
}

// SetTarget define o target de build multi-stage
func (di *DockerImage) SetTarget(target string) {
    di.Target = target
}

// Validate valida a entidade DockerImage
func (di *DockerImage) Validate() error {
    if err := validateDockerImageData(di.Name, di.Tag); err != nil {
        return err
    }
    
    if di.Dockerfile != "" && !strings.HasSuffix(di.Dockerfile, "Dockerfile") {
        return fmt.Errorf("invalid dockerfile path: %s", di.Dockerfile)
    }
    
    return nil
}

// NewDockerBuildOptions cria novas opções de build
func NewDockerBuildOptions(imageName, tag string) (*DockerBuildOptions, error) {
    if err := validateDockerImageData(imageName, tag); err != nil {
        return nil, err
    }
    
    return &DockerBuildOptions{
        ImageName: imageName,
        Tag:       tag,
        Context:   ".",
        BuildArgs: make(map[string]string),
        NoCache:   false,
        Progress:  "auto",
    }, nil
}

// AddBuildArg adiciona um argumento de build às opções
func (dbo *DockerBuildOptions) AddBuildArg(key, value string) error {
    if key == "" {
        return fmt.Errorf("build arg key cannot be empty")
    }
    
    dbo.BuildArgs[key] = value
    return nil
}

// SetDockerfile define o Dockerfile nas opções
func (dbo *DockerBuildOptions) SetDockerfile(path string) error {
    if path == "" {
        return fmt.Errorf("dockerfile path cannot be empty")
    }
    
    dbo.Dockerfile = path
    return nil
}

// GetFullImageName retorna o nome completo da imagem
func (dbo *DockerBuildOptions) GetFullImageName() string {
    return buildFullImageName(dbo.ImageName, dbo.Tag)
}

// Validate valida as opções de build
func (dbo *DockerBuildOptions) Validate() error {
    if err := validateDockerImageData(dbo.ImageName, dbo.Tag); err != nil {
        return err
    }
    
    if dbo.Context == "" {
        return fmt.Errorf("build context cannot be empty")
    }
    
    return nil
}

// NewDockerImageStatus cria um novo status de imagem
func NewDockerImageStatus(image *DockerImage) *DockerImageStatus {
    return &DockerImageStatus{
        Image:       image,
        Exists:      false,
        IsBuilding:  false,
        LastChecked: time.Now(),
    }
}

// SetExists marca a imagem como existente
func (dis *DockerImageStatus) SetExists(exists bool) {
    dis.Exists = exists
    dis.LastChecked = time.Now()
}

// SetBuilding marca a imagem como em processo de build
func (dis *DockerImageStatus) SetBuilding(building bool) {
    dis.IsBuilding = building
    if building {
        dis.BuildError = ""
    }
    dis.LastChecked = time.Now()
}

// SetBuildError define um erro de build
func (dis *DockerImageStatus) SetBuildError(err string) {
    dis.BuildError = err
    dis.IsBuilding = false
    dis.LastChecked = time.Now()
}

// Helper functions
func validateDockerImageData(name, tag string) error {
    if name == "" {
        return fmt.Errorf("image name cannot be empty")
    }
    
    if tag == "" {
        return fmt.Errorf("image tag cannot be empty")
    }
    
    // Validação básica de nome de imagem Docker
    if strings.Contains(name, " ") {
        return fmt.Errorf("image name cannot contain spaces")
    }
    
    return nil
}

func buildFullImageName(name, tag string) string {
    return fmt.Sprintf("%s:%s", name, tag)
}

func generateDockerImageID(fullName string) string {
    return fmt.Sprintf("img_%d", len(fullName)*int(time.Now().Unix())%10000)
}
```

### `internal/domain/entities/environment.go`

```go
package entities

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

// Environment representa o ambiente de desenvolvimento
type Environment struct {
    ID                string
    HomeDir           string
    WorkingDir        string
    ConfigFile        string
    EnvFile           string
    SSHKeyPath        string
    SSHKeyName        string
    GitUser           string
    GitToken          string
    DockerImageName   string
    DockerImageVersion string
    DockerfilePath    string
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

// EnvironmentConfig contém configurações do ambiente
type EnvironmentConfig struct {
    ConfigFile         string
    SSHKeyPath        string
    SSHKeyName        string
    GitUser           string
    GitToken          string
    ImageName         string
    ImageVersion      string
    DockerfilePath    string
    BuildTarget       string
}

// NewEnvironment cria uma nova instância de Environment
func NewEnvironment() (*Environment, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("failed to get home directory: %w", err)
    }
    
    workingDir, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("failed to get working directory: %w", err)
    }
    
    return &Environment{
        ID:                generateEnvironmentID(),
        HomeDir:           homeDir,
        WorkingDir:        workingDir,
        ConfigFile:        filepath.Join(homeDir, ".config.json"),
        EnvFile:           filepath.Join(homeDir, ".env"),
        SSHKeyPath:        filepath.Join(homeDir, ".ssh"),
        SSHKeyName:        "id_rsa",
        DockerImageName:   "my_app",
        DockerImageVersion: "latest",
        DockerfilePath:    filepath.Join(workingDir, "Dockerfile"),
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }, nil
}

// LoadFromEnv carrega configurações das variáveis de ambiente
func (e *Environment) LoadFromEnv() {
    if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
        e.ConfigFile = configFile
    }
    
    if sshKeyPath := os.Getenv("SSH_KEY_PATH"); sshKeyPath != "" {
        e.SSHKeyPath = sshKeyPath
    }
    
    if sshKeyName := os.Getenv("SSH_KEY_NAME"); sshKeyName != "" {
        e.SSHKeyName = sshKeyName
    }
    
    if gitUser := os.Getenv("GIT_USER"); gitUser != "" {
        e.GitUser = gitUser
    }
    
    if gitToken := os.Getenv("GIT_TOKEN"); gitToken != "" {
        e.GitToken = gitToken
    }
    
    if imageName := os.Getenv("IMAGE_NAME"); imageName != "" {
        e.DockerImageName = imageName
    }
    
    if imageVersion := os.Getenv("IMAGE_VERSION"); imageVersion != "" {
        e.DockerImageVersion = imageVersion
    }
    
    if dockerfilePath := os.Getenv("DOCKERFILE_PATH"); dockerfilePath != "" {
        e.DockerfilePath = dockerfilePath
    }
    
    e.UpdatedAt = time.Now()
}

// GetFullSSHKeyPath retorna o caminho completo da chave SSH
func (e *Environment) GetFullSSHKeyPath() string {
    return filepath.Join(e.SSHKeyPath, e.SSHKeyName)
}

// GetFullImageName retorna o nome completo da imagem Docker
func (e *Environment) GetFullImageName() string {
    return fmt.Sprintf("%s:%s", e.DockerImageName, e.DockerImageVersion)
}

// Validate validates the environment configuration
func (e *Environment) Validate() error {
    if e.HomeDir == "" {
        return fmt.Errorf("home directory is required")
    }
    
    if e.WorkingDir == "" {
        return fmt.Errorf("working directory is required")
    }
    
    // Check if configuration file exists
    if _, err := os.Stat(e.ConfigFile); os.IsNotExist(err) {
        return fmt.Errorf("config file not found: %s", e.ConfigFile)
    }
    
    // Check if .env file exists
    if _, err := os.Stat(e.EnvFile); os.IsNotExist(err) {
        return fmt.Errorf("env file not found: %s", e.EnvFile)
    }
    
    return nil
}

// ToConfig converte para EnvironmentConfig
func (e *Environment) ToConfig() *EnvironmentConfig {
    return &EnvironmentConfig{
        ConfigFile:      e.ConfigFile,
        SSHKeyPath:     e.SSHKeyPath,
        SSHKeyName:     e.SSHKeyName,
        GitUser:        e.GitUser,
        GitToken:       e.GitToken,
        ImageName:      e.DockerImageName,
        ImageVersion:   e.DockerImageVersion,
        DockerfilePath: e.DockerfilePath,
    }
}

// NewEnvironmentConfig creates a new environment configuration
func NewEnvironmentConfig() *EnvironmentConfig {
    return &EnvironmentConfig{
        SSHKeyName:   "id_rsa",
        ImageName:    "my_app",
        ImageVersion: "latest",
        BuildTarget:  "development",
    }
}

// SetDefaults define valores padrão para campos vazios
func (ec *EnvironmentConfig) SetDefaults() {
    homeDir, _ := os.UserHomeDir()
    workingDir, _ := os.Getwd()
    
    if ec.ConfigFile == "" {
        ec.ConfigFile = filepath.Join(homeDir, ".config.json")
    }
    
    if ec.SSHKeyPath == "" {
        ec.SSHKeyPath = filepath.Join(homeDir, ".ssh")
    }
    
    if ec.SSHKeyName == "" {
        ec.SSHKeyName = "id_rsa"
    }
    
    if ec.ImageName == "" {
        ec.ImageName = "my_app"
    }
    
    if ec.ImageVersion == "" {
        ec.ImageVersion = "latest"
    }
    
    if ec.DockerfilePath == "" {
        ec.DockerfilePath = filepath.Join(workingDir, "Dockerfile")
    }
    
    if ec.BuildTarget == "" {
        ec.BuildTarget = "development"
    }
}

// Validate validates the configuration
func (ec *EnvironmentConfig) Validate() error {
    if ec.ConfigFile == "" {
        return fmt.Errorf("config file path is required")
    }
    
    if ec.ImageName == "" {
        return fmt.Errorf("image name is required")
    }
    
    if ec.ImageVersion == "" {
        return fmt.Errorf("image version is required")
    }
    
    return nil
}

// Helper function
func generateEnvironmentID() string {
    return fmt.Sprintf("env_%d", time.Now().Unix()%10000)
}
```

### `internal/domain/errors/domain_errors.go`

```go
package errors

import (
    "fmt"
)

// DomainErrorCode representa códigos de erro do domínio
type DomainErrorCode string

const (
    // Repository errors
    ErrCodeRepositoryNotFound     DomainErrorCode = "REPOSITORY_NOT_FOUND"
    ErrCodeRepositoryInvalid      DomainErrorCode = "REPOSITORY_INVALID"
    ErrCodeRepositoryExists       DomainErrorCode = "REPOSITORY_EXISTS"
    ErrCodeBranchNotFound         DomainErrorCode = "BRANCH_NOT_FOUND"
    ErrCodeAuthenticationFailed   DomainErrorCode = "AUTHENTICATION_FAILED"
    
    // Application errors
    ErrCodeApplicationNotFound    DomainErrorCode = "APPLICATION_NOT_FOUND"
    ErrCodeApplicationInvalid     DomainErrorCode = "APPLICATION_INVALID"
    ErrCodeVersionIncompatible    DomainErrorCode = "VERSION_INCOMPATIBLE"
    
    // Docker errors
    ErrCodeDockerImageNotFound    DomainErrorCode = "DOCKER_IMAGE_NOT_FOUND"
    ErrCodeDockerImageInvalid     DomainErrorCode = "DOCKER_IMAGE_INVALID"
    ErrCodeDockerBuildFailed      DomainErrorCode = "DOCKER_BUILD_FAILED"
    ErrCodeDockerfileNotFound     DomainErrorCode = "DOCKERFILE_NOT_FOUND"
    
    // Environment errors
    ErrCodeEnvironmentInvalid     DomainErrorCode = "ENVIRONMENT_INVALID"
    ErrCodeConfigNotFound         DomainErrorCode = "CONFIG_NOT_FOUND"
    ErrCodeConfigInvalid          DomainErrorCode = "CONFIG_INVALID"
    
    // General errors
    ErrCodeValidationFailed       DomainErrorCode = "VALIDATION_FAILED"
    ErrCodeOperationFailed        DomainErrorCode = "OPERATION_FAILED"
    ErrCodePermissionDenied       DomainErrorCode = "PERMISSION_DENIED"
)

// DomainError representa um erro do domínio
type DomainError struct {
    Code     DomainErrorCode
    Message  string
    Cause    error
    Context  map[string]interface{}
}

// Error implementa a interface error
func (de *DomainError) Error() string {
    if de.Cause != nil {
        return fmt.Sprintf("%s: %s (caused by: %v)", de.Code, de.Message, de.Cause)
    }
    return fmt.Sprintf("%s: %s", de.Code, de.Message)
}

// Unwrap permite o uso com errors.Is e errors.As
func (de *DomainError) Unwrap() error {
    return de.Cause
}

// WithContext adiciona contexto ao erro
func (de *DomainError) WithContext(key string, value interface{}) *DomainError {
    if de.Context == nil {
        de.Context = make(map[string]interface{})
    }
    de.Context[key] = value
    return de
}

// NewDomainError cria um novo erro de domínio
func NewDomainError(code DomainErrorCode, message string) *DomainError {
    return &DomainError{
        Code:    code,
        Message: message,
        Context: make(map[string]interface{}),
    }
}

// NewDomainErrorWithCause cria um novo erro de domínio com causa
func NewDomainErrorWithCause(code DomainErrorCode, message string, cause error) *DomainError {
    return &DomainError{
        Code:    code,
        Message: message,
        Cause:   cause,
        Context: make(map[string]interface{}),
    }
}

// Repository Errors
func NewRepositoryNotFoundError(url string) *DomainError {
    return NewDomainError(ErrCodeRepositoryNotFound, "repository not found").
        WithContext("url", url)
}

func NewRepositoryInvalidError(reason string) *DomainError {
    return NewDomainError(ErrCodeRepositoryInvalid, "repository is invalid").
        WithContext("reason", reason)
}

func NewRepositoryExistsError(directory string) *DomainError {
    return NewDomainError(ErrCodeRepositoryExists, "repository already exists").
        WithContext("directory", directory)
}

func NewBranchNotFoundError(branch string) *DomainError {
    return NewDomainError(ErrCodeBranchNotFound, "branch not found").
        WithContext("branch", branch)
}

func NewAuthenticationFailedError(authType string, cause error) *DomainError {
    return NewDomainErrorWithCause(ErrCodeAuthenticationFailed, "authentication failed", cause).
        WithContext("authType", authType)
}

// Application Errors
func NewApplicationNotFoundError(name string) *DomainError {
    return NewDomainError(ErrCodeApplicationNotFound, "application not found").
        WithContext("name", name)
}

func NewApplicationInvalidError(reason string) *DomainError {
    return NewDomainError(ErrCodeApplicationInvalid, "application is invalid").
        WithContext("reason", reason)
}

func NewVersionIncompatibleError(installed, recommended string) *DomainError {
    return NewDomainError(ErrCodeVersionIncompatible, "version is incompatible").
        WithContext("installed", installed).
        WithContext("recommended", recommended)
}

// Docker Errors
func NewDockerImageNotFoundError(imageName string) *DomainError {
    return NewDomainError(ErrCodeDockerImageNotFound, "docker image not found").
        WithContext("imageName", imageName)
}

func NewDockerImageInvalidError(reason string) *DomainError {
    return NewDomainError(ErrCodeDockerImageInvalid, "docker image is invalid").
        WithContext("reason", reason)
}

func NewDockerBuildFailedError(cause error) *DomainError {
    return NewDomainErrorWithCause(ErrCodeDockerBuildFailed, "docker build failed", cause)
}

func NewDockerfileNotFoundError(path string) *DomainError {
    return NewDomainError(ErrCodeDockerfileNotFound, "dockerfile not found").
        WithContext("path", path)
}

// Environment Errors
func NewEnvironmentInvalidError(reason string) *DomainError {
    return NewDomainError(ErrCodeEnvironmentInvalid, "environment is invalid").
        WithContext("reason", reason)
}

func NewConfigNotFoundError(path string) *DomainError {
    return NewDomainError(ErrCodeConfigNotFound, "configuration file not found").
        WithContext("path", path)
}

func NewConfigInvalidError(reason string, cause error) *DomainError {
    return NewDomainErrorWithCause(ErrCodeConfigInvalid, "configuration is invalid", cause).
        WithContext("reason", reason)
}

// General Errors
func NewValidationFailedError(field, reason string) *DomainError {
    return NewDomainError(ErrCodeValidationFailed, "validation failed").
        WithContext("field", field).
        WithContext("reason", reason)
}

func NewOperationFailedError(operation string, cause error) *DomainError {
    return NewDomainErrorWithCause(ErrCodeOperationFailed, "operation failed", cause).
        WithContext("operation", operation)
}

func NewPermissionDeniedError(resource string) *DomainError {
    return NewDomainError(ErrCodePermissionDenied, "permission denied").
        WithContext("resource", resource)
}

// IsDomainError verifica se um erro é um DomainError
func IsDomainError(err error) bool {
    _, ok := err.(*DomainError)
    return ok
}

// GetDomainError extrai um DomainError de um erro
func GetDomainError(err error) (*DomainError, bool) {
    domainErr, ok := err.(*DomainError)
    return domainErr, ok
}

// HasErrorCode verifica se um erro tem um código específico
func HasErrorCode(err error, code DomainErrorCode) bool {
    if domainErr, ok := GetDomainError(err); ok {
        return domainErr.Code == code
    }
    return false
}
```

## ✅ Acceptance Criteria

- [ ] All entities modeled and validated
- [ ] Repository interfaces defined
- [ ] Business rules centralized
- [ ] Error system working
- [ ] Entity unit tests (>85% coverage)

## 🧪 Required Tests

### Entity Tests

```go
// internal/domain/entities/repository_test.go
func TestNewRepository(t *testing.T) {
    // Test valid repository creation
    // Test invalid URL
    // Test empty directory
}

func TestRepository_SetBranch(t *testing.T) {
    // Test valid branch
    // Test empty branch
}

func TestRepository_SetAuthentication(t *testing.T) {
    // Test SSH auth
    // Test HTTPS auth
    // Test invalid auth
}
```

## 🚀 Next Step

After completing this phase, continue to [Phase 2: Application Layer](phase-2-application.md).
