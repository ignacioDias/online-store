package repositories

import (
	"context"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domains.Order) error
	GetOrderByID(ctx context.Context, id string) (*domains.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]*domains.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status domains.OrderStatus) error
}

type orderRepository struct{ db *sqlx.DB }

var ErrOrderNotFound = errors.New("order not found")

const orderColumns = `id, user_id, status, currency, subtotal, discount_total, shipping_total, tax_total, total, coupon_id, shipping_address, created_at, updated_at`
const orderItemColumns = `id, order_id, product_id, variant_id, sku, product_name, quantity, unit_price, total_price, created_at`

func NewOrderRepository(db *sqlx.DB) OrderRepository { return &orderRepository{db: db} }

func (r *orderRepository) CreateOrder(ctx context.Context, order *domains.Order) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := `INSERT INTO orders (id, user_id, status, currency, subtotal, discount_total, shipping_total, tax_total, total, coupon_id, shipping_address) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING created_at, updated_at`
	if err := tx.QueryRowContext(ctx, query, order.ID, order.UserID, order.Status, order.Currency, order.Subtotal, order.DiscountTotal, order.ShippingTotal, order.TaxTotal, order.Total, order.CouponID, order.ShippingAddress).Scan(&order.CreatedAt, &order.UpdatedAt); err != nil {
		return err
	}

	itemQuery := `INSERT INTO order_items (id, order_id, product_id, variant_id, sku, product_name, quantity, unit_price, total_price) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING created_at`
	for _, item := range order.Items {
		item.OrderID = order.ID
		if err := tx.QueryRowContext(ctx, itemQuery, item.ID, item.OrderID, item.ProductID, item.VariantID, item.SKU, item.ProductName, item.Quantity, item.UnitPrice, item.TotalPrice).Scan(&item.CreatedAt); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *orderRepository) GetOrderByID(ctx context.Context, id string) (*domains.Order, error) {
	var order domains.Order
	if err := r.db.GetContext(ctx, &order, `SELECT `+orderColumns+` FROM orders WHERE id=$1`, id); err != nil {
		return nil, ErrOrderNotFound
	}
	if err := r.db.SelectContext(ctx, &order.Items, `SELECT `+orderItemColumns+` FROM order_items WHERE order_id=$1 ORDER BY created_at`, id); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]*domains.Order, error) {
	var orders []*domains.Order
	if err := r.db.SelectContext(ctx, &orders, `SELECT `+orderColumns+` FROM orders WHERE user_id=$1 ORDER BY created_at DESC`, userID); err != nil {
		return nil, err
	}
	for _, order := range orders {
		if err := r.db.SelectContext(ctx, &order.Items, `SELECT `+orderItemColumns+` FROM order_items WHERE order_id=$1 ORDER BY created_at`, order.ID); err != nil {
			return nil, err
		}
	}
	return orders, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, id string, status domains.OrderStatus) error {
	result, err := r.db.ExecContext(ctx, `UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	return CheckErrResult(result, err, ErrOrderNotFound)
}
