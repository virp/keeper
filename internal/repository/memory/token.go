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

type TokenRepository struct {
	mu       *sync.RWMutex
	tokens   map[string]entity.Token
	lifetime time.Duration
}

func NewTokenRepository(lifetime time.Duration) *TokenRepository {
	return &TokenRepository{
		mu:       new(sync.RWMutex),
		tokens:   map[string]entity.Token{},
		lifetime: lifetime,
	}
}

func (r *TokenRepository) CreateToken(_ context.Context, user entity.User) (string, error) {
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
		UserLogin: user.Login,
		CreatedAt: time.Now(),
	}

	return token, nil
}

func (r *TokenRepository) GetToken(_ context.Context, id string) (entity.Token, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, ok := r.tokens[id]
	if !ok {
		return entity.Token{}, repository.ErrTokenNotFound
	}
	if time.Since(token.CreatedAt) > r.lifetime {
		r.mu.RUnlock()
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.tokens, id)
		return entity.Token{}, repository.ErrTokenExpired
	}

	return token, nil
}
