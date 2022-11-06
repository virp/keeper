package services

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"keeper/internal/entity"
)

type IdGenerator interface {
	Generate() string
}

type PasswordHasher interface {
	Generate(password string) (string, error)
	Check(password string, hash string) bool
}

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	GetByLogin(ctx context.Context, login string) (entity.User, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
}

type TokenRepository interface {
	CreateToken(ctx context.Context, user entity.User) (string, error)
	GetToken(ctx context.Context, id string) (entity.Token, error)
}

type ItemRepository interface {
	Create(ctx context.Context, item entity.Item) error
	Update(ctx context.Context, item entity.Item) error
	GetByUserIDAndName(ctx context.Context, userID string, name string) (entity.Item, error)
	Delete(ctx context.Context, item entity.Item) error
	FindByUser(ctx context.Context, userID string) ([]entity.Item, error)
}

type UuidGenerator struct {
}

func (g *UuidGenerator) Generate() string {
	return uuid.NewString()
}

type BCryptPasswordHasher struct {
	Cost int
}

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

func (ph *BCryptPasswordHasher) Check(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
