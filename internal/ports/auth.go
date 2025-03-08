package ports

import "stress-relief-ai-chat-back/internal/domain"

type AuthPort interface {
	Authenticate(token string) (*domain.User, error)
}
