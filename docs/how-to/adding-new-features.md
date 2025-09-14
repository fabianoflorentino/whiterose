# 📋 How to Add New Features

## 🎯 Step-by-Step Guide for Implementing Features

This guide shows the exact implementation order for adding a new feature following hexagonal architecture. We'll use the implementation of a **"Configuration Backup"** feature as an example.

## 🏗️ Implementation Order

### 1️⃣ **Domain Layer** (First - Solid foundation)

#### 📍 1.1 Create Entity

**File**: `internal/domain/entities/backup.go`

```go
package entities

import (
    "fmt"
    "time"
)

// Backup represents a configuration backup
type Backup struct {
    ID          string
    Name        string
    FilePath    string
    BackupPath  string
    CreatedAt   time.Time
    Size        int64
    Checksum    string
}

// NewBackup creates a new Backup instance
func NewBackup(name, filePath, backupPath string) (*Backup, error) {
    if name == "" {
        return nil, fmt.Errorf("backup name cannot be empty")
    }
    
    if filePath == "" {
        return nil, fmt.Errorf("file path cannot be empty")
    }
    
    if backupPath == "" {
        return nil, fmt.Errorf("backup path cannot be empty")
    }
    
    return &Backup{
        ID:         generateBackupID(name),
        Name:       name,
        FilePath:   filePath,
        BackupPath: backupPath,
        CreatedAt:  time.Now(),
    }, nil
}

// Validate validates the Backup entity
func (b *Backup) Validate() error {
    if b.Name == "" {
        return fmt.Errorf("backup name is required")
    }
    
    if b.FilePath == "" {
        return fmt.Errorf("file path is required")
    }
    
    if b.BackupPath == "" {
        return fmt.Errorf("backup path is required")
    }
    
    return nil
}

func generateBackupID(name string) string {
    return fmt.Sprintf("backup_%s_%d", name, time.Now().Unix())
}
```

#### 📍 1.2 Define Repository Interface

**File**: `internal/domain/repositories/backup_repository.go`

```go
package repositories

import (
    "context"
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
)

// BackupRepository defines persistence operations for backups
type BackupRepository interface {
    Create(ctx context.Context, backup *entities.Backup) error
    GetByID(ctx context.Context, id string) (*entities.Backup, error)
    List(ctx context.Context) ([]*entities.Backup, error)
    Delete(ctx context.Context, id string) error
    Exists(ctx context.Context, filePath string) (bool, error)
}
```

#### 📍 1.3 Implement Domain Service

**File**: `internal/domain/services/backup_service.go`

```go
package services

import (
    "context"
    "fmt"
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
    "github.com/fabianoflorentino/whiterose/internal/domain/errors"
)

// BackupDomainService contains business rules for backups
type BackupDomainService struct{}

// NewBackupDomainService creates a new backup service
func NewBackupDomainService() *BackupDomainService {
    return &BackupDomainService{}
}

// ValidateBackup validates business rules for backup
func (s *BackupDomainService) ValidateBackup(ctx context.Context, backup *entities.Backup) error {
    if err := backup.Validate(); err != nil {
        return errors.NewValidationFailedError("backup", err.Error())
    }
    
    // Business rule: Backup cannot overwrite existing file
    if backup.FilePath == backup.BackupPath {
        return errors.NewValidationFailedError("backup", "backup path cannot be the same as source file")
    }
    
    return nil
}

// GenerateBackupName generates unique name for backup
func (s *BackupDomainService) GenerateBackupName(ctx context.Context, baseName string) string {
    timestamp := time.Now().Format("20060102_150405")
    return fmt.Sprintf("%s_%s", baseName, timestamp)
}
```

#### 📍 1.4 Add Domain Errors

**File**: `internal/domain/errors/domain_errors.go` (add)

```go
// Backup errors
const (
    ErrCodeBackupNotFound     DomainErrorCode = "BACKUP_NOT_FOUND"
    ErrCodeBackupExists       DomainErrorCode = "BACKUP_EXISTS"
    ErrCodeBackupFailed       DomainErrorCode = "BACKUP_FAILED"
    ErrCodeRestoreFailed      DomainErrorCode = "RESTORE_FAILED"
)

func NewBackupNotFoundError(id string) *DomainError {
    return NewDomainError(ErrCodeBackupNotFound, "backup not found").
        WithContext("id", id)
}

func NewBackupExistsError(path string) *DomainError {
    return NewDomainError(ErrCodeBackupExists, "backup already exists").
        WithContext("path", path)
}
```

