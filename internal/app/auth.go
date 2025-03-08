package app

import (
	"errors"
	"stress-relief-ai-chat-back/internal/domain"
	"stress-relief-ai-chat-back/internal/ports"
)

type AuthService struct {
	authAdapter ports.AuthPort
}

func NewAuthService(authAdapter ports.AuthPort) *AuthService {
	return &AuthService{
		authAdapter: authAdapter,
	}
}

func (s *AuthService) ValidateToken(token string) (*domain.User, error) {
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	return s.authAdapter.Authenticate(token)
}
