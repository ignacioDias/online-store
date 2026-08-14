package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type FavoriteRepository interface {
	AddFavorite(ctx context.Context, favorite *domains.Favorite) error
	RemoveFavorite(ctx context.Context, favorite *domains.Favorite) error
	GetFavoritesByUserID(ctx context.Context, userID string) ([]string, error)
}

type favoriteRepository struct {
	db *sqlx.DB
}

var ErrFavoriteNotFound = errors.New("Favorite not found")
var ErrFavoriteAlreadyExists = errors.New("Favorite already exists")

func NewFavoriteRepository(db *sqlx.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (fr *favoriteRepository) AddFavorite(ctx context.Context, favorite *domains.Favorite) error {
	query := `INSERT INTO favorites (user_id, product_id) VALUES ($1, $2) ON CONFLICT (user_id, product_id) DO NOTHING RETURNING created_at`
	err := fr.db.GetContext(ctx, &favorite.CreatedAt, query, favorite.UserID, favorite.ProductID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrFavoriteAlreadyExists
	}
	return err
}

func (fr *favoriteRepository) RemoveFavorite(ctx context.Context, favorite *domains.Favorite) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND product_id = $2`
	result, err := fr.db.ExecContext(ctx, query, favorite.UserID, favorite.ProductID)
	return CheckErrResult(result, err, ErrFavoriteNotFound)
}

func (fr *favoriteRepository) GetFavoritesByUserID(ctx context.Context, userID string) ([]string, error) {
	query := `SELECT product_id FROM favorites WHERE user_id = $1`
	var productIDS []string
	err := fr.db.SelectContext(ctx, &productIDS, query, userID)
	if err != nil {
		return nil, fmt.Errorf("select favorites: %w", err)
	}
	if len(productIDS) == 0 {
		return []string{}, nil
	}
	return productIDS, nil
}