### 2️⃣ **Application Layer** (Segundo - Orquestração)

#### 📍 2.1 Criar DTOs

**Arquivo**: `internal/application/dto/backup_dto.go`

```go
package dto

import "time"

// BackupRequest representa uma solicitação de backup
type BackupRequest struct {
    Name       string   `json:"name" validate:"required"`
    FilePaths  []string `json:"file_paths" validate:"required,min=1"`
    BackupDir  string   `json:"backup_dir" validate:"required"`
    Compress   bool     `json:"compress"`
}

// BackupResponse representa a resposta de uma operação de backup
type BackupResponse struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    FilePaths   []string  `json:"file_paths"`
    BackupPath  string    `json:"backup_path"`
    Size        int64     `json:"size"`
    CreatedAt   time.Time `json:"created_at"`
    Success     bool      `json:"success"`
    Message     string    `json:"message"`
}

// ListBackupsResponse representa a resposta de listagem de backups
type ListBackupsResponse struct {
    Backups []BackupInfo `json:"backups"`
    Total   int          `json:"total"`
}

// BackupInfo representa informações resumidas de um backup
type BackupInfo struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Size      int64     `json:"size"`
    CreatedAt time.Time `json:"created_at"`
}

// RestoreRequest representa uma solicitação de restore
type RestoreRequest struct {
    BackupID    string `json:"backup_id" validate:"required"`
    RestorePath string `json:"restore_path"`
    Overwrite   bool   `json:"overwrite"`
}

// RestoreResponse representa a resposta de uma operação de restore
type RestoreResponse struct {
    BackupID     string `json:"backup_id"`
    RestorePath  string `json:"restore_path"`
    Success      bool   `json:"success"`
    Message      string `json:"message"`
}
```

#### 📍 2.2 Definir Input Port

**Arquivo**: `internal/application/ports/input/backup_service.go`

```go
package input

import (
    "context"
    "github.com/fabianoflorentino/whiterose/internal/application/dto"
)

// BackupService define operações de backup disponíveis
type BackupService interface {
    CreateBackup(ctx context.Context, req *dto.BackupRequest) (*dto.BackupResponse, error)
    ListBackups(ctx context.Context) (*dto.ListBackupsResponse, error)
    RestoreBackup(ctx context.Context, req *dto.RestoreRequest) (*dto.RestoreResponse, error)
    DeleteBackup(ctx context.Context, backupID string) error
}
```

#### 📍 2.3 Definir Output Ports

**Arquivo**: `internal/application/ports/output/backup_port.go`

```go
package output

import (
    "context"
)

// BackupPort define operações de backup externas
type BackupPort interface {
    CreateBackup(ctx context.Context, sourcePath, backupPath string, compress bool) error
    RestoreBackup(ctx context.Context, backupPath, destPath string, overwrite bool) error
    CalculateChecksum(ctx context.Context, filePath string) (string, error)
    GetFileSize(ctx context.Context, filePath string) (int64, error)
    FileExists(ctx context.Context, filePath string) (bool, error)
}
```

#### 📍 2.4 Implementar Use Case

**Arquivo**: `internal/application/usecases/backup_usecase.go`

