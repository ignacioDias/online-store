package repositories

import (
	"context"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type ShoppingCartRepository interface {
	AddItemsToCart(ctx context.Context, item *domains.ShoppingCartItem) error
	RemoveItemFromCart(ctx context.Context, userID, productID string) error
	UpdateItemQuantity(ctx context.Context, item *domains.ShoppingCartItem) error
	GetCartItemsFromUser(ctx context.Context, userID string) ([]*domains.ShoppingCartItem, error)
	ClearCart(ctx context.Context, userID string) error
}

type shoppingCartRepository struct {
	db *sqlx.DB
}

var ErrShoppingCartItemNotFound = errors.New("Shopping cart item not found")

func NewShoppingCartRepository(db *sqlx.DB) ShoppingCartRepository {
	return &shoppingCartRepository{db: db}
}

func (scr *shoppingCartRepository) AddItemsToCart(ctx context.Context, item *domains.ShoppingCartItem) error {
	query := `INSERT INTO shopping_cart (user_id, product_id, quantity, added_at) VALUES ($1, $2, $3, NOW())`
	_, err := scr.db.ExecContext(ctx, query, item.UserID, item.ProductID, item.Quantity)
	return err
}

func (scr *shoppingCartRepository) RemoveItemFromCart(ctx context.Context, userID, productID string) error {
	query := `DELETE FROM shopping_cart WHERE user_id = $1 AND product_id = $2`
	result, err := scr.db.ExecContext(ctx, query, userID, productID)
	return CheckErrResult(result, err, ErrShoppingCartItemNotFound)
}

func (scr *shoppingCartRepository) UpdateItemQuantity(ctx context.Context, item *domains.ShoppingCartItem) error {
	query := `UPDATE shopping_cart SET quantity = $1, added_at = NOW() WHERE user_id = $2 AND product_id = $3`
	result, err := scr.db.ExecContext(ctx, query, item.Quantity, item.UserID, item.ProductID)
	return CheckErrResult(result, err, ErrShoppingCartItemNotFound)
}

func (scr *shoppingCartRepository) GetCartItemsFromUser(ctx context.Context, userID string) ([]*domains.ShoppingCartItem, error) {
	query := `SELECT user_id, product_id, quantity, added_at FROM shopping_cart WHERE user_id = $1`
	var items []*domains.ShoppingCartItem
	err := scr.db.SelectContext(ctx, &items, query, userID)
	if err != nil {
		return []*domains.ShoppingCartItem{}, err
	}
	if len(items) == 0 {
		return items, ErrShoppingCartItemNotFound
	}
	return items, nil
}

func (scr *shoppingCartRepository) ClearCart(ctx context.Context, userID string) error {
	query := `DELETE FROM shopping_cart WHERE user_id = $1`
	result, err := scr.db.ExecContext(ctx, query, userID)
	return CheckErrResult(result, err, ErrShoppingCartItemNotFound)
}
