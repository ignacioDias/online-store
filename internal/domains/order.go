package domains

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderConfirmed  OrderStatus = "confirmed"
	OrderProcessing OrderStatus = "processing"
	OrderShipped    OrderStatus = "shipped"
	OrderDelivered  OrderStatus = "delivered"
	OrderCancelled  OrderStatus = "cancelled"
	OrderRefunded   OrderStatus = "refunded"
)

type Order struct {
	ID              string          `json:"id" db:"id"`
	UserID          string          `json:"user_id" db:"user_id"`
	Status          OrderStatus     `json:"status" db:"status"`
	Currency        string          `json:"currency" db:"currency"`
	Subtotal        float64         `json:"subtotal" db:"subtotal"`
	DiscountTotal   float64         `json:"discount_total" db:"discount_total"`
	ShippingTotal   float64         `json:"shipping_total" db:"shipping_total"`
	TaxTotal        float64         `json:"tax_total" db:"tax_total"`
	Total           float64         `json:"total" db:"total"`
	CouponID        *string         `json:"coupon_id,omitempty" db:"coupon_id"`
	ShippingAddress json.RawMessage `json:"shipping_address" db:"shipping_address"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
	Items           []*OrderItem    `json:"items" db:"-"`
}

type OrderItem struct {
	ID          string    `json:"id" db:"id"`
	OrderID     string    `json:"order_id" db:"order_id"`
	ProductID   string    `json:"product_id" db:"product_id"`
	VariantID   string    `json:"variant_id" db:"variant_id"`
	SKU         string    `json:"sku" db:"sku"`
	ProductName string    `json:"product_name" db:"product_name"`
	Quantity    int       `json:"quantity" db:"quantity"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	TotalPrice  float64   `json:"total_price" db:"total_price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func NewOrder(userID, currency string, shippingAddress json.RawMessage) *Order {
	return &Order{ID: uuid.NewString(), UserID: userID, Status: OrderPending, Currency: currency, ShippingAddress: shippingAddress}
}
