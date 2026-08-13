package domains

import (
	"time"
)

type User struct {
	ID             string `db:"id" json:"id"`
	Username       string `db:"username" json:"username"`
	Name           string `db:"name" json:"name"`
	Email          string `db:"email" json:"email"`
	HashedPassword string `db:"hashed_password" json:"-"`
	ProfilePicture string `db:"profile_picture" json:"profile_picture,omitempty"`
	IsActive       bool   `db:"is_active" json:"is_active"`
	EmailVerified  bool   `db:"email_verified" json:"email_verified"`
	CreatedAt      string `db:"created_at" json:"created_at"`
	UpdatedAt      string `db:"updated_at" json:"updated_at"`
}

func NewUser(id, username, name, email, hashedPassword string) *User {
	return &User{
		ID:             id,
		Username:       username,
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
		IsActive:       true,
		EmailVerified:  false,
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}
}
