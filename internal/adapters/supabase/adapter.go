package supabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"stress-relief-ai-chat-back/internal/domain"
)

type AuthAdapter struct {
	supabaseURL string
	supabaseKey string
}

func NewAuthAdapter(url, key string) *AuthAdapter {
	return &AuthAdapter{
		supabaseURL: url,
		supabaseKey: key,
	}
}

func (a *AuthAdapter) Authenticate(token string) (*domain.User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/v1/user", a.supabaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", a.supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token")
	}

	var user domain.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
