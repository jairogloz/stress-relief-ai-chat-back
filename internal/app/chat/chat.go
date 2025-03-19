package chat

import (
	"context"
	"errors"
	"fmt"
	"stress-relief-ai-chat-back/internal/domain"
	"stress-relief-ai-chat-back/internal/ports"
)

type service struct {
	chatAdapter     ports.ChatHandler
	logger          ports.Logger
	userDataHandler ports.UserDataAPIHandler
}

func NewChatService(chatAdapter ports.ChatHandler, l ports.Logger, u ports.UserDataAPIHandler) ports.ChatService {
	ch := &service{
		chatAdapter:     chatAdapter,
		logger:          l,
		userDataHandler: u,
	}

	if ch.chatAdapter == nil {
		panic("Cannot create service without a ChatHandler")
	}
	if ch.logger == nil {
		panic("Cannot create service without a Logger")
	}
	if ch.userDataHandler == nil {
		panic("Cannot create service without a UserDataAPIHandler")
	}

	return ch
}

func (s *service) ProcessMessage(ctx context.Context, message *domain.ChatMessage, userID string) (*domain.ChatResponse, error) {
	if message == nil {
		return nil, errors.New("message cannot be nil")
	}
	// Get the user_data information from the database, to get the threadID if exists
	userData, err := s.userDataHandler.GetByID(ctx, userID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		s.logger.Warn(ctx, "could not get user_data information", "error", err.Error())
		return nil, fmt.Errorf("could not get user_data information: %w", err)
	}

	var threadId *string
	if userData != nil {
		threadId = userData.ThreadID
	}
	chatResponse, err := s.chatAdapter.ProcessMessage(ctx, message, threadId)
	if err != nil {
		s.logger.Error(ctx, "error processing message", "error", err.Error())
		return nil, fmt.Errorf("error processing message: %w", err)
	}

	// Update the user_data information with the new threadID, if needed
	if userData == nil {
		s.logger.Debug(ctx, "user_data information not found, creating new entry")
		err := s.userDataHandler.Insert(ctx, userID, &domain.UserData{
			UserID:   userID,
			ThreadID: &chatResponse.ThreadID,
		})
		if err != nil {
			s.logger.Warn(ctx, "could not update user_data information", "error", err.Error())
			return nil, fmt.Errorf("could not update user_data information: %w", err)
		}
	}

	return chatResponse, nil
}
