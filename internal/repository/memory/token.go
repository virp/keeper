package memory

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
	"time"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

const (
	tokenLifetime = 24 * time.Hour
)

type TokenMemoryRepository struct {
	mu     *sync.RWMutex
	tokens map[string]entity.Token
}

func NewTokenMemoryRepository() *TokenMemoryRepository {
	return &TokenMemoryRepository{
		mu:     new(sync.RWMutex),
		tokens: map[string]entity.Token{},
	}
}

func (r *TokenMemoryRepository) CreateToken(_ context.Context, user entity.User) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var sb strings.Builder
	sb.WriteString(user.ID)
	sb.WriteString(user.Login)
	sb.WriteString(strconv.Itoa(int(time.Now().UnixNano())))
	hash := md5.Sum([]byte(sb.String()))
	token := hex.EncodeToString(hash[:])

	r.tokens[token] = entity.Token{
		ID:        token,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}

	return token, nil
}

func (r *TokenMemoryRepository) GetToken(_ context.Context, id string) (entity.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, ok := r.tokens[id]
	if !ok {
		return entity.Token{}, repository.ErrTokenNotFound
	}
	if time.Since(token.CreatedAt) > tokenLifetime {
		r.mu.RUnlock()
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.tokens, id)
		return entity.Token{}, repository.ErrTokenExpired
	}

	return token, nil
}