```go
package usecases

import (
    "context"
    "fmt"
    
    "github.com/fabianoflorentino/whiterose/internal/application/dto"
    "github.com/fabianoflorentino/whiterose/internal/application/ports/input"
    "github.com/fabianoflorentino/whiterose/internal/application/ports/output"
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
    "github.com/fabianoflorentino/whiterose/internal/domain/repositories"
    "github.com/fabianoflorentino/whiterose/internal/domain/services"
    "github.com/fabianoflorentino/whiterose/internal/domain/errors"
)

// BackupUseCase implementa as operações de backup
type BackupUseCase struct {
    backupRepo    repositories.BackupRepository
    backupPort    output.BackupPort
    backupService *services.BackupDomainService
    logger        output.LoggerPort
}

// NewBackupUseCase cria um novo use case de backup
func NewBackupUseCase(
    backupRepo repositories.BackupRepository,
    backupPort output.BackupPort,
    backupService *services.BackupDomainService,
    logger output.LoggerPort,
) input.BackupService {
    return &BackupUseCase{
        backupRepo:    backupRepo,
        backupPort:    backupPort,
        backupService: backupService,
        logger:        logger,
    }
}

// CreateBackup implementa a criação de backup
func (uc *BackupUseCase) CreateBackup(ctx context.Context, req *dto.BackupRequest) (*dto.BackupResponse, error) {
    uc.logger.Info(ctx, "Creating backup", map[string]interface{}{
        "name": req.Name,
        "files": len(req.FilePaths),
    })
    
    // Validar entrada
    if err := uc.validateBackupRequest(req); err != nil {
        return nil, err
    }
    
    // Gerar nome único se necessário
    backupName := uc.backupService.GenerateBackupName(ctx, req.Name)
    
    var responses []*dto.BackupResponse
    
    // Processar cada arquivo
    for _, filePath := range req.FilePaths {
        response, err := uc.processFileBackup(ctx, filePath, backupName, req)
        if err != nil {
            uc.logger.Error(ctx, "Failed to backup file", map[string]interface{}{
                "file": filePath,
                "error": err.Error(),
            })
            continue // Continua com próximo arquivo
        }
        responses = append(responses, response)
    }
    
    if len(responses) == 0 {
        return nil, errors.NewOperationFailedError("backup", fmt.Errorf("no files were backed up"))
    }
    
    // Retorna o primeiro backup criado (pode ser ajustado conforme necessidade)
    return responses[0], nil
}

// processFileBackup processa o backup de um arquivo individual
func (uc *BackupUseCase) processFileBackup(ctx context.Context, filePath, backupName string, req *dto.BackupRequest) (*dto.BackupResponse, error) {
    // Criar entidade de backup
    backupPath := fmt.Sprintf("%s/%s_%s", req.BackupDir, backupName, filepath.Base(filePath))
    backup, err := entities.NewBackup(backupName, filePath, backupPath)
    if err != nil {
        return nil, err
    }
    
    // Validar regras de negócio
    if err := uc.backupService.ValidateBackup(ctx, backup); err != nil {
        return nil, err
    }
    
    // Verificar se arquivo fonte existe
    exists, err := uc.backupPort.FileExists(ctx, filePath)
    if err != nil || !exists {
        return nil, errors.NewOperationFailedError("backup", fmt.Errorf("source file not found: %s", filePath))
    }
    
    // Executar backup
    if err := uc.backupPort.CreateBackup(ctx, filePath, backupPath, req.Compress); err != nil {
        return nil, errors.NewOperationFailedError("backup", err)
    }
    
    // Calcular checksum e tamanho
    checksum, err := uc.backupPort.CalculateChecksum(ctx, backupPath)
    if err != nil {
        uc.logger.Warn(ctx, "Failed to calculate checksum", map[string]interface{}{
            "file": backupPath,
            "error": err.Error(),
        })
    }
    backup.Checksum = checksum
    
    size, err := uc.backupPort.GetFileSize(ctx, backupPath)
    if err != nil {
        uc.logger.Warn(ctx, "Failed to get file size", map[string]interface{}{
            "file": backupPath,
            "error": err.Error(),
        })
    }
    backup.Size = size
    
    // Salvar no repositório
    if err := uc.backupRepo.Create(ctx, backup); err != nil {
        return nil, errors.NewOperationFailedError("backup", err)
    }
    
    return &dto.BackupResponse{
        ID:         backup.ID,
        Name:       backup.Name,
        FilePaths:  []string{filePath},
        BackupPath: backupPath,
        Size:       backup.Size,
        CreatedAt:  backup.CreatedAt,
        Success:    true,
        Message:    "Backup created successfully",
    }, nil
}

// ListBackups implementa a listagem de backups
func (uc *BackupUseCase) ListBackups(ctx context.Context) (*dto.ListBackupsResponse, error) {
    backups, err := uc.backupRepo.List(ctx)
    if err != nil {
        return nil, errors.NewOperationFailedError("list_backups", err)
    }
    
    var backupInfos []dto.BackupInfo
    for _, backup := range backups {
        backupInfos = append(backupInfos, dto.BackupInfo{
            ID:        backup.ID,
            Name:      backup.Name,
            Size:      backup.Size,
            CreatedAt: backup.CreatedAt,
        })
    }
    
    return &dto.ListBackupsResponse{
        Backups: backupInfos,
        Total:   len(backupInfos),
    }, nil
}

// RestoreBackup implementa o restore de backup
func (uc *BackupUseCase) RestoreBackup(ctx context.Context, req *dto.RestoreRequest) (*dto.RestoreResponse, error) {
    // Buscar backup
    backup, err := uc.backupRepo.GetByID(ctx, req.BackupID)
    if err != nil {
        return nil, errors.NewBackupNotFoundError(req.BackupID)
    }
    
    // Definir caminho de restore
    restorePath := req.RestorePath
    if restorePath == "" {
        restorePath = backup.FilePath // Restore no local original
    }
    
    // Executar restore
    if err := uc.backupPort.RestoreBackup(ctx, backup.BackupPath, restorePath, req.Overwrite); err != nil {
        return nil, errors.NewOperationFailedError("restore", err)
    }
    
    return &dto.RestoreResponse{
        BackupID:    req.BackupID,
        RestorePath: restorePath,
        Success:     true,
        Message:     "Backup restored successfully",
    }, nil
}

// DeleteBackup implementa a remoção de backup
func (uc *BackupUseCase) DeleteBackup(ctx context.Context, backupID string) error {
    // Buscar backup
    backup, err := uc.backupRepo.GetByID(ctx, backupID)
    if err != nil {
        return errors.NewBackupNotFoundError(backupID)
    }
    
    // Remover arquivo físico (opcional - pode ser configurável)
    exists, err := uc.backupPort.FileExists(ctx, backup.BackupPath)
    if err == nil && exists {
        // Implementar remoção física se necessário
        uc.logger.Info(ctx, "Physical backup file found", map[string]interface{}{
            "path": backup.BackupPath,
        })
    }
    
    // Remover do repositório
    return uc.backupRepo.Delete(ctx, backupID)
}

// validateBackupRequest valida a requisição de backup
func (uc *BackupUseCase) validateBackupRequest(req *dto.BackupRequest) error {
    if req.Name == "" {
        return errors.NewValidationFailedError("name", "backup name is required")
    }
    
    if len(req.FilePaths) == 0 {
        return errors.NewValidationFailedError("file_paths", "at least one file path is required")
    }
    
    if req.BackupDir == "" {
        return errors.NewValidationFailedError("backup_dir", "backup directory is required")
    }
    
    return nil
}
```

