package repositories

import (
	"context"
	"errors"
	"fmt"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type InventoryRepository interface {
	GetByVariantID(ctx context.Context, variantID string) (*domains.Inventory, error)
	SetQuantity(ctx context.Context, inventory *domains.Inventory) error
	Reserve(ctx context.Context, variantID string, quantity int) error
	Release(ctx context.Context, variantID string, quantity int) error
	AddMovement(ctx context.Context, movement *domains.InventoryMovement) error
	ListMovementsByVariantID(ctx context.Context, variantID string) ([]*domains.InventoryMovement, error)
}

type inventoryRepository struct{ db *sqlx.DB }

var ErrInventoryNotFound = errors.New("inventory not found")
var ErrInsufficientInventory = errors.New("insufficient inventory")

func NewInventoryRepository(db *sqlx.DB) InventoryRepository { return &inventoryRepository{db: db} }

func (r *inventoryRepository) GetByVariantID(ctx context.Context, variantID string) (*domains.Inventory, error) {
	var inventory domains.Inventory
	err := r.db.GetContext(ctx, &inventory, `SELECT id, variant_id, quantity, reserved_quantity, updated_at FROM inventory WHERE variant_id=$1`, variantID)
	if err != nil {
		return nil, ErrInventoryNotFound
	}
	return &inventory, nil
}

func (r *inventoryRepository) SetQuantity(ctx context.Context, inventory *domains.Inventory) error {
	result, err := r.db.ExecContext(ctx, `UPDATE inventory SET quantity=$1, reserved_quantity=$2, updated_at=NOW() WHERE variant_id=$3`, inventory.Quantity, inventory.ReservedQuantity, inventory.VariantID)
	return CheckErrResult(result, err, ErrInventoryNotFound)
}

func (r *inventoryRepository) Reserve(ctx context.Context, variantID string, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("reserve quantity must be positive")
	}
	result, err := r.db.ExecContext(ctx, `UPDATE inventory SET reserved_quantity=reserved_quantity+$1, updated_at=NOW() WHERE variant_id=$2 AND quantity-reserved_quantity >= $1`, quantity, variantID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		if _, err := r.GetByVariantID(ctx, variantID); err != nil {
			return err
		}
		return ErrInsufficientInventory
	}
	return nil
}

func (r *inventoryRepository) Release(ctx context.Context, variantID string, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("release quantity must be positive")
	}
	result, err := r.db.ExecContext(ctx, `UPDATE inventory SET reserved_quantity=reserved_quantity-$1, updated_at=NOW() WHERE variant_id=$2 AND reserved_quantity >= $1`, quantity, variantID)
	if err != nil {
		return err
	}
	return CheckErrResult(result, nil, ErrInventoryNotFound)
}

func (r *inventoryRepository) AddMovement(ctx context.Context, movement *domains.InventoryMovement) error {
	query := `INSERT INTO inventory_movements (id, variant_id, type, quantity, order_id, reason) VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at`
	return r.db.QueryRowContext(ctx, query, movement.ID, movement.VariantID, movement.Type, movement.Quantity, movement.OrderID, movement.Reason).Scan(&movement.CreatedAt)
}

func (r *inventoryRepository) ListMovementsByVariantID(ctx context.Context, variantID string) ([]*domains.InventoryMovement, error) {
	var movements []*domains.InventoryMovement
	err := r.db.SelectContext(ctx, &movements, `SELECT id, variant_id, type, quantity, order_id, reason, created_at FROM inventory_movements WHERE variant_id=$1 ORDER BY created_at DESC`, variantID)
	if err != nil {
		return nil, err
	}
	return movements, nil
}
