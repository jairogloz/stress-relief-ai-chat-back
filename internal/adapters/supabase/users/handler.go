package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stress-relief-ai-chat-back/internal/domain"
	"stress-relief-ai-chat-back/internal/ports"
)

type handler struct {
	apiKey     string
	logger     ports.Logger
	projectURL string
}

func NewUserAPIHandler(apiKey, projectURL string, logger ports.Logger) (ports.UserDataAPIHandler, error) {
	s := &handler{
		apiKey:     apiKey,
		logger:     logger,
		projectURL: projectURL,
	}
	if s.apiKey == "" {
		return nil, fmt.Errorf("apiKey can't be empty")
	}
	if s.projectURL == "" {
		return nil, fmt.Errorf("projectURL can't be empty")
	}
	if s.logger == nil {
		return nil, fmt.Errorf("logger can't be nil")
	}
	return s, nil
}

func (s handler) GetByID(ctx context.Context, userID string) (*domain.UserData, error) {
	if userID == "" {
		s.logger.Debug(ctx, "Can't get user with empty userID")
		return nil, fmt.Errorf("can't get user with empty userID")
	}
	url := fmt.Sprintf("%s/rest/v1/user_data?user_id=eq.%s", s.projectURL, userID)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		s.logger.Error(ctx, "Error creating request", "error", err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("apikey", s.apiKey)

	res, err := client.Do(req)
	if err != nil {
		s.logger.Error(ctx, "Error getting user", "error", err)
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			s.logger.Error(ctx, "Error closing response body", "error", err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Error(ctx, "Error reading response body", "error", err)
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Objects are returned as an array, so we unmarshal it as an array
	var userArray []domain.UserData
	err = json.Unmarshal(body, &userArray)
	if err != nil {
		s.logger.Error(ctx, "Error unmarshalling response body", "error", err)
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	if len(userArray) == 0 {
		return nil, fmt.Errorf("%w: user_data not found", domain.ErrNotFound)
	}

	return &userArray[0], nil
}
