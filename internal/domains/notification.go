package domains

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        string           `json:"id" db:"id"`
	UserID    string           `json:"userId" db:"user_id"`
	Type      NotificationType `json:"type" db:"type"`
	Title     string           `json:"title" db:"title"`
	Message   string           `json:"message" db:"message"`
	Metadata  json.RawMessage  `json:"metadata" db:"metadata"`
	IsSeen    bool             `json:"isSeen" db:"is_seen"`
	SeenAt    *time.Time       `json:"seenAt,omitempty" db:"seen_at"`
	CreatedAt time.Time        `json:"createdAt" db:"created_at"`
}

type NotificationType string

const (
	OrderPlacedNotification    NotificationType = "order_placed"
	OrderShippedNotification   NotificationType = "order_shipped"
	OrderDeliveredNotification NotificationType = "order_delivered"
	OrderCancelledNotification NotificationType = "order_cancelled"
	PromotionNotification      NotificationType = "promotion"
	PasswordResetNotification  NotificationType = "password_reset"
	NewReviewNotification      NotificationType = "new_review"
	AbandonedCartNotification  NotificationType = "abandoned_cart"
	WelcomeNotification        NotificationType = "welcome"
	OtherNotification          NotificationType = "other"
	BackInStockNotification    NotificationType = "back_in_stock"
)

func NewNotification(userID string, notifType NotificationType, title, message string, metadata json.RawMessage) *Notification {
	return &Notification{
		ID:       uuid.NewString(),
		UserID:   userID,
		Type:     notifType,
		Title:    title,
		Message:  message,
		Metadata: metadata,
		IsSeen:   false,
	}
}
func (t NotificationType) IsValid() bool {
	switch t {
	case OrderPlacedNotification, OrderShippedNotification, OrderDeliveredNotification,
		OrderCancelledNotification, PromotionNotification, PasswordResetNotification,
		NewReviewNotification, AbandonedCartNotification, WelcomeNotification,
		BackInStockNotification, OtherNotification:
		return true
	default:
		return false
	}
}
