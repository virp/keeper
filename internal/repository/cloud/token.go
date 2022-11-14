package cloud

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

type TokenRepository struct {
	bucket   string
	client   *s3.Client
	lifetime time.Duration
}

func NewTokenRepository(client *s3.Client, bucket string, lifetime time.Duration) *TokenRepository {
	return &TokenRepository{
		bucket:   bucket,
		client:   client,
		lifetime: lifetime,
	}
}

func (r *TokenRepository) CreateToken(ctx context.Context, user entity.User) (string, error) {
	var sb strings.Builder
	sb.WriteString(user.ID)
	sb.WriteString(user.Login)
	sb.WriteString(strconv.Itoa(int(time.Now().UnixNano())))
	hash := md5.Sum([]byte(sb.String()))
	token := hex.EncodeToString(hash[:])

	tokenEntity := entity.Token{
		ID:        token,
		UserLogin: user.Login,
		CreatedAt: time.Now(),
	}

	tokenFileName := getTokenFileName(token)
	tokenData, err := json.Marshal(tokenEntity)
	if err != nil {
		return "", fmt.Errorf("marshal token entity: %w", err)
	}
	tokenReader := bytes.NewReader(tokenData)
	params := s3.PutObjectInput{
		Bucket: &r.bucket,
		Key:    &tokenFileName,
		Body:   tokenReader,
	}
	_, err = r.client.PutObject(ctx, &params)
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}

	return token, nil
}

func (r *TokenRepository) GetToken(ctx context.Context, id string) (entity.Token, error) {
	tokenFileName := getTokenFileName(id)
	params := s3.GetObjectInput{
		Bucket: &r.bucket,
		Key:    &tokenFileName,
	}
	out, err := r.client.GetObject(ctx, &params)
	if err != nil {
		return entity.Token{}, repository.ErrTokenNotFound
	}
	tokenFileData, err := io.ReadAll(out.Body)
	if err != nil {
		return entity.Token{}, fmt.Errorf("read token file body: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer out.Body.Close()

	var token entity.Token
	err = json.Unmarshal(tokenFileData, &token)
	if err != nil {
		return entity.Token{}, fmt.Errorf("unmarshal token data: %w", err)
	}

	if time.Since(token.CreatedAt) > r.lifetime {
		if err := r.deleteToken(ctx, token.ID); err != nil {
			return entity.Token{}, err
		}
		return entity.Token{}, repository.ErrTokenExpired
	}

	return token, nil
}

func (r *TokenRepository) deleteToken(ctx context.Context, id string) error {
	tokenFileName := getTokenFileName(id)
	params := s3.DeleteObjectInput{
		Bucket: &r.bucket,
		Key:    &tokenFileName,
	}
	_, err := r.client.DeleteObject(ctx, &params)
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}

func getTokenFileName(id string) string {
	return fmt.Sprintf("_tokens/%s.json", id)
}
