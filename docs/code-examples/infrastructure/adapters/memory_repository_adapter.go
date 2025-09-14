package adapters

import (
	"context"
	"sync"

	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/entities"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/errors"
	"github.com/fabianoflorentino/whiterose/docs/code-examples/domain/repositories"
)

// InMemoryRepositoryAdapter implements RepositoryRepository interface for testing/examples
type InMemoryRepositoryAdapter struct {
	repositories map[string]*entities.Repository
	mutex        sync.RWMutex
}

// NewInMemoryRepositoryAdapter creates a new in-memory repository adapter
func NewInMemoryRepositoryAdapter() *InMemoryRepositoryAdapter {
	return &InMemoryRepositoryAdapter{
		repositories: make(map[string]*entities.Repository),
	}
}

// Save persists a repository entity
func (r *InMemoryRepositoryAdapter) Save(ctx context.Context, repo *entities.Repository) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.repositories[repo.ID()] = repo
	return nil
}

// FindByID retrieves a repository by its ID
func (r *InMemoryRepositoryAdapter) FindByID(ctx context.Context, id string) (*entities.Repository, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	repo, exists := r.repositories[id]
	if !exists {
		return nil, errors.NewNotFoundError("repository not found")
	}

	return repo, nil
}

// FindByName retrieves a repository by its name
func (r *InMemoryRepositoryAdapter) FindByName(ctx context.Context, name string) (*entities.Repository, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, repo := range r.repositories {
		if repo.Name() == name {
			return repo, nil
		}
	}

	return nil, errors.NewNotFoundError("repository not found")
}

// FindAll retrieves all repositories
func (r *InMemoryRepositoryAdapter) FindAll(ctx context.Context) ([]*entities.Repository, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	repos := make([]*entities.Repository, 0, len(r.repositories))
	for _, repo := range r.repositories {
		repos = append(repos, repo)
	}

	return repos, nil
}

// Update updates an existing repository
func (r *InMemoryRepositoryAdapter) Update(ctx context.Context, repo *entities.Repository) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.repositories[repo.ID()]; !exists {
		return errors.NewNotFoundError("repository not found")
	}

	r.repositories[repo.ID()] = repo
	return nil
}

// Delete removes a repository by ID
func (r *InMemoryRepositoryAdapter) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.repositories[id]; !exists {
		return errors.NewNotFoundError("repository not found")
	}

	delete(r.repositories, id)
	return nil
}

// Exists checks if a repository exists by name
func (r *InMemoryRepositoryAdapter) Exists(ctx context.Context, name string) (bool, error) {
	_, err := r.FindByName(ctx, name)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Compile-time check to ensure InMemoryRepositoryAdapter implements RepositoryRepository
var _ repositories.RepositoryRepository = (*InMemoryRepositoryAdapter)(nil)
