package memory

import (
	"context"
	"sync"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

type ItemMemoryRepository struct {
	mu    *sync.RWMutex
	items map[string]map[string]entity.Item
}

func NewItemMemoryRepository() *ItemMemoryRepository {
	return &ItemMemoryRepository{
		mu:    new(sync.RWMutex),
		items: map[string]map[string]entity.Item{},
	}
}

func (r *ItemMemoryRepository) Create(_ context.Context, item entity.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[item.UserID]; !ok {
		r.items[item.UserID] = map[string]entity.Item{}
	}
	if _, ok := r.items[item.UserID][item.Name]; ok {
		return repository.ErrItemAlreadyExist
	}
	r.items[item.UserID][item.Name] = item

	return nil
}

func (r *ItemMemoryRepository) Update(_ context.Context, item entity.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[item.UserID]; !ok {
		return repository.ErrItemNotFound
	}
	if _, ok := r.items[item.UserID][item.Name]; ok {
		r.items[item.UserID][item.Name] = item
		return nil
	}

	return repository.ErrItemNotFound
}

func (r *ItemMemoryRepository) GetByUserIDAndName(_ context.Context, userID string, name string) (entity.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if item, ok := r.items[userID][name]; ok {
		return item, nil
	}

	return entity.Item{}, repository.ErrItemNotFound
}

func (r *ItemMemoryRepository) Delete(_ context.Context, item entity.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[item.UserID][item.Name]; !ok {
		return repository.ErrItemNotFound
	}
	delete(r.items[item.UserID], item.Name)

	return nil
}

func (r *ItemMemoryRepository) FindByUser(_ context.Context, userID string) ([]entity.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if items, ok := r.items[userID]; ok {
		list := make([]entity.Item, 0, len(items))
		for _, item := range items {
			list = append(list, item)
		}
		return list, nil
	}

	return []entity.Item{}, nil
}
