package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

// ItemRepository S3 item storage.
type ItemRepository struct {
	bucket string
	client *s3.Client
}

// NewItemRepository construct ItemRepository.
func NewItemRepository(client *s3.Client, bucket string) *ItemRepository {
	return &ItemRepository{
		bucket: bucket,
		client: client,
	}
}

// Create store item in storage.
func (r *ItemRepository) Create(ctx context.Context, item entity.Item) error {
	if _, err := r.GetByUserIDAndName(ctx, item.UserID, item.Name); err == nil {
		return repository.ErrItemAlreadyExist
	}

	itemFileName := getItemFileName(item.UserID, item.Name)
	itemFileData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("marshal item entity: %w", err)
	}
	itemReader := bytes.NewReader(itemFileData)
	params := s3.PutObjectInput{
		Bucket: &r.bucket,
		Key:    &itemFileName,
		Body:   itemReader,
	}
	_, err = r.client.PutObject(ctx, &params)
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	return nil
}

// Update store new version item in storage.
func (r *ItemRepository) Update(ctx context.Context, item entity.Item) error {
	if _, err := r.GetByUserIDAndName(ctx, item.UserID, item.Name); err != nil {
		return repository.ErrItemNotFound
	}

	itemFileName := getItemFileName(item.UserID, item.Name)
	itemFileData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("marshal item entity: %w", err)
	}
	itemReader := bytes.NewReader(itemFileData)
	params := s3.PutObjectInput{
		Bucket: &r.bucket,
		Key:    &itemFileName,
		Body:   itemReader,
	}
	_, err = r.client.PutObject(ctx, &params)
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	return nil
}

// GetByUserIDAndName return item from storage by user ID and item name.
func (r *ItemRepository) GetByUserIDAndName(ctx context.Context, userID string, name string) (entity.Item, error) {
	itemFileName := getItemFileName(userID, name)
	params := s3.GetObjectInput{
		Bucket: &r.bucket,
		Key:    &itemFileName,
	}
	out, err := r.client.GetObject(ctx, &params)
	if err != nil {
		return entity.Item{}, repository.ErrItemNotFound
	}
	itemFileData, err := io.ReadAll(out.Body)
	if err != nil {
		return entity.Item{}, fmt.Errorf("read item file body: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer out.Body.Close()

	var item entity.Item
	err = json.Unmarshal(itemFileData, &item)
	if err != nil {
		return entity.Item{}, fmt.Errorf("unmarshal item data: %w", err)
	}

	return item, nil
}

// Delete remove item from storage.
func (r *ItemRepository) Delete(ctx context.Context, item entity.Item) error {
	itemFileName := getItemFileName(item.UserID, item.Name)
	params := s3.DeleteObjectInput{
		Bucket: &r.bucket,
		Key:    &itemFileName,
	}
	_, err := r.client.DeleteObject(ctx, &params)
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}

// FindByUser return item names list from storage for user ID.
func (r *ItemRepository) FindByUser(ctx context.Context, userID string) ([]string, error) {
	userFolderName := getUserFolderName(userID)
	params := s3.ListObjectsV2Input{
		Bucket: &r.bucket,
		Prefix: &userFolderName,
	}
	out, err := r.client.ListObjectsV2(ctx, &params)
	if err != nil {
		return []string{}, fmt.Errorf("get items list: %w", err)
	}

	re, _ := regexp.Compile(`^_items/[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}/(.*)\.json$`)
	var list []string
	for _, obj := range out.Contents {
		sm := re.FindStringSubmatch(*obj.Key)
		if len(sm) == 2 {
			list = append(list, sm[1])
		}
	}

	return list, nil
}

func getItemFileName(userID string, name string) string {
	return fmt.Sprintf("_items/%s/%s.json", userID, name)
}

func getUserFolderName(userID string) string {
	return fmt.Sprintf("_items/%s", userID)
}
