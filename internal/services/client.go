package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"google.golang.org/grpc/metadata"

	pb "keeper/gen/service"
)

// ClientService implement logic for Keeper client application and work with GRPC server calls.
type ClientService struct {
	client pb.KeeperServiceClient
}

// NewClientService construct new ClientService.
func NewClientService(client pb.KeeperServiceClient) *ClientService {
	return &ClientService{
		client: client,
	}
}

// Register make Register rpc call.
func (s *ClientService) Register(ctx context.Context, login string, password string) (string, error) {
	req := pb.LoginRequest{
		Login:    login,
		Password: password,
	}
	res, err := s.client.Register(ctx, &req)
	if err != nil {
		return "", err
	}
	token := res.GetToken()
	return token, nil
}

// Login make Login rpc call.
func (s *ClientService) Login(ctx context.Context, login string, password string) (string, error) {
	req := pb.LoginRequest{
		Login:    login,
		Password: password,
	}
	res, err := s.client.Login(ctx, &req)
	if err != nil {
		return "", err
	}
	token := res.GetToken()
	return token, nil
}

// List make GetItemsList rpc call.
func (s *ClientService) List(ctx context.Context, token string) ([]string, error) {
	ctx = getOutgoingContext(ctx, token)
	req := pb.GetItemsListRequest{}
	res, err := s.client.GetItemsList(ctx, &req)
	if err != nil {
		return nil, err
	}
	items := res.GetNames()
	return items, nil
}

// Get make GetItem rpc call.
func (s *ClientService) Get(ctx context.Context, token string, secret string, name string) (Item, error) {
	ctx = getOutgoingContext(ctx, token)
	req := pb.GetItemRequest{
		Name: name,
	}
	res, err := s.client.GetItem(ctx, &req)
	if err != nil {
		return Item{}, err
	}
	pbItem := res.GetItem()
	data, err := decryptData(secret, pbItem.GetData())
	if err != nil {
		return Item{}, err
	}
	item := Item{
		Name: pbItem.GetName(),
		Type: pbItem.GetType(),
		Data: data,
	}
	item.Metadata = make([]Metadata, 0, len(pbItem.GetMetadata()))
	for _, m := range pbItem.GetMetadata() {
		mtd := Metadata{
			Key:   m.GetKey(),
			Value: m.GetValue(),
		}
		item.Metadata = append(item.Metadata, mtd)
	}
	return item, nil
}

// Add make CreateItem rpc call.
func (s *ClientService) Add(ctx context.Context, token string, secret string, item Item) error {
	ctx = getOutgoingContext(ctx, token)
	data, err := encryptData(secret, item.Data)
	if err != nil {
		return err
	}
	pbItem := pb.Item{
		Name: item.Name,
		Type: item.Type,
		Data: data,
	}
	for _, m := range item.Metadata {
		mtd := pb.Metadata{
			Key:   m.Key,
			Value: m.Value,
		}
		pbItem.Metadata = append(pbItem.Metadata, &mtd)
	}
	req := pb.CreateItemRequest{
		Item: &pbItem,
	}
	_, err = s.client.CreateItem(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}

// Delete make DeleteItem rpc call.
func (s *ClientService) Delete(ctx context.Context, token string, name string) error {
	ctx = getOutgoingContext(ctx, token)
	req := pb.DeleteItemRequest{
		Name: name,
	}
	_, err := s.client.DeleteItem(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}

func getOutgoingContext(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}

func encryptData(secret string, data []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	enc := gcm.Seal(nil, nonce, data, nil)
	enc = append(enc, nonce...)
	return enc, nil
}

func decryptData(secret string, data []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := data[len(data)-gcm.NonceSize():]
	dec, err := gcm.Open(nil, nonce, data[:len(data)-gcm.NonceSize()], nil)
	if err != nil {
		return nil, err
	}
	return dec, nil
}
