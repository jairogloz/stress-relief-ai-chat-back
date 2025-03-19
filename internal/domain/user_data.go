package domain

type UserData struct {
	UserID   string  `json:"user_id"`
	ThreadID *string `json:"thread_id"`
}
