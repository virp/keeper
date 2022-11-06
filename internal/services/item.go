package services

import (
	"context"
	"time"

	"keeper/internal/entity"
)

type ItemService struct {
	idGenerator    IdGenerator
	itemRepository ItemRepository
}

func NewItemService(idGenerator IdGenerator, itemRepository ItemRepository) *ItemService {
	return &ItemService{
		idGenerator:    idGenerator,
		itemRepository: itemRepository,
	}
}

func (s *ItemService) Create(ctx context.Context, userID string, item Item) error {
	createdItem := itemServiceToItemEntity(item, s.idGenerator.Generate(), userID)
	createdItem.CreatedAt = time.Now()
	createdItem.UpdatedAt = createdItem.CreatedAt
	return s.itemRepository.Create(ctx, createdItem)
}

func (s *ItemService) Update(ctx context.Context, userID string, item Item) error {
	existedItem, err := s.itemRepository.GetByUserIDAndName(ctx, userID, item.Name)
	if err != nil {
		return err
	}
	updatedItem := itemServiceToItemEntity(item, existedItem.ID, userID)
	updatedItem.UpdatedAt = time.Now()
	return s.itemRepository.Update(ctx, updatedItem)
}

func (s *ItemService) Get(ctx context.Context, userID string, name string) (Item, error) {
	item, err := s.itemRepository.GetByUserIDAndName(ctx, userID, name)
	if err != nil {
		return Item{}, err
	}
	return itemEntityToItemService(item), nil
}

func (s *ItemService) Delete(ctx context.Context, userID string, name string) error {
	item, err := s.itemRepository.GetByUserIDAndName(ctx, userID, name)
	if err != nil {
		return err
	}
	return s.itemRepository.Delete(ctx, item)
}

func (s *ItemService) List(ctx context.Context, userID string) ([]Item, error) {
	items, err := s.itemRepository.FindByUser(ctx, userID)
	if err != nil {
		return []Item{}, err
	}
	list := make([]Item, 0, len(items))
	for _, item := range items {
		list = append(list, itemEntityToItemService(item))
	}
	return list, nil
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
