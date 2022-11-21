package services

import (
	"context"

	"keeper/internal/entity"
)

// TokenService implement logic for working with user tokens.
type TokenService struct {
	userRepository  UserRepository
	tokenRepository TokenRepository
}

// NewTokenService construct new TokenService.
func NewTokenService(userRepository UserRepository, tokenRepository TokenRepository) *TokenService {
	return &TokenService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}

// GetUser returns user by token.
func (s *TokenService) GetUser(ctx context.Context, token string) (entity.User, error) {
	te, err := s.tokenRepository.GetToken(ctx, token)
	if err != nil {
		return entity.User{}, err
	}
	return s.userRepository.GetByLogin(ctx, te.UserLogin)
}
