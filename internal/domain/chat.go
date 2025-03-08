package domain

import "errors"

// ChatMessage represents a message sent by the end user.
type ChatMessage struct {
	Content string `json:"content"`
}

func (chM *ChatMessage) Validate() error {
	if chM == nil {
		return errors.New("message cannot be nil")
	}
	if chM.Content == "" {
		return errors.New("message content cannot be empty")
	}
	return nil
}

// ChatResponse represents a message sent by the chatbot.
type ChatResponse struct {
	Content  string `json:"content"`
	ThreadID string `json:"threadId"`
}
