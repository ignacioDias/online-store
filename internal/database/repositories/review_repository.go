package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *domains.Review) error
	GetReviewByID(ctx context.Context, id string) (*domains.Review, error)
	GetReviewsFromUser(ctx context.Context, userID string) ([]*domains.Review, error)
	GetReviewsFromProduct(ctx context.Context, productID string) ([]*domains.Review, error)
	UpdateReview(ctx context.Context, review *domains.Review) error
	DeleteReview(ctx context.Context, id string) error
}

type reviewRepository struct {
	db *sqlx.DB
}

var ErrReviewAlreadyExists = errors.New("Review already exists")
var ErrReviewNotFound = errors.New("Review not found")

func NewReviewRepository(db *sqlx.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (rr *reviewRepository) CreateReview(ctx context.Context, review *domains.Review) error {
	query := `INSERT INTO reviews (id, user_id, product_id, score, comment, image_url, likes, dislikes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (user_id, product_id) DO NOTHING RETURNING created_at`

	err := rr.db.QueryRowContext(ctx, query, review.ID, review.UserID, review.ProductID, review.Score, review.Comment, review.ImageURL, review.Likes, review.Dislikes).Scan(&review.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrReviewAlreadyExists
	}
	return err
}

func (rr *reviewRepository) GetReviewByID(ctx context.Context, id string) (*domains.Review, error) {
	query := `SELECT id, user_id, product_id, score, comment, image_url, likes, dislikes, created_at FROM reviews WHERE id = $1`
	var review domains.Review
	err := rr.db.GetContext(ctx, &review, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}
	return &review, nil
}

func (rr *reviewRepository) GetReviewsFromUser(ctx context.Context, userID string) ([]*domains.Review, error) {
	return rr.getReviewsByField(ctx, "user_id", userID)
}

func (rr *reviewRepository) GetReviewsFromProduct(ctx context.Context, productID string) ([]*domains.Review, error) {
	return rr.getReviewsByField(ctx, "product_id", productID)
}

func (rr *reviewRepository) getReviewsByField(ctx context.Context, field, value string) ([]*domains.Review, error) {
	query := `SELECT id, user_id, product_id, score, comment, image_url, likes, dislikes, created_at FROM reviews WHERE ` + field + ` = $1`
	var reviews []*domains.Review
	err := rr.db.SelectContext(ctx, &reviews, query, value)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (rr *reviewRepository) UpdateReview(ctx context.Context, review *domains.Review) error {
	query := `UPDATE reviews SET score = $1, comment = $2, image_url = $3, likes = $4, dislikes = $5 WHERE id = $6`
	result, err := rr.db.ExecContext(ctx, query, review.Score, review.Comment, review.ImageURL, review.Likes, review.Dislikes, review.ID)
	return CheckErrResult(result, err, ErrReviewNotFound)
}

func (rr *reviewRepository) DeleteReview(ctx context.Context, id string) error {
	query := `DELETE FROM reviews WHERE id = $1`
	result, err := rr.db.ExecContext(ctx, query, id)
	return CheckErrResult(result, err, ErrReviewNotFound)
}
