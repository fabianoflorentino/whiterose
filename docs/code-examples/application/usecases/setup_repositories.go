package usecases

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/entities"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/errors"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/repositories"
)

// SetupRepositoriesUseCase handles repository setup operations
type SetupRepositoriesUseCase struct {
	repositoryRepo repositories.RepositoryRepository
	gitRepo        repositories.GitRepository
	configRepo     repositories.ConfigurationRepository
	workingDir     string
}

// NewSetupRepositoriesUseCase creates a new use case instance
func NewSetupRepositoriesUseCase(
	repositoryRepo repositories.RepositoryRepository,
	gitRepo repositories.GitRepository,
	configRepo repositories.ConfigurationRepository,
	workingDir string,
) *SetupRepositoriesUseCase {
	return &SetupRepositoriesUseCase{
		repositoryRepo: repositoryRepo,
		gitRepo:        gitRepo,
		configRepo:     configRepo,
		workingDir:     workingDir,
	}
}

// SetupRepositoriesRequest represents the input for repository setup
type SetupRepositoriesRequest struct {
	Repositories []RepositorySetupData `json:"repositories"`
	ForceClone   bool                  `json:"force_clone"`
}

// RepositorySetupData represents data for setting up a repository
type RepositorySetupData struct {
	Name   string `json:"name" validate:"required"`
	URL    string `json:"url" validate:"required,url"`
	Branch string `json:"branch" validate:"required"`
}

// SetupRepositoriesResponse represents the output of repository setup
type SetupRepositoriesResponse struct {
	SetupResults []RepositorySetupResult `json:"setup_results"`
	TotalCount   int                     `json:"total_count"`
	SuccessCount int                     `json:"success_count"`
	FailureCount int                     `json:"failure_count"`
}

// RepositorySetupResult represents the result of setting up a single repository
type RepositorySetupResult struct {
	Name      string `json:"name"`
	Status    string `json:"status"` // "success", "failed", "skipped"
	Message   string `json:"message"`
	LocalPath string `json:"local_path,omitempty"`
	Error     string `json:"error,omitempty"`
}

// Execute performs the repository setup operation
func (uc *SetupRepositoriesUseCase) Execute(ctx context.Context, request SetupRepositoriesRequest) (*SetupRepositoriesResponse, error) {
	if err := uc.validateRequest(request); err != nil {
		return nil, err
	}

	response := &SetupRepositoriesResponse{
		SetupResults: make([]RepositorySetupResult, 0, len(request.Repositories)),
		TotalCount:   len(request.Repositories),
	}

	for _, repoData := range request.Repositories {
		result := uc.setupSingleRepository(ctx, repoData, request.ForceClone)
		response.SetupResults = append(response.SetupResults, result)

		if result.Status == "success" {
			response.SuccessCount++
		} else {
			response.FailureCount++
		}
	}

	return response, nil
}

// validateRequest validates the input request
func (uc *SetupRepositoriesUseCase) validateRequest(request SetupRepositoriesRequest) error {
	if len(request.Repositories) == 0 {
		return errors.NewValidationError("at least one repository must be specified", nil)
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
func (uc *SetupRepositoriesUseCase) setupSingleRepository(ctx context.Context, repoData RepositorySetupData, forceClone bool) RepositorySetupResult {
	// Create repository entity
	repo, err := entities.NewRepository(repoData.Name, repoData.URL, repoData.Branch)
	if err != nil {
		return RepositorySetupResult{
			Name:    repoData.Name,
			Status:  "failed",
			Message: "Failed to create repository entity",
			Error:   err.Error(),
		}
	}

	// Check if repository already exists
	existingRepo, err := uc.repositoryRepo.FindByName(ctx, repoData.Name)
	if err != nil && !errors.IsNotFoundError(err) {
		return RepositorySetupResult{
			Name:    repoData.Name,
			Status:  "failed",
			Message: "Failed to check existing repository",
			Error:   err.Error(),
		}
	}

	// Skip if repository exists and not forcing clone
	if existingRepo != nil && !forceClone {
		return RepositorySetupResult{
			Name:      repoData.Name,
			Status:    "skipped",
			Message:   "Repository already exists (use force_clone to override)",
			LocalPath: existingRepo.LocalPath(),
		}
	}

	// Set local path
	localPath := filepath.Join(uc.workingDir, repoData.Name)
	if err := repo.SetLocalPath(localPath); err != nil {
		return RepositorySetupResult{
			Name:    repoData.Name,
			Status:  "failed",
			Message: "Failed to set local path",
			Error:   err.Error(),
		}
	}

	// Clone repository
	if err := uc.gitRepo.Clone(ctx, repo, localPath); err != nil {
		return RepositorySetupResult{
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
		return RepositorySetupResult{
			Name:    repoData.Name,
			Status:  "failed",
			Message: "Failed to save repository",
			Error:   err.Error(),
		}
	}

	return RepositorySetupResult{
		Name:      repoData.Name,
		Status:    "success",
		Message:   "Repository successfully cloned and configured",
		LocalPath: localPath,
	}
}

// GetRepositoryStatus returns the current status of repositories
func (uc *SetupRepositoriesUseCase) GetRepositoryStatus(ctx context.Context) ([]RepositoryStatus, error) {
	repositories, err := uc.repositoryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repositories: %w", err)
	}

	statuses := make([]RepositoryStatus, 0, len(repositories))
	for _, repo := range repositories {
		status := RepositoryStatus{
			Name:        repo.Name(),
			URL:         repo.URL().String(),
			Branch:      repo.Branch(),
			LocalPath:   repo.LocalPath(),
			IsCloned:    repo.IsCloned(),
			LastUpdated: repo.LastUpdated().Format(time.RFC3339),
		}

		// Check if local path exists and get current branch
		if repo.IsCloned() {
			currentBranch, err := uc.gitRepo.GetCurrentBranch(ctx, repo.LocalPath())
			if err == nil {
				status.CurrentBranch = currentBranch
			}

			isClean, err := uc.gitRepo.IsClean(ctx, repo.LocalPath())
			if err == nil {
				status.IsClean = isClean
			}
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// RepositoryStatus represents the current status of a repository
type RepositoryStatus struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Branch        string `json:"target_branch"`
	CurrentBranch string `json:"current_branch"`
	LocalPath     string `json:"local_path"`
	IsCloned      bool   `json:"is_cloned"`
	IsClean       bool   `json:"is_clean"`
	LastUpdated   string `json:"last_updated"`
}
