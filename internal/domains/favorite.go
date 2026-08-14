package domains

import "time"

type Favorite struct {
	UserID    string    `json:"userId" db:"user_id"`
	ProductID string    `json:"productId" db:"product_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

func NewFavorite(userID, productID string) *Favorite {
	return &Favorite{
		UserID:    userID,
		ProductID: productID,
	}
}
