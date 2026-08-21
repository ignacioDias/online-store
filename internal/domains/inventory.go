package domains

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID               string    `json:"id" db:"id"`
	VariantID        string    `json:"variant_id" db:"variant_id"`
	Quantity         int       `json:"quantity" db:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity" db:"reserved_quantity"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type InventoryMovement struct {
	ID        string    `json:"id" db:"id"`
	VariantID string    `json:"variant_id" db:"variant_id"`
	Type      string    `json:"type" db:"type"`
	Quantity  int       `json:"quantity" db:"quantity"`
	OrderID   *string   `json:"order_id,omitempty" db:"order_id"`
	Reason    *string   `json:"reason,omitempty" db:"reason"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func NewInventory(variantID string) *Inventory {
	return &Inventory{ID: uuid.NewString(), VariantID: variantID}
}

func NewInventoryMovement(variantID, movementType string, quantity int, orderID, reason *string) *InventoryMovement {
	return &InventoryMovement{ID: uuid.NewString(), VariantID: variantID, Type: movementType, Quantity: quantity, OrderID: orderID, Reason: reason}
}
