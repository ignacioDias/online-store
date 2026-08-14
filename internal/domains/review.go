package domains

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	Score     int       `db:"score" json:"score"`
	Comment   *string   `db:"comment" json:"comment"`
	ImageURL  *string   `db:"image_url" json:"image_url"`
	Likes     int       `db:"likes" json:"likes"`
	Dislikes  int       `db:"dislikes" json:"dislikes"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewReview(userID, productID string, score int, comment, imageURL *string) *Review {
	return &Review{
		ID:        uuid.NewString(),
		UserID:    userID,
		ProductID: productID,
		Score:     score,
		Comment:   comment,
		ImageURL:  imageURL,
		Likes:     0,
		Dislikes:  0,
	}
}