### 3️⃣ **Infrastructure Layer** (Terceiro - Adaptadores)

#### 📍 3.1 Implementar Adapter

**Arquivo**: `internal/infrastructure/adapters/backup/filesystem_backup.go`

```go
package backup

import (
    "context"
    "crypto/sha256"
    "fmt"
    "io"
    "os"
    "path/filepath"
    
    "github.com/fabianoflorentino/whiterose/internal/application/ports/output"
)

// FilesystemBackupAdapter implementa operações de backup no filesystem
type FilesystemBackupAdapter struct{}

// NewFilesystemBackupAdapter cria um novo adapter de backup
func NewFilesystemBackupAdapter() output.BackupPort {
    return &FilesystemBackupAdapter{}
}

// CreateBackup cria um backup de arquivo
func (a *FilesystemBackupAdapter) CreateBackup(ctx context.Context, sourcePath, backupPath string, compress bool) error {
    // Criar diretório de backup se não existir
    backupDir := filepath.Dir(backupPath)
    if err := os.MkdirAll(backupDir, 0755); err != nil {
        return fmt.Errorf("failed to create backup directory: %w", err)
    }
    
    // Abrir arquivo fonte
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return fmt.Errorf("failed to open source file: %w", err)
    }
    defer sourceFile.Close()
    
    // Criar arquivo de backup
    backupFile, err := os.Create(backupPath)
    if err != nil {
        return fmt.Errorf("failed to create backup file: %w", err)
    }
    defer backupFile.Close()
    
    // Copiar conteúdo
    _, err = io.Copy(backupFile, sourceFile)
    if err != nil {
        return fmt.Errorf("failed to copy file content: %w", err)
    }
    
    return nil
}

// RestoreBackup restaura um backup
func (a *FilesystemBackupAdapter) RestoreBackup(ctx context.Context, backupPath, destPath string, overwrite bool) error {
    // Verificar se arquivo de destino existe
    if !overwrite {
        if _, err := os.Stat(destPath); err == nil {
            return fmt.Errorf("destination file exists and overwrite is false")
        }
    }
    
    // Criar diretório de destino se não existir
    destDir := filepath.Dir(destPath)
    if err := os.MkdirAll(destDir, 0755); err != nil {
        return fmt.Errorf("failed to create destination directory: %w", err)
    }
    
    // Abrir arquivo de backup
    backupFile, err := os.Open(backupPath)
    if err != nil {
        return fmt.Errorf("failed to open backup file: %w", err)
    }
    defer backupFile.Close()
    
    // Criar arquivo de destino
    destFile, err := os.Create(destPath)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %w", err)
    }
    defer destFile.Close()
    
    // Copiar conteúdo
    _, err = io.Copy(destFile, backupFile)
    if err != nil {
        return fmt.Errorf("failed to copy backup content: %w", err)
    }
    
    return nil
}

// CalculateChecksum calcula o checksum de um arquivo
func (a *FilesystemBackupAdapter) CalculateChecksum(ctx context.Context, filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("failed to open file for checksum: %w", err)
    }
    defer file.Close()
    
    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", fmt.Errorf("failed to calculate checksum: %w", err)
    }
    
    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// GetFileSize retorna o tamanho de um arquivo
func (a *FilesystemBackupAdapter) GetFileSize(ctx context.Context, filePath string) (int64, error) {
    info, err := os.Stat(filePath)
    if err != nil {
        return 0, fmt.Errorf("failed to get file info: %w", err)
    }
    
    return info.Size(), nil
}

// FileExists verifica se um arquivo existe
func (a *FilesystemBackupAdapter) FileExists(ctx context.Context, filePath string) (bool, error) {
    _, err := os.Stat(filePath)
    if err == nil {
        return true, nil
    }
    
    if os.IsNotExist(err) {
        return false, nil
    }
    
    return false, err
}
```

