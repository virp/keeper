package server

import (
	"context"

	pb "keeper/gen/service"
	"keeper/internal/services"
)

func (s *KeeperServer) CreateItem(ctx context.Context, in *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	userID := getUserIDFromContext(ctx)
	item := itemMessageToItemEntity(in.GetItem())
	if err := s.itemService.Create(ctx, userID, item); err != nil {
		return nil, err
	}

	var response pb.CreateItemResponse
	return &response, nil
}

func (s *KeeperServer) UpdateItem(ctx context.Context, in *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	userID := getUserIDFromContext(ctx)
	item := itemMessageToItemEntity(in.GetItem())
	if err := s.itemService.Update(ctx, userID, item); err != nil {
		return nil, err
	}

	var response pb.UpdateItemResponse
	return &response, nil
}

func (s *KeeperServer) GetItem(ctx context.Context, in *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	userID := getUserIDFromContext(ctx)
	item, err := s.itemService.Get(ctx, userID, in.GetName())
	if err != nil {
		return nil, err
	}

	response := pb.GetItemResponse{
		Item: itemEntityToItemMessage(item),
	}
	return &response, nil
}

func (s *KeeperServer) DeleteItem(ctx context.Context, in *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	userID := getUserIDFromContext(ctx)
	if err := s.itemService.Delete(ctx, userID, in.GetName()); err != nil {
		return nil, err
	}

	var response pb.DeleteItemResponse
	return &response, nil
}

func (s *KeeperServer) GetItemsList(ctx context.Context, _ *pb.GetItemsListRequest) (*pb.GetItemsListResponse, error) {
	userID := getUserIDFromContext(ctx)
	items, err := s.itemService.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	response := pb.GetItemsListResponse{
		Names: items,
	}
	return &response, nil
}

func itemMessageToItemEntity(msg *pb.Item) services.Item {
	item := services.Item{
		Name: msg.GetName(),
		Type: msg.GetType(),
		Data: msg.GetData(),
	}
	for _, m := range msg.GetMetadata() {
		metadata := services.Metadata{
			Key:   m.GetKey(),
			Value: m.GetValue(),
		}
		item.Metadata = append(item.Metadata, metadata)
	}
	return item
}

func itemEntityToItemMessage(item services.Item) *pb.Item {
	msg := pb.Item{
		Name: item.Name,
		Type: item.Type,
		Data: item.Data,
	}
	for _, m := range item.Metadata {
		metadata := pb.Metadata{
			Key:   m.Key,
			Value: m.Value,
		}
		msg.Metadata = append(msg.Metadata, &metadata)
	}
	return &msg
}
