package domains

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	ParentID  *string   `json:"parent_id,omitempty" db:"parent_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func NewCategory(name, slug string, parentID *string) *Category {
	return &Category{
		ID:       uuid.NewString(),
		Name:     name,
		Slug:     slug,
		ParentID: parentID,
	}
}
