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

type GitCloneOptions struct {
	URL        string
	Directory  string
	Username   string
	Password   string
	SSHKeyPath string
	SSHKeyName string
}

func FetchRepositories(repos []GitCloneOptions) error {
	for _, opts := range repos {
		fmt.Printf("Cloning %s into %s...\n", opts.URL, opts.Directory)
		if err := clone(opts); err != nil {
			return fmt.Errorf("error cloning %s: %w", opts.URL, err)
		}
	}

	return nil
}

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

	// Após o clone, tente o checkout da branch development
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

	// Se não existe, cria a branch local development/<nome_do_usuario>
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
