package ports

import (
	"context"
	"stress-relief-ai-chat-back/internal/domain"
)

// ChatService exposes the services provided by this application around chat.
type ChatService interface {
	ProcessMessage(ctx context.Context, message *domain.ChatMessage, threadID *string) (*domain.ChatResponse, error)
}

// ChatHandler is an interface for handling chat messages against an AI service.
type ChatHandler interface {
	ProcessMessage(ctx context.Context, message *domain.ChatMessage, threadID *string) (*domain.ChatResponse, error)
}
