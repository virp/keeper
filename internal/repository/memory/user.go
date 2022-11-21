package memory

import (
	"context"
	"sync"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

// UserRepository in memory user storage.
type UserRepository struct {
	mu    *sync.RWMutex
	users map[string]entity.User
}

// NewUserRepository construct UserRepository.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		mu:    new(sync.RWMutex),
		users: map[string]entity.User{},
	}
}

// Create store user entity in storage.
func (r *UserRepository) Create(_ context.Context, user entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.Login]; ok {
		return repository.ErrUserAlreadyExist
	}

	r.users[user.Login] = user

	return nil
}

// GetByLogin return user by login from storage.
func (r *UserRepository) GetByLogin(_ context.Context, login string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, ok := r.users[login]; ok {
		return user, nil
	}
	return entity.User{}, repository.ErrUserNotFound
}
