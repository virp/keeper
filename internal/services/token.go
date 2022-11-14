package services

import (
	"context"

	"keeper/internal/entity"
)

type TokenService struct {
	userRepository  UserRepository
	tokenRepository TokenRepository
}

func NewTokenService(userRepository UserRepository, tokenRepository TokenRepository) *TokenService {
	return &TokenService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}

func (s *TokenService) GetUser(ctx context.Context, token string) (entity.User, error) {
	te, err := s.tokenRepository.GetToken(ctx, token)
	if err != nil {
		return entity.User{}, err
	}
	return s.userRepository.GetByLogin(ctx, te.UserLogin)
}
