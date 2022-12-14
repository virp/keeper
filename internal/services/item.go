package services

import (
	"context"
	"fmt"
	"time"

	"keeper/internal/entity"
)

const (
	itemNameMinLength = 1
)

// ItemService implement logic for working with items.
type ItemService struct {
	idGenerator    IdGenerator
	itemRepository ItemRepository
}

// NewItemService construct new ItemService.
func NewItemService(idGenerator IdGenerator, itemRepository ItemRepository) *ItemService {
	return &ItemService{
		idGenerator:    idGenerator,
		itemRepository: itemRepository,
	}
}

// Create save new item for user in storage.
func (s *ItemService) Create(ctx context.Context, userID string, item Item) error {
	var fields FieldErrors
	if len(item.Name) < itemNameMinLength {
		field := FieldError{
			Field: "item.name",
			Error: fmt.Sprintf("length should be greater or equal %d", itemNameMinLength),
		}
		fields = append(fields, field)
	}
	if len(fields) > 0 {
		return fields
	}

	createdItem := itemServiceToItemEntity(item, s.idGenerator.Generate(), userID)
	createdItem.CreatedAt = time.Now()
	createdItem.UpdatedAt = createdItem.CreatedAt
	return s.itemRepository.Create(ctx, createdItem)
}

// Update save new version of item for user in storage.
func (s *ItemService) Update(ctx context.Context, userID string, item Item) error {
	existedItem, err := s.itemRepository.GetByUserIDAndName(ctx, userID, item.Name)
	if err != nil {
		return err
	}
	updatedItem := itemServiceToItemEntity(item, existedItem.ID, userID)
	updatedItem.Name = existedItem.Name // Name can't be changed
	updatedItem.CreatedAt = existedItem.CreatedAt
	updatedItem.UpdatedAt = time.Now()
	return s.itemRepository.Update(ctx, updatedItem)
}

// Get receive user item from storage.
func (s *ItemService) Get(ctx context.Context, userID string, name string) (Item, error) {
	item, err := s.itemRepository.GetByUserIDAndName(ctx, userID, name)
	if err != nil {
		return Item{}, err
	}
	return itemEntityToItemService(item), nil
}

// Delete remove user item from storage.
func (s *ItemService) Delete(ctx context.Context, userID string, name string) error {
	item, err := s.itemRepository.GetByUserIDAndName(ctx, userID, name)
	if err != nil {
		return err
	}
	return s.itemRepository.Delete(ctx, item)
}

// List receive user items list from storage.
func (s *ItemService) List(ctx context.Context, userID string) ([]string, error) {
	items, err := s.itemRepository.FindByUser(ctx, userID)
	if err != nil {
		return []string{}, err
	}
	return items, nil
}

func itemServiceToItemEntity(in Item, id string, userID string) entity.Item {
	out := entity.Item{
		ID:     id,
		UserID: userID,
		Name:   in.Name,
		Type:   in.Type,
		Data:   in.Data,
	}
	for _, m := range in.Metadata {
		md := entity.Metadata{
			Key:   m.Key,
			Value: m.Value,
		}
		out.Metadata = append(out.Metadata, md)
	}
	return out
}

func itemEntityToItemService(in entity.Item) Item {
	out := Item{
		Name: in.Name,
		Type: in.Type,
		Data: in.Data,
	}
	for _, m := range in.Metadata {
		md := Metadata{
			Key:   m.Key,
			Value: m.Value,
		}
		out.Metadata = append(out.Metadata, md)
	}
	return out
}
