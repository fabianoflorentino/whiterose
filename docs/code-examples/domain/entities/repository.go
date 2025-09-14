package entities

import (
	"net/url"
	"strings"
	"time"

	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/errors"
)

// Repository represents a Git repository in the domain
type Repository struct {
	id          string
	name        string
	url         *url.URL
	branch      string
	localPath   string
	isCloned    bool
	lastUpdated time.Time
}

// NewRepository creates a new Repository entity with validation
func NewRepository(name, urlStr, branch string) (*Repository, error) {
	if err := validateRepositoryName(name); err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.NewValidationError("invalid repository URL", err)
	}

	if err := validateRepositoryURL(parsedURL); err != nil {
		return nil, err
	}

	if err := validateBranchName(branch); err != nil {
		return nil, err
	}

	return &Repository{
		id:          generateRepositoryID(name),
		name:        name,
		url:         parsedURL,
		branch:      branch,
		localPath:   "",
		isCloned:    false,
		lastUpdated: time.Now(),
	}, nil
}

// ID returns the repository unique identifier
func (r *Repository) ID() string {
	return r.id
}

// Name returns the repository name
func (r *Repository) Name() string {
	return r.name
}

// URL returns the repository URL
func (r *Repository) URL() *url.URL {
	return r.url
}

// Branch returns the target branch
func (r *Repository) Branch() string {
	return r.branch
}

// LocalPath returns the local filesystem path
func (r *Repository) LocalPath() string {
	return r.localPath
}

// IsCloned returns whether the repository is cloned locally
func (r *Repository) IsCloned() bool {
	return r.isCloned
}

// LastUpdated returns the last update timestamp
func (r *Repository) LastUpdated() time.Time {
	return r.lastUpdated
}

// SetLocalPath sets the local filesystem path
func (r *Repository) SetLocalPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.NewValidationError("local path cannot be empty", nil)
	}
	r.localPath = path
	return nil
}

// MarkAsCloned marks the repository as successfully cloned
func (r *Repository) MarkAsCloned() {
	r.isCloned = true
	r.lastUpdated = time.Now()
}

// UpdateBranch changes the target branch
func (r *Repository) UpdateBranch(newBranch string) error {
	if err := validateBranchName(newBranch); err != nil {
		return err
	}
	r.branch = newBranch
	r.lastUpdated = time.Now()
	return nil
}

// Clone validates that the repository can be cloned
func (r *Repository) Clone() error {
	if r.isCloned {
		return errors.NewBusinessRuleError("repository is already cloned")
	}
	if r.localPath == "" {
		return errors.NewBusinessRuleError("local path must be set before cloning")
	}
	return nil
}

// Validation functions

func validateRepositoryName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.NewValidationError("repository name cannot be empty", nil)
	}
	if len(name) > 255 {
		return errors.NewValidationError("repository name too long (max 255 characters)", nil)
	}
	return nil
}

func validateRepositoryURL(u *url.URL) error {
	if u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "ssh" {
		return errors.NewValidationError("repository URL must use http, https, or ssh scheme", nil)
	}
	if u.Host == "" {
		return errors.NewValidationError("repository URL must have a host", nil)
	}
	return nil
}

func validateBranchName(branch string) error {
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return errors.NewValidationError("branch name cannot be empty", nil)
	}
	if strings.Contains(branch, " ") {
		return errors.NewValidationError("branch name cannot contain spaces", nil)
	}
	return nil
}

func generateRepositoryID(name string) string {
	// Simple ID generation - in real implementation, use UUID or similar
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
