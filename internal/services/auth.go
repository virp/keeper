package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"keeper/internal/entity"
)

const (
	loginMinLength    = 4
	passwordMinLength = 6
)

var (
	ErrUserInvalidPassword = errors.New("invalid user password")
)

// AuthService implement logic for working with users.
type AuthService struct {
	idGenerator     IdGenerator
	passwordHasher  PasswordHasher
	userRepository  UserRepository
	tokenRepository TokenRepository
}

// NewAuthService construct new AuthService.
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

// Auth check login and password and return new user token.
func (s *AuthService) Auth(ctx context.Context, login string, password string) (string, error) {
	user, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
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

// Register creates new user and return new user token.
func (s *AuthService) Register(ctx context.Context, login string, password string) (string, error) {
	var fields FieldErrors
	if len(login) < loginMinLength {
		field := FieldError{
			Field: "login",
			Error: fmt.Sprintf("length should be greater or equal %d", loginMinLength),
		}
		fields = append(fields, field)
	}
	if len(password) < passwordMinLength {
		field := FieldError{
			Field: "password",
			Error: fmt.Sprintf("length should be greater or equal %d", passwordMinLength),
		}
		fields = append(fields, field)
	}
	if len(fields) > 0 {
		return "", fields
	}

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
		return "", err
	}
	token, err := s.tokenRepository.CreateToken(ctx, user)
	if err != nil {
		return "", fmt.Errorf("new user token generate: %w", err)
	}
	return token, nil
}