#### 📍 3.2 Implementar Repository

**Arquivo**: `internal/infrastructure/adapters/backup/memory_backup_repository.go`

```go
package backup

import (
    "context"
    "fmt"
    "sync"
    
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
    "github.com/fabianoflorentino/whiterose/internal/domain/repositories"
)

// MemoryBackupRepository implementa backup repository em memória
type MemoryBackupRepository struct {
    backups map[string]*entities.Backup
    mutex   sync.RWMutex
}

// NewMemoryBackupRepository cria um novo repository em memória
func NewMemoryBackupRepository() repositories.BackupRepository {
    return &MemoryBackupRepository{
        backups: make(map[string]*entities.Backup),
    }
}

// Create salva um backup
func (r *MemoryBackupRepository) Create(ctx context.Context, backup *entities.Backup) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    if _, exists := r.backups[backup.ID]; exists {
        return fmt.Errorf("backup with ID %s already exists", backup.ID)
    }
    
    r.backups[backup.ID] = backup
    return nil
}

// GetByID busca um backup por ID
func (r *MemoryBackupRepository) GetByID(ctx context.Context, id string) (*entities.Backup, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    backup, exists := r.backups[id]
    if !exists {
        return nil, fmt.Errorf("backup with ID %s not found", id)
    }
    
    return backup, nil
}

// List retorna todos os backups
func (r *MemoryBackupRepository) List(ctx context.Context) ([]*entities.Backup, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    var backups []*entities.Backup
    for _, backup := range r.backups {
        backups = append(backups, backup)
    }
    
    return backups, nil
}

// Delete remove um backup
func (r *MemoryBackupRepository) Delete(ctx context.Context, id string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    if _, exists := r.backups[id]; !exists {
        return fmt.Errorf("backup with ID %s not found", id)
    }
    
    delete(r.backups, id)
    return nil
}

// Exists verifica se existe backup para um arquivo
func (r *MemoryBackupRepository) Exists(ctx context.Context, filePath string) (bool, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    for _, backup := range r.backups {
        if backup.FilePath == filePath {
            return true, nil
        }
    }
    
    return false, nil
}
```

### 4️⃣ **CLI Interface** (Quarto - Interface do usuário)

#### 📍 4.1 Criar CLI Handler

