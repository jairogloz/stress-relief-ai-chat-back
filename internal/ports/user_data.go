package ports

import (
	"context"
	"stress-relief-ai-chat-back/internal/domain"
)

type UserDataAPIHandler interface {
	GetByID(ctx context.Context, userID string) (*domain.UserData, error)
	Insert(ctx context.Context, userID string, userData *domain.UserData) error
}
