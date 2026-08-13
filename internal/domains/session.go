package domains

import (
	"time"
)

type Session struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
}

func NewSession(id, userID string) *Session {
	return &Session{
		ID:        id,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(30 * time.Minute),
	}
}