**Arquivo**: `cmd/handlers/backup_handler.go`

```go
package handlers

import (
    "context"
    "fmt"
    "strings"
    
    "github.com/spf13/cobra"
    
    "github.com/fabianoflorentino/whiterose/internal/application/dto"
    "github.com/fabianoflorentino/whiterose/internal/application/ports/input"
)

// BackupHandler gerencia comandos relacionados a backup
type BackupHandler struct {
    backupService input.BackupService
}

// NewBackupHandler cria um novo handler de backup
func NewBackupHandler(backupService input.BackupService) *BackupHandler {
    return &BackupHandler{
        backupService: backupService,
    }
}

// CreateBackupCommand cria o comando de backup
func (h *BackupHandler) CreateBackupCommand() *cobra.Command {
    var (
        name      string
        files     []string
        backupDir string
        compress  bool
    )
    
    cmd := &cobra.Command{
        Use:   "create",
        Short: "Create a backup of configuration files",
        Long: `Create a backup of specified configuration files.
        
Examples:
  whiterose backup create --name "config-backup" --files ".env,.config.json" --backup-dir "./backups"
  whiterose backup create -n "daily" -f ".env" -f ".config.json" -d "./backups" --compress`,
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            
            req := &dto.BackupRequest{
                Name:      name,
                FilePaths: files,
                BackupDir: backupDir,
                Compress:  compress,
            }
            
            response, err := h.backupService.CreateBackup(ctx, req)
            if err != nil {
                return fmt.Errorf("failed to create backup: %w", err)
            }
            
            h.printBackupResult(response)
            return nil
        },
    }
    
    cmd.Flags().StringVarP(&name, "name", "n", "", "Name for the backup (required)")
    cmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "Files to backup (can be specified multiple times)")
    cmd.Flags().StringVarP(&backupDir, "backup-dir", "d", "./backups", "Directory to store backups")
    cmd.Flags().BoolVar(&compress, "compress", false, "Compress backup files")
    
    cmd.MarkFlagRequired("name")
    cmd.MarkFlagRequired("files")
    
    return cmd
}

// ListBackupsCommand cria o comando de listagem de backups
func (h *BackupHandler) ListBackupsCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "list",
        Short: "List all available backups",
        Long:  "List all available backups with their details.",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            
            response, err := h.backupService.ListBackups(ctx)
            if err != nil {
                return fmt.Errorf("failed to list backups: %w", err)
            }
            
            h.printBackupList(response)
            return nil
        },
    }
    
    return cmd
}

// RestoreBackupCommand cria o comando de restore
func (h *BackupHandler) RestoreBackupCommand() *cobra.Command {
    var (
        backupID    string
        restorePath string
        overwrite   bool
    )
    
    cmd := &cobra.Command{
        Use:   "restore",
        Short: "Restore a backup",
        Long: `Restore a backup to original location or specified path.
        
