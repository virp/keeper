package services

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"keeper/internal/entity"
)

// IdGenerator interface describe required logic for generating ID for business entities.
type IdGenerator interface {
	Generate() string
}

// PasswordHasher interface describe required logic for hashing and checking user passwords.
type PasswordHasher interface {
	Generate(password string) (string, error)
	Check(password string, hash string) bool
}

// UserRepository interface describe required logic for storing users.
type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	GetByLogin(ctx context.Context, login string) (entity.User, error)
}

// TokenRepository interface describe required logic for storing user tokens.
type TokenRepository interface {
	CreateToken(ctx context.Context, user entity.User) (string, error)
	GetToken(ctx context.Context, id string) (entity.Token, error)
}

// ItemRepository interface describe required logic for storing items.
type ItemRepository interface {
	Create(ctx context.Context, item entity.Item) error
	Update(ctx context.Context, item entity.Item) error
	GetByUserIDAndName(ctx context.Context, userID string, name string) (entity.Item, error)
	Delete(ctx context.Context, item entity.Item) error
	FindByUser(ctx context.Context, userID string) ([]string, error)
}

// UuidGenerator contains implementation for UUID generation.
type UuidGenerator struct {
}

func (g *UuidGenerator) Generate() string {
	return uuid.NewString()
}

// BCryptPasswordHasher contains implementation for creating password hash and checks.
type BCryptPasswordHasher struct {
	Cost int
}

// Generate creates password hash.
func (ph *BCryptPasswordHasher) Generate(password string) (string, error) {
	cost := ph.Cost
	if cost == 0 {
		cost = 13
	}
	b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Check compare password and password hash.
func (ph *BCryptPasswordHasher) Check(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
