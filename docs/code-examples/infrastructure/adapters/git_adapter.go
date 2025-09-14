package adapters

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/entities"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/errors"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/repositories"
)

// GitAdapter implements the GitRepository interface using go-git
type GitAdapter struct {
	// Can include configuration options like credentials, timeouts, etc.
}

// NewGitAdapter creates a new Git adapter
func NewGitAdapter() *GitAdapter {
	return &GitAdapter{}
}

// Clone clones a repository to the specified local path
func (g *GitAdapter) Clone(ctx context.Context, repo *entities.Repository, localPath string) error {
	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Clone options
	cloneOptions := &git.CloneOptions{
		URL:           repo.URL().String(),
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", repo.Branch())),
		SingleBranch:  true,
	}

	// Perform the clone with context
	_, err := git.PlainCloneContext(ctx, localPath, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository %s: %w", repo.Name(), err)
	}

	return nil
}

// Pull updates the local repository with remote changes
func (g *GitAdapter) Pull(ctx context.Context, localPath string) error {
	// Open the repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return fmt.Errorf("failed to open repository at %s: %w", localPath, err)
	}

	// Get the working directory
	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Pull with context
	err = workTree.PullContext(ctx, &git.PullOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull changes: %w", err)
	}

	return nil
}

// Checkout switches to the specified branch
func (g *GitAdapter) Checkout(ctx context.Context, localPath, branch string) error {
	// Open the repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return fmt.Errorf("failed to open repository at %s: %w", localPath, err)
	}

	// Get the working directory
	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Checkout the branch
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	})
	if err != nil {
		return fmt.Errorf("failed to checkout branch %s: %w", branch, err)
	}

	return nil
}

// GetCurrentBranch returns the current branch name
func (g *GitAdapter) GetCurrentBranch(ctx context.Context, localPath string) (string, error) {
	// Open the repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository at %s: %w", localPath, err)
	}

	// Get the HEAD reference
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Extract branch name from reference
	if head.Name().IsBranch() {
		return head.Name().Short(), nil
	}

	return "", errors.NewBusinessRuleError("HEAD is not pointing to a branch")
}

// ListBranches returns all available branches
func (g *GitAdapter) ListBranches(ctx context.Context, localPath string) ([]string, error) {
	// Open the repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository at %s: %w", localPath, err)
	}

	// Get all references
	refs, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to get references: %w", err)
	}

	var branches []string
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			branches = append(branches, ref.Name().Short())
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate references: %w", err)
	}

	return branches, nil
}

// IsClean checks if the repository has uncommitted changes
func (g *GitAdapter) IsClean(ctx context.Context, localPath string) (bool, error) {
	// Open the repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return false, fmt.Errorf("failed to open repository at %s: %w", localPath, err)
	}

	// Get the working directory
	workTree, err := repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get the status
	status, err := workTree.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get status: %w", err)
	}

	// Repository is clean if status is empty
	return status.IsClean(), nil
}

// GetLastCommit returns information about the last commit
func (g *GitAdapter) GetLastCommit(ctx context.Context, localPath string) (*repositories.CommitInfo, error) {
	// Open the repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository at %s: %w", localPath, err)
	}

	// Get the HEAD reference
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Get the commit object
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	return &repositories.CommitInfo{
		Hash:      commit.Hash.String(),
		Message:   commit.Message,
		Author:    commit.Author.Name,
		Email:     commit.Author.Email,
		Timestamp: commit.Author.When.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// Compile-time check to ensure GitAdapter implements GitRepository
var _ repositories.GitRepository = (*GitAdapter)(nil)