Examples:
  whiterose backup restore --id "backup_123" 
  whiterose backup restore --id "backup_123" --restore-path "/new/location" --overwrite`,
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            
            req := &dto.RestoreRequest{
                BackupID:    backupID,
                RestorePath: restorePath,
                Overwrite:   overwrite,
            }
            
            response, err := h.backupService.RestoreBackup(ctx, req)
            if err != nil {
                return fmt.Errorf("failed to restore backup: %w", err)
            }
            
            h.printRestoreResult(response)
            return nil
        },
    }
    
    cmd.Flags().StringVar(&backupID, "id", "", "Backup ID to restore (required)")
    cmd.Flags().StringVar(&restorePath, "restore-path", "", "Path to restore backup (default: original location)")
    cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing files")
    
    cmd.MarkFlagRequired("id")
    
    return cmd
}

// DeleteBackupCommand cria o comando de remoção
func (h *BackupHandler) DeleteBackupCommand() *cobra.Command {
    var backupID string
    
    cmd := &cobra.Command{
        Use:   "delete",
        Short: "Delete a backup",
        Long:  "Delete a backup by its ID.",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx := context.Background()
            
            err := h.backupService.DeleteBackup(ctx, backupID)
            if err != nil {
                return fmt.Errorf("failed to delete backup: %w", err)
            }
            
            fmt.Printf("✅ Backup %s deleted successfully\n", backupID)
            return nil
        },
    }
    
    cmd.Flags().StringVar(&backupID, "id", "", "Backup ID to delete (required)")
    cmd.MarkFlagRequired("id")
    
    return cmd
}

// printBackupResult imprime o resultado do backup
func (h *BackupHandler) printBackupResult(response *dto.BackupResponse) {
    fmt.Printf("🎯 Backup Created Successfully\n")
    fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
    fmt.Printf("📋 ID: %s\n", response.ID)
    fmt.Printf("📝 Name: %s\n", response.Name)
    fmt.Printf("📁 Files: %s\n", strings.Join(response.FilePaths, ", "))
    fmt.Printf("💾 Backup Path: %s\n", response.BackupPath)
    fmt.Printf("📏 Size: %d bytes\n", response.Size)
    fmt.Printf("🕒 Created: %s\n", response.CreatedAt.Format("2006-01-02 15:04:05"))
    fmt.Printf("✅ Status: %s\n", response.Message)
    fmt.Printf("\n")
}

// printBackupList imprime a lista de backups
func (h *BackupHandler) printBackupList(response *dto.ListBackupsResponse) {
    fmt.Printf("📋 Available Backups (%d total)\n", response.Total)
    fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
    
    if len(response.Backups) == 0 {
        fmt.Printf("No backups found.\n")
        return
    }
    
    for i, backup := range response.Backups {
        fmt.Printf("%d. %s\n", i+1, backup.Name)
        fmt.Printf("   ID: %s\n", backup.ID)
        fmt.Printf("   Size: %d bytes\n", backup.Size)
        fmt.Printf("   Created: %s\n", backup.CreatedAt.Format("2006-01-02 15:04:05"))
        fmt.Printf("\n")
    }
}

// printRestoreResult imprime o resultado do restore
func (h *BackupHandler) printRestoreResult(response *dto.RestoreResponse) {
    fmt.Printf("🔄 Backup Restored Successfully\n")
    fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
    fmt.Printf("📋 Backup ID: %s\n", response.BackupID)
    fmt.Printf("📁 Restored to: %s\n", response.RestorePath)
    fmt.Printf("✅ Status: %s\n", response.Message)
    fmt.Printf("\n")
}
```

#### 📍 4.2 Adicionar Command ao Root

**Arquivo**: `cmd/backup.go`

```go
package cmd

import (
    "github.com/spf13/cobra"
    
    "github.com/fabianoflorentino/whiterose/cmd/handlers"
    "github.com/fabianoflorentino/whiterose/internal/infrastructure/config"
)

// backupCmd representa o comando de backup
var backupCmd = &cobra.Command{
    Use:   "backup",
    Short: "Manage configuration backups",
    Long: `Manage configuration backups including creation, listing, restoration and deletion.
    
The backup command helps you create and manage backups of your configuration files
to ensure you can recover from configuration changes or corruption.`,
}

func init() {
    rootCmd.AddCommand(backupCmd)
    
    // Configurar dependency injection
    container := config.NewContainer()
    backupHandler := handlers.NewBackupHandler(container.GetBackupService())
    
    // Adicionar subcomandos
    backupCmd.AddCommand(backupHandler.CreateBackupCommand())
    backupCmd.AddCommand(backupHandler.ListBackupsCommand())
    backupCmd.AddCommand(backupHandler.RestoreBackupCommand())
    backupCmd.AddCommand(backupHandler.DeleteBackupCommand())
}
```

### 5️⃣ **Dependency Injection** (Quinto - Conectar tudo)

#### 📍 5.1 Atualizar Container

**Arquivo**: `internal/infrastructure/config/container.go` (adicionar)

```go
// Adicionar ao container existente

func (c *Container) wireBackupComponents() {
    // Repository
    c.backupRepo = backup.NewMemoryBackupRepository()
    
    // Adapters
    c.backupPort = backup.NewFilesystemBackupAdapter()
    
    // Domain Services
    c.backupDomainService = services.NewBackupDomainService()
    
    // Use Cases
    c.backupService = usecases.NewBackupUseCase(
        c.backupRepo,
        c.backupPort,
        c.backupDomainService,
        c.logger,
    )
}

func (c *Container) GetBackupService() input.BackupService {
    return c.backupService
}
```

