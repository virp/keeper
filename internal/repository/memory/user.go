package memory

import (
	"context"
	"sync"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

type UserRepository struct {
	mu    *sync.RWMutex
	users map[string]entity.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		mu:    new(sync.RWMutex),
		users: map[string]entity.User{},
	}
}

func (r *UserRepository) Create(_ context.Context, user entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.Login]; ok {
		return repository.ErrUserAlreadyExist
	}

	r.users[user.Login] = user

	return nil
}

func (r *UserRepository) GetByLogin(_ context.Context, login string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, ok := r.users[login]; ok {
		return user, nil
	}
	return entity.User{}, repository.ErrUserNotFound
}
