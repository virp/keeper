package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

type UserRepository struct {
	bucket string
	client *s3.Client
}

func NewUserRepository(client *s3.Client, bucket string) *UserRepository {
	return &UserRepository{
		bucket: bucket,
		client: client,
	}
}

//func (r *UserRepository) List(ctx context.Context) ([]string, error) {
//	bucketName := "keeper"
//	params := s3.ListObjectsV2Input{
//		Bucket: &bucketName,
//	}
//	p := s3.NewListObjectsV2Paginator(r.client, &params, func(options *s3.ListObjectsV2PaginatorOptions) {
//		options.Limit = 10
//	})
//
//	var items []string
//	for p.HasMorePages() {
//		page, err := p.NextPage(ctx)
//		if err != nil {
//			return []string{}, fmt.Errorf("get list page: %w", err)
//		}
//		for _, item := range page.Contents {
//			items = append(items, *item.Key)
//		}
//	}
//	return items, nil
//}

func (r *UserRepository) Create(ctx context.Context, user entity.User) error {
	if _, err := r.GetByLogin(ctx, user.Login); err == nil {
		return repository.ErrUserAlreadyExist
	}

	userFileName := getUserFileName(user.Login)
	userFileData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal user entity: %w", err)
	}
	userReader := bytes.NewReader(userFileData)
	params := s3.PutObjectInput{
		Bucket: &r.bucket,
		Key:    &userFileName,
		Body:   userReader,
	}
	_, err = r.client.PutObject(ctx, &params)
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	userFileName := getUserFileName(login)
	params := s3.GetObjectInput{
		Bucket: &r.bucket,
		Key:    &userFileName,
	}
	out, err := r.client.GetObject(ctx, &params)
	if err != nil {
		return entity.User{}, repository.ErrUserNotFound
	}
	userFileData, err := io.ReadAll(out.Body)
	if err != nil {
		return entity.User{}, fmt.Errorf("read user file body: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer out.Body.Close()

	var user entity.User
	err = json.Unmarshal(userFileData, &user)
	if err != nil {
		return entity.User{}, fmt.Errorf("unmarshal user data: %w", err)
	}

	return user, nil
}

func getUserFileName(login string) string {
	return fmt.Sprintf("_users/%s.json", login)
}
