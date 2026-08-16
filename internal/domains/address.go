package domains

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Street     string    `json:"street" db:"street"`
	Apartment  *string   `json:"apartment,omitempty" db:"apartment"`
	City       string    `json:"city" db:"city"`
	State      string    `json:"state" db:"state"`
	PostalCode string    `json:"postal_code" db:"postal_code"`
	Country    string    `json:"country" db:"country"`
	IsDefault  bool      `json:"is_default" db:"is_default"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func NewAddress(userID, street, apartment, city, state, postalCode, country string, isDefault bool) *Address {
	return &Address{
		ID:         uuid.NewString(),
		UserID:     userID,
		Street:     street,
		Apartment:  &apartment,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
		IsDefault:  isDefault,
	}
}
