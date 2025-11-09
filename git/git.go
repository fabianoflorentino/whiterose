// Package git provides utilities for cloning and managing Git repositories,
// supporting both HTTPS and SSH authentication methods.
//
// Types:
//   - GitCloneOptions: Options for cloning a Git repository, including URL, directory, credentials, and SSH key information.
//
// Functions:
//   - FetchRepositories: Clones multiple repositories based on provided options.
//   - LoadRepositoriesFromFile: Loads repository clone options from a JSON file.
//   - clone: Clones a single repository and checks out the 'development' branch or creates a user-specific branch if not present.
//   - createSSHAuth: Creates SSH authentication using a private key file, with support for default key locations and names.
package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fabianoflorentino/whiterose/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

// GitCloneOptions holds options for cloning a Git repository, including URL, directory, credentials, and SSH key information.
type GitCloneOptions struct {
	URL        string
	Directory  string
	Username   string
	Password   string
	SSHKeyPath string
	SSHKeyName string
}

// NewGitRepository creates and returns a new GitCloneOptions instance.
func NewGitRepository() *GitCloneOptions {
	return &GitCloneOptions{}
}

// Setup loads repository configuration, sets authentication options from environment variables, and clones repositories.
func (g *GitCloneOptions) Setup() {
	cfg := g.loadConfigFile(filepath.Base(os.Getenv("CONFIG_FILE")))
	repos, err := LoadRepositoriesFromFile(cfg)
	if err != nil {
		fmt.Printf("failed to load repositories: %v", err)
		os.Exit(1)
	}

	for i := range repos {
		repos[i].Username = utils.GetEnvOrDefault("GIT_USER", "")
		repos[i].Password = utils.GetEnvOrDefault("GIT_TOKEN", "")
		repos[i].SSHKeyPath = utils.GetEnvOrDefault("SSH_KEY_PATH", "")
		repos[i].SSHKeyName = utils.GetEnvOrDefault("SSH_KEY_NAME", "id_rsa")
	}

	if err := g.fetchRepositories(repos); err != nil {
		fmt.Printf("failed to fetch repositories: %v", err)
		os.Exit(1)
	}
}

// fetchRepositories clones multiple repositories based on the provided options.
func (g *GitCloneOptions) fetchRepositories(repos []GitCloneOptions) error {
	for _, opts := range repos {
		fmt.Printf("Cloning %s into %s...\n", opts.URL, opts.Directory)
		if err := clone(opts); err != nil {
			return fmt.Errorf("error cloning %s: %w", opts.URL, err)
		}
	}

	return nil
}

// LoadRepositoriesFromFile loads repository clone options from a configuration file and returns a slice of GitCloneOptions.
func LoadRepositoriesFromFile(file string) ([]GitCloneOptions, error) {
	repoInfos, err := utils.FetchRepositories(file)
	if err != nil {
		return nil, err
	}
	var opts []GitCloneOptions
	for _, r := range repoInfos {
		opts = append(opts, GitCloneOptions{
			URL:       r.URL,
			Directory: r.Directory,
			// Username, Password, SSHKeyPath, SSHKeyName can be set later or via env
		})
	}
	return opts, nil
}

// clone clones a single Git repository into the specified directory, checks out the 'development' branch, or creates a user-specific branch if not present.
func clone(opts GitCloneOptions) error {
	if _, err := os.Stat(opts.Directory); err == nil {
		return fmt.Errorf("directory %s already exists", opts.Directory)
	}

	cloneOpts := &git.CloneOptions{
		URL:      opts.URL,
		Progress: os.Stdout,
	}

	if strings.HasPrefix(opts.URL, "https://") {
		cloneOpts.Auth = &http.BasicAuth{
			Username: opts.Username,
			Password: opts.Password,
		}
	} else if strings.HasPrefix(opts.URL, "git@") || strings.HasPrefix(opts.URL, "ssh://") {
		auth, err := createSSHAuth(opts.SSHKeyPath)
		if err != nil {
			return fmt.Errorf("failed to create SSH auth: %w", err)
		}
		cloneOpts.Auth = auth
	}

	fmt.Println("Cloning repository...")
	repo, err := git.PlainClone(opts.Directory, false, cloneOpts)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// After cloning, try to checkout the development branch
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/development",
	})
	if err == nil {
		fmt.Println("Checked out to development branch.")
		return nil
	}

	// If it does not exist, create the local branch development/<user_name>
	newBranch := fmt.Sprintf("development/%s", os.Getenv("USER"))
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + newBranch),
		Create: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create and checkout branch %s: %w", newBranch, err)
	}
	fmt.Printf("Created and checked out to branch %s.\n", newBranch)

	return nil
}

// createSSHAuth creates SSH authentication using a private key file, supporting default key locations and names.
func createSSHAuth(keyPath string) (*ssh.PublicKeys, error) {
	if keyPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		keyName := utils.GetEnvOrDefault("SSH_KEY_NAME", "id_rsa")
		keyPath = filepath.Join(homeDir, ".ssh", keyName)

	} else {
		fi, err := os.Stat(keyPath)
		if err == nil && fi.IsDir() {
			keyName := utils.GetEnvOrDefault("SSH_KEY_NAME", "id_rsa")
			keyPath = filepath.Join(keyPath, keyName)
		}
	}

	sshKey, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key file: %w", err)
	}

	publickeys, err := ssh.NewPublicKeys("git", sshKey, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH public keys: %w", err)
	}

	return publickeys, nil
}

// loadConfigFile determines the configuration file path based on file extension and existence of YAML/YML files in the user's home directory.
func (g *GitCloneOptions) loadConfigFile(f string) string {
	if utils.YmlOrYamlExistsInHomeDir() {
		if strings.HasSuffix(f, ".yaml") {
			return utils.GetEnvOrDefault("CONFIG_FILE", os.Getenv("HOME")+"/.config.yaml")
		}

		return utils.GetEnvOrDefault("CONFIG_FILE", os.Getenv("HOME")+"/.config.yml")
	}

	return utils.GetEnvOrDefault("CONFIG_FILE", os.Getenv("HOME")+"/.config.json")
}
