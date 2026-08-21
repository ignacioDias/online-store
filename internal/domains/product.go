package domains

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string            `json:"id" db:"id"`
	Name        string            `json:"name" db:"name"`
	Description *string           `json:"description,omitempty" db:"description"`
	BasePrice   float64           `json:"base_price" db:"base_price"`
	CategoryID  *string           `json:"category_id,omitempty" db:"category_id"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	Variants    []*ProductVariant `json:"variants,omitempty" db:"-"`
}

type ProductVariant struct {
	ID         string          `json:"id" db:"id"`
	ProductID  string          `json:"product_id" db:"product_id"`
	SKU        string          `json:"sku" db:"sku"`
	Price      *float64        `json:"price,omitempty" db:"price"`
	Attributes json.RawMessage `json:"attributes" db:"attributes"`
	ImageURL   *string         `json:"image_url,omitempty" db:"image_url"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}

func NewProduct(name string, basePrice float64, categoryID *string) *Product {
	return &Product{ID: uuid.NewString(), Name: name, BasePrice: basePrice, CategoryID: categoryID}
}

func NewProductVariant(productID, sku string, price *float64, attributes json.RawMessage) *ProductVariant {
	return &ProductVariant{ID: uuid.NewString(), ProductID: productID, SKU: sku, Price: price, Attributes: attributes}
}
