package repositories

import (
	"context"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domains.Product) error
	GetProductByID(ctx context.Context, id string) (*domains.Product, error)
	ListProducts(ctx context.Context) ([]*domains.Product, error)
	UpdateProduct(ctx context.Context, product *domains.Product) error
	DeleteProduct(ctx context.Context, id string) error
	CreateVariant(ctx context.Context, variant *domains.ProductVariant) error
	GetVariantByID(ctx context.Context, id string) (*domains.ProductVariant, error)
	ListVariantsByProductID(ctx context.Context, productID string) ([]*domains.ProductVariant, error)
	UpdateVariant(ctx context.Context, variant *domains.ProductVariant) error
	DeleteVariant(ctx context.Context, id string) error
}

type productRepository struct{ db *sqlx.DB }

var ErrProductNotFound = errors.New("product not found")
var ErrVariantNotFound = errors.New("product variant not found")

const productColumns = `id, name, description, base_price, category_id, created_at`
const variantColumns = `id, product_id, sku, price, attributes, image_url, created_at`

func NewProductRepository(db *sqlx.DB) ProductRepository { return &productRepository{db: db} }

func (r *productRepository) CreateProduct(ctx context.Context, p *domains.Product) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO products (id, name, description, base_price, category_id) VALUES ($1,$2,$3,$4,$5) RETURNING created_at`, p.ID, p.Name, p.Description, p.BasePrice, p.CategoryID).Scan(&p.CreatedAt)
}

func (r *productRepository) GetProductByID(ctx context.Context, id string) (*domains.Product, error) {
	var p domains.Product
	if err := r.db.GetContext(ctx, &p, `SELECT `+productColumns+` FROM products WHERE id = $1`, id); err != nil {
		return nil, ErrProductNotFound
	}
	variants, err := r.ListVariantsByProductID(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Variants = variants
	return &p, nil
}

func (r *productRepository) ListProducts(ctx context.Context) ([]*domains.Product, error) {
	var products []*domains.Product
	if err := r.db.SelectContext(ctx, &products, `SELECT `+productColumns+` FROM products ORDER BY created_at DESC`); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, p *domains.Product) error {
	result, err := r.db.ExecContext(ctx, `UPDATE products SET name=$1, description=$2, base_price=$3, category_id=$4 WHERE id=$5`, p.Name, p.Description, p.BasePrice, p.CategoryID, p.ID)
	return CheckErrResult(result, err, ErrProductNotFound)
}

func (r *productRepository) DeleteProduct(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	return CheckErrResult(result, err, ErrProductNotFound)
}

func (r *productRepository) CreateVariant(ctx context.Context, v *domains.ProductVariant) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO product_variants (id, product_id, sku, price, attributes, image_url) VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at`, v.ID, v.ProductID, v.SKU, v.Price, v.Attributes, v.ImageURL).Scan(&v.CreatedAt)
}

func (r *productRepository) GetVariantByID(ctx context.Context, id string) (*domains.ProductVariant, error) {
	var v domains.ProductVariant
	if err := r.db.GetContext(ctx, &v, `SELECT `+variantColumns+` FROM product_variants WHERE id=$1`, id); err != nil {
		return nil, ErrVariantNotFound
	}
	return &v, nil
}

func (r *productRepository) ListVariantsByProductID(ctx context.Context, productID string) ([]*domains.ProductVariant, error) {
	var variants []*domains.ProductVariant
	if err := r.db.SelectContext(ctx, &variants, `SELECT `+variantColumns+` FROM product_variants WHERE product_id=$1 ORDER BY created_at`, productID); err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *productRepository) UpdateVariant(ctx context.Context, v *domains.ProductVariant) error {
	result, err := r.db.ExecContext(ctx, `UPDATE product_variants SET sku=$1, price=$2, attributes=$3, image_url=$4 WHERE id=$5`, v.SKU, v.Price, v.Attributes, v.ImageURL, v.ID)
	return CheckErrResult(result, err, ErrVariantNotFound)
}

func (r *productRepository) DeleteVariant(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM product_variants WHERE id=$1`, id)
	return CheckErrResult(result, err, ErrVariantNotFound)
}
