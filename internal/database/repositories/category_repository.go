package repositories

import (
	"context"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domains.Category) error
	GetCategoryByID(ctx context.Context, id string) (*domains.Category, error)
	GetAllCategories(ctx context.Context) ([]*domains.Category, error)
	GetCategoriesByParentID(ctx context.Context, parentID string) ([]*domains.Category, error)
	UpdateCategory(ctx context.Context, category *domains.Category) error
	DeleteCategoryByID(ctx context.Context, id string) error
}

type categoryRepository struct {
	db *sqlx.DB
}

var ErrCategoryNotFound = errors.New("Category not found")

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (c *categoryRepository) CreateCategory(ctx context.Context, category *domains.Category) error {
	query := `INSERT INTO categories (id, name, slug, parent_id) VALUES ($1, $2, $3, $4) RETURNING created_at`
	err := c.db.QueryRowContext(ctx, query, category.ID, category.Name, category.Slug, category.ParentID).Scan(&category.CreatedAt)
	return err
}

func (c *categoryRepository) GetCategoryByID(ctx context.Context, id string) (*domains.Category, error) {
	query := `SELECT id, name, slug, parent_id, created_at FROM categories WHERE id = $1`
	var category domains.Category
	err := c.db.GetContext(ctx, &category, query, id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *categoryRepository) GetAllCategories(ctx context.Context) ([]*domains.Category, error) {
	query := `SELECT id, name, slug, parent_id, created_at FROM categories`
	var categories []*domains.Category
	err := c.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *categoryRepository) GetCategoriesByParentID(ctx context.Context, parentID string) ([]*domains.Category, error) {
	query := `SELECT id, name, slug, parent_id, created_at FROM categories WHERE parent_id = $1`
	var categories []*domains.Category
	err := c.db.SelectContext(ctx, &categories, query, parentID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *categoryRepository) UpdateCategory(ctx context.Context, category *domains.Category) error {
	query := `UPDATE categories SET name = $1, slug = $2, parent_id = $3 WHERE id = $4`
	result, err := c.db.ExecContext(ctx, query, category.Name, category.Slug, category.ParentID, category.ID)
	return CheckErrResult(result, err, ErrCategoryNotFound)
}

func (c *categoryRepository) DeleteCategoryByID(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := c.db.ExecContext(ctx, query, id)
	return CheckErrResult(result, err, ErrCategoryNotFound)
}
