package repositories

import (
	"context"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type CouponRepository interface {
	CreateCoupon(ctx context.Context, coupon *domains.Coupon) error
	GetCouponByID(ctx context.Context, id string) (*domains.Coupon, error)
	GetCouponByCode(ctx context.Context, code string) (*domains.Coupon, error)
	UseCoupon(ctx context.Context, id string) error
	DeactivateCoupon(ctx context.Context, id string) error
	UpdateCoupon(ctx context.Context, coupon *domains.Coupon) error
	DeleteCoupon(ctx context.Context, id string) error
	ListCoupons(ctx context.Context) ([]*domains.Coupon, error)
}

type couponRepository struct {
	db *sqlx.DB
}

var ErrCouponUsageLimitReached = errors.New("coupon usage limit reached")
var ErrCouponNotFound = errors.New("coupon not found")

func NewCouponRepository(db *sqlx.DB) CouponRepository {
	return &couponRepository{db: db}
}

func (cr *couponRepository) CreateCoupon(ctx context.Context, coupon *domains.Coupon) error {
	query := `INSERT INTO coupons (id, code, discount_type, discount_value, min_purchase_amount, max_discount_amount, usage_limit, usage_count, usage_limit_per_user, starts_at, expires_at, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING created_at`
	return cr.db.QueryRowContext(ctx, query, coupon.ID, coupon.Code, coupon.DiscountType, coupon.DiscountValue, coupon.MinPurchaseAmount, coupon.MaxDiscountAmount, coupon.UsageLimit, coupon.UsageCount, coupon.UsageLimitPerUser, coupon.StartsAt, coupon.ExpiresAt, coupon.IsActive).Scan(&coupon.CreatedAt)
}

func (cr *couponRepository) ListCoupons(ctx context.Context) ([]*domains.Coupon, error) {
	query := `SELECT id, code, discount_type, discount_value, min_purchase_amount, max_discount_amount, usage_limit, usage_count, usage_limit_per_user, starts_at, expires_at, is_active, created_at FROM coupons`
	var coupons []*domains.Coupon
	err := cr.db.SelectContext(ctx, &coupons, query)
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

func (cr *couponRepository) GetCouponByID(ctx context.Context, id string) (*domains.Coupon, error) {
	query := `SELECT id, code, discount_type, discount_value, min_purchase_amount, max_discount_amount, usage_limit, usage_count, usage_limit_per_user, starts_at, expires_at, is_active, created_at FROM coupons WHERE id = $1`
	var coupon domains.Coupon
	err := cr.db.GetContext(ctx, &coupon, query, id)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (cr *couponRepository) GetCouponByCode(ctx context.Context, code string) (*domains.Coupon, error) {
	query := `SELECT id, code, discount_type, discount_value, min_purchase_amount, max_discount_amount, usage_limit, usage_count, usage_limit_per_user, starts_at, expires_at, is_active, created_at FROM coupons WHERE code = $1`
	var coupon domains.Coupon
	err := cr.db.GetContext(ctx, &coupon, query, code)
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (cr *couponRepository) UseCoupon(ctx context.Context, id string) error {
	query := `UPDATE coupons
		SET usage_count = usage_count + 1
		WHERE id = $1
		  AND is_active = TRUE
		  AND NOW() >= starts_at
		  AND NOW() < expires_at
		  AND (usage_limit IS NULL OR usage_count < usage_limit)`
	result, err := cr.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrCouponUsageLimitReached
	}
	return nil
}

func (cr *couponRepository) DeactivateCoupon(ctx context.Context, id string) error {
	query := `UPDATE coupons SET is_active = FALSE WHERE id = $1`
	result, err := cr.db.ExecContext(ctx, query, id)
	return CheckErrResult(result, err, ErrCouponNotFound)
}

func (cr *couponRepository) UpdateCoupon(ctx context.Context, coupon *domains.Coupon) error {
	query := `UPDATE coupons SET code = $1, discount_type = $2, discount_value = $3, min_purchase_amount = $4, max_discount_amount = $5, usage_limit = $6, usage_limit_per_user = $7, starts_at = $8, expires_at = $9, is_active = $10 WHERE id = $11`
	result, err := cr.db.ExecContext(ctx, query, coupon.Code, coupon.DiscountType, coupon.DiscountValue, coupon.MinPurchaseAmount, coupon.MaxDiscountAmount, coupon.UsageLimit, coupon.UsageLimitPerUser, coupon.StartsAt, coupon.ExpiresAt, coupon.IsActive, coupon.ID)
	return CheckErrResult(result, err, ErrCouponNotFound)
}

func (cr *couponRepository) DeleteCoupon(ctx context.Context, id string) error {
	query := `DELETE FROM coupons WHERE id = $1`
	result, err := cr.db.ExecContext(ctx, query, id)
	return CheckErrResult(result, err, ErrCouponNotFound)
}
