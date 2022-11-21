package memory

import (
	"context"
	"sync"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

// ItemRepository in memory item storage.
type ItemRepository struct {
	mu    *sync.RWMutex
	items map[string]map[string]entity.Item
}

// NewItemRepository construct ItemRepository.
func NewItemRepository() *ItemRepository {
	return &ItemRepository{
		mu:    new(sync.RWMutex),
		items: map[string]map[string]entity.Item{},
	}
}

// Create store item in storage.
func (r *ItemRepository) Create(_ context.Context, item entity.Item) error {
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

// Update store new version item in storage.
func (r *ItemRepository) Update(_ context.Context, item entity.Item) error {
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

// GetByUserIDAndName return item from storage by user ID and item name.
func (r *ItemRepository) GetByUserIDAndName(_ context.Context, userID string, name string) (entity.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if item, ok := r.items[userID][name]; ok {
		return item, nil
	}

	return entity.Item{}, repository.ErrItemNotFound
}

// Delete remove item from storage.
func (r *ItemRepository) Delete(_ context.Context, item entity.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[item.UserID][item.Name]; !ok {
		return repository.ErrItemNotFound
	}
	delete(r.items[item.UserID], item.Name)

	return nil
}

// FindByUser return item names list from storage for user ID.
func (r *ItemRepository) FindByUser(_ context.Context, userID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if items, ok := r.items[userID]; ok {
		list := make([]string, 0, len(items))
		for _, item := range items {
			list = append(list, item.Name)
		}
		return list, nil
	}

	return []string{}, nil
}
