package chat

import (
	"context"
	"errors"
	"stress-relief-ai-chat-back/internal/domain"
	"stress-relief-ai-chat-back/internal/ports"
)

type service struct {
	chatAdapter ports.ChatHandler
	logger      ports.Logger
}

func NewChatService(chatAdapter ports.ChatHandler, l ports.Logger) ports.ChatService {
	ch := &service{
		chatAdapter: chatAdapter,
		logger:      l,
	}

	if ch.chatAdapter == nil {
		panic("Cannot create service without a ChatHandler")
	}
	if ch.logger == nil {
		panic("Cannot create service without a Logger")
	}

	return ch
}

func (s *service) ProcessMessage(ctx context.Context, message *domain.ChatMessage, threadID *string) (*domain.ChatResponse, error) {
	if message == nil {
		return nil, errors.New("message cannot be nil")
	}
	return s.chatAdapter.ProcessMessage(ctx, message, threadID)
}
