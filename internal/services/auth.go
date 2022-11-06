package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"keeper/internal/entity"
	"keeper/internal/repository"
)

var (
	ErrUserInvalidPassword = errors.New("invalid user password")
)

type AuthService struct {
	idGenerator     IdGenerator
	passwordHasher  PasswordHasher
	userRepository  UserRepository
	tokenRepository TokenRepository
}

func NewAuthService(
	idGenerator IdGenerator,
	passwordHasher PasswordHasher,
	userRepository UserRepository,
	tokenRepository TokenRepository,
) *AuthService {
	return &AuthService{
		idGenerator:     idGenerator,
		passwordHasher:  passwordHasher,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}

func (s *AuthService) Auth(ctx context.Context, login string, password string) (string, error) {
	user, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			passwordHash, err := s.passwordHasher.Generate(password)
			if err != nil {
				return "", fmt.Errorf("password hash generate: %w", err)
			}
			user := entity.User{
				ID:           s.idGenerator.Generate(),
				Login:        login,
				PasswordHash: passwordHash,
				CreatedAt:    time.Now(),
			}
			if err := s.userRepository.Create(ctx, user); err != nil {
				return "", fmt.Errorf("user creare: %w", err)
			}
			token, err := s.tokenRepository.CreateToken(ctx, user)
			if err != nil {
				return "", fmt.Errorf("new user token generate: %w", err)
			}
			return token, nil
		}
		return "", err
	}

	if !s.passwordHasher.Check(password, user.PasswordHash) {
		return "", ErrUserInvalidPassword
	}

	token, err := s.tokenRepository.CreateToken(ctx, user)
	if err != nil {
		return "", fmt.Errorf("user token generate: %w", err)
	}
	return token, nil
}