### 6️⃣ **Testes** (Sexto - Garantir qualidade)

#### 📍 6.1 Testes de Entidade

**Arquivo**: `internal/domain/entities/backup_test.go`

```go
package entities_test

import (
    "testing"
    
    "github.com/fabianoflorentino/whiterose/internal/domain/entities"
)

func TestNewBackup(t *testing.T) {
    tests := []struct {
        name        string
        backupName  string
        filePath    string
        backupPath  string
        expectError bool
    }{
        {
            name:       "valid backup",
            backupName: "test-backup",
            filePath:   "/path/to/file.txt",
            backupPath: "/backup/path/file.txt",
            expectError: false,
        },
        {
            name:       "empty backup name",
            backupName: "",
            filePath:   "/path/to/file.txt",
            backupPath: "/backup/path/file.txt",
            expectError: true,
        },
        // ... mais casos de teste
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            backup, err := entities.NewBackup(tt.backupName, tt.filePath, tt.backupPath)
            
            if tt.expectError {
                if err == nil {
                    t.Errorf("expected error but got none")
                }
                return
            }
            
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            
            if backup.Name != tt.backupName {
                t.Errorf("expected name %s, got %s", tt.backupName, backup.Name)
            }
        })
    }
}
```

#### 📍 6.2 Testes de Use Case

**Arquivo**: `internal/application/usecases/backup_usecase_test.go`

```go
package usecases_test

import (
    "context"
    "testing"
    
    "github.com/fabianoflorentino/whiterose/internal/application/dto"
    "github.com/fabianoflorentino/whiterose/internal/application/usecases"
    // ... mocks
)

func TestBackupUseCase_CreateBackup(t *testing.T) {
    // Configurar mocks
    mockRepo := &MockBackupRepository{}
    mockPort := &MockBackupPort{}
    mockService := &MockBackupDomainService{}
    mockLogger := &MockLogger{}
    
    useCase := usecases.NewBackupUseCase(mockRepo, mockPort, mockService, mockLogger)
    
    // Configurar expectativas dos mocks
    mockPort.On("FileExists", mock.Anything, "/path/to/file.txt").Return(true, nil)
    mockPort.On("CreateBackup", mock.Anything, mock.Anything, mock.Anything, false).Return(nil)
    // ... mais configurações
    
    // Executar teste
    req := &dto.BackupRequest{
        Name:      "test-backup",
        FilePaths: []string{"/path/to/file.txt"},
        BackupDir: "/backup/dir",
        Compress:  false,
    }
    
    response, err := useCase.CreateBackup(context.Background(), req)
    
    // Verificar resultados
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    
    if response == nil {
        t.Error("expected response but got nil")
    }
    
    // Verificar que mocks foram chamados
    mockRepo.AssertExpectations(t)
    mockPort.AssertExpectations(t)
}
```

## 🎯 Implementation Order Summary

### ✅ Implementation Checklist

1. **Domain Layer** (Solid foundation)
   - [ ] Create entity
   - [ ] Define repository interface
   - [ ] Implement domain service
   - [ ] Add domain errors

2. **Application Layer** (Orchestration)
   - [ ] Create DTOs
   - [ ] Define input ports
   - [ ] Define output ports
   - [ ] Implement use case

3. **Infrastructure Layer** (Adapters)
   - [ ] Implement adapters
   - [ ] Implement repositories
   - [ ] Configure logging

4. **CLI Interface** (User interface)
   - [ ] Create CLI handler
   - [ ] Add commands to root
   - [ ] Configure flags and validations

5. **Dependency Injection** (Connect)
   - [ ] Update container
   - [ ] Configure wirings
   - [ ] Expose services

6. **Tests** (Ensure quality)
   - [ ] Entity tests
   - [ ] Use case tests
   - [ ] Adapter tests
   - [ ] Integration tests

## 🚀 Final Commands

After implementing everything, the new commands will be available:

```bash
# Create backup
whiterose backup create --name "config-backup" --files ".env,.config.json" --backup-dir "./backups"

# List backups
whiterose backup list

# Restore backup
whiterose backup restore --id "backup_123"

# Delete backup
whiterose backup delete --id "backup_123"
```

Esta ordem garante que cada camada seja implementada na sequência correta, mantendo a arquitetura hexagonal e permitindo testes incrementais em cada etapa.
