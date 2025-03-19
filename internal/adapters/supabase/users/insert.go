package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stress-relief-ai-chat-back/internal/domain"
)

func (s handler) Insert(ctx context.Context, userID string, userData *domain.UserData) error {
	if userID == "" {
		s.logger.Debug(ctx, "Can't insert user with empty userID")
		return fmt.Errorf("can't insert user with empty userID")
	}
	if userData == nil {
		s.logger.Debug(ctx, "Can't insert nil userData")
		return fmt.Errorf("can't insert nil userData")
	}

	url := fmt.Sprintf("%s/rest/v1/user_data", s.projectURL)

	userData.UserID = userID
	data, err := json.Marshal(userData)
	if err != nil {
		s.logger.Error(ctx, "Error marshalling userData", "error", err)
		return fmt.Errorf("error marshalling userData: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		s.logger.Error(ctx, "Error creating request", "error", err)
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("apikey", s.apiKey)

	res, err := client.Do(req)
	if err != nil {
		s.logger.Error(ctx, "Error inserting user", "error", err)
		return fmt.Errorf("error inserting user: %w", err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			s.logger.Error(ctx, "Error closing response body", "error", err)
		}
	}()

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		s.logger.Error(ctx, "Error response from server", "status", res.StatusCode, "body", string(body))
		return fmt.Errorf("error response from server: %s", res.Status)
	}

	return nil
}
