package domains

import "time"

type ShoppingCartItem struct {
	UserID    string    `db:"user_id" json:"user_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	AddedAt   time.Time `db:"added_at" json:"added_at"`
}

func NewShoppingCartItem(userID, productID string, quantity int) *ShoppingCartItem {
	return &ShoppingCartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
}
