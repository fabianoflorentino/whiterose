// Package entities defines types and functions related to the system's domain entities,
// including source code repositories and authentication methods.
package entities

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// AuthType represents the type of authentication used to access a repository.
type AuthType string

const (
	AuthTypeSSH   AuthType = "ssh"
	AuthTypeHTTPS AuthType = "https"
)

type Repository struct {
	ID         string
	URL        string
	Name       string
	Directory  string
	Branch     string
	AuthMethod AuthenticationMethod
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type SSHKeyConfig struct {
	Path       string
	Name       string
	Passphrase string
}

type AuthenticationMethod struct {
	Type     AuthType
	Username string
	Token    string
	SSHKey   SSHKeyConfig
}

// NewRepository creates a new Repository instance from a URL and local directory.
// It validates the URL and sets default values for branch and authentication.
func NewRepository(url, directory string) (*Repository, error) {
	if err := validateURL(url); err != nil {
		return nil, fmt.Errorf("invalid repository URL: %w", err)
	}

	if directory == "" {
		return nil, fmt.Errorf("directorty cannot be empty")
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

// SetBranch updates the branch name of the repository.
func (r *Repository) SetBranch(branch string) error {
	if branch == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	if err := r.Validate(); err != nil {
		return fmt.Errorf("cannot set branch on invalid repository: %w", err)
	}

	r.Branch = branch
	r.UpdatedAt = time.Now()

	return nil
}

// SetAuthentication updates the authentication method for the repository.
func (r *Repository) SetAuthentication(auth AuthenticationMethod) error {
	if err := auth.Validate(); err != nil {
		return fmt.Errorf("invalid authentication method: %w", err)
	}

	r.AuthMethod = auth
	r.UpdatedAt = time.Now()

	return nil
}

// IsSSH returns true if the repository uses SSH authentication.
func (r *Repository) IsSSH() bool {
	return r.AuthMethod.Type == AuthTypeSSH
}

// IsHTTPS returns true if the repository uses HTTPS authentication.
func (r *Repository) IsHTTPS() bool {
	return r.AuthMethod.Type == AuthTypeHTTPS
}

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

// Validate checks if the authentication method is valid based on its type.
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

// private helper functions

// validateURL checks if the repository URL is valid and supported (SSH or HTTPS).
func validateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if strings.HasPrefix(rawURL, "git@") {
		return nil
	}

	if strings.HasPrefix(rawURL, "https://") {
		_, err := url.Parse(rawURL)
		return err
	}

	return fmt.Errorf("unsupported URL format")
}

// extractRepositoryName extracts the repository name from the given URL.
func extractRepositoryName(rawURL string) string {
	if strings.HasPrefix(rawURL, "git@") {
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

// determineAuthMethod determines the appropriate authentication method based on the repository URL.
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

// generateID generates a unique identifier for the repository based on the URL and current timestamp.
func generateID(url string) string {
	return fmt.Sprintf("repo_%d", len(url)*int(time.Now().Unix())%10000)
}
