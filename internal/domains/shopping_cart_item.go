package domains

import "time"

type ShoppingCartItem struct {
	UserID    string    `db:"user_id" json:"user_id"`
	VariantID string    `db:"variant_id" json:"variant_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	AddedAt   time.Time `db:"added_at" json:"added_at"`
}

func NewShoppingCartItem(userID, variantID string, quantity int) *ShoppingCartItem {
	return &ShoppingCartItem{
		UserID:    userID,
		VariantID: variantID,
		Quantity:  quantity,
	}
}
