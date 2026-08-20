package domains

import (
	"time"

	"github.com/google/uuid"
)

type Discounts string

const (
	Percentage  Discounts = "percentage"
	FixedAmount Discounts = "fixed_amount"
)

type Coupon struct {
	ID                string    `json:"id" db:"id"`
	Code              string    `json:"code" db:"code"`
	DiscountType      Discounts `json:"discount_type" db:"discount_type"`
	DiscountValue     float64   `json:"discount_value" db:"discount_value"`
	MinPurchaseAmount *float64  `json:"min_purchase_amount,omitempty" db:"min_purchase_amount"`
	MaxDiscountAmount *float64  `json:"max_discount_amount,omitempty" db:"max_discount_amount"`
	UsageLimit        *int      `json:"usage_limit,omitempty" db:"usage_limit"`
	UsageCount        int       `json:"usage_count" db:"usage_count"`
	UsageLimitPerUser int       `json:"usage_limit_per_user" db:"usage_limit_per_user"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	StartsAt          time.Time `json:"starts_at" db:"starts_at"`
	ExpiresAt         time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

func NewCoupon(code string, discountType Discounts, discountValue float64, minPurchaseAmount *float64, maxDiscountAmount *float64, usageLimit *int, usageLimitPerUser int, startsAt time.Time, expiresAt time.Time) *Coupon {
	return &Coupon{
		ID:                uuid.NewString(),
		Code:              code,
		DiscountType:      discountType,
		DiscountValue:     discountValue,
		MinPurchaseAmount: minPurchaseAmount,
		MaxDiscountAmount: maxDiscountAmount,
		UsageLimit:        usageLimit,
		UsageCount:        0,
		UsageLimitPerUser: usageLimitPerUser,
		IsActive:          true,
		StartsAt:          startsAt,
		ExpiresAt:         expiresAt,
	}
}
