package memory

import (
	"context"
	"sync"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

type UserMemoryRepository struct {
	mu     *sync.RWMutex
	users  map[string]entity.User
	logins map[string]string
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		mu:     new(sync.RWMutex),
		users:  map[string]entity.User{},
		logins: map[string]string{},
	}
}

func (r *UserMemoryRepository) Create(_ context.Context, user entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; ok {
		return repository.ErrUserAlreadyExist
	}
	if _, ok := r.logins[user.Login]; ok {
		return repository.ErrUserAlreadyExist
	}

	r.users[user.ID] = user
	r.logins[user.Login] = user.ID

	return nil
}

func (r *UserMemoryRepository) GetByLogin(_ context.Context, login string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, ok := r.logins[login]
	if !ok {
		return entity.User{}, repository.ErrUserNotFound
	}
	if user, ok := r.users[userID]; ok {
		return user, nil
	}
	return entity.User{}, repository.ErrUserNotFound
}

func (r *UserMemoryRepository) GetByID(_ context.Context, id string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, ok := r.users[id]; ok {
		return user, nil
	}

	return entity.User{}, repository.ErrUserNotFound
}
