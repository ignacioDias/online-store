package repositories

import (
	"context"
	"errors"
	"fmt"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *domains.Notification) error
	GetNotificationsByUserID(ctx context.Context, userID string, limit, offset int) ([]*domains.Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID string) error
	DeleteNotification(ctx context.Context, notificationID string) error
}

type notificationRepository struct {
	db *sqlx.DB
}

var ErrNotificationNotFound = errors.New("Notification not found")

func NewNotificationRepository(db *sqlx.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (nr *notificationRepository) CreateNotification(ctx context.Context, notification *domains.Notification) error {
	query := `INSERT INTO notifications (id, user_id, type, title, message, metadata, is_seen, seen_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING created_at`
	err := nr.db.GetContext(ctx, &notification.CreatedAt, query, notification.ID, notification.UserID, notification.Type, notification.Title, notification.Message, notification.Metadata, notification.IsSeen, notification.SeenAt)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}
	return nil
}

func (nr *notificationRepository) GetNotificationsByUserID(ctx context.Context, userID string, limit, offset int) ([]*domains.Notification, error) {
	query := `SELECT id, user_id, type, title, message, metadata, is_seen, seen_at, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	var notifications []*domains.Notification
	err := nr.db.SelectContext(ctx, &notifications, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select notifications: %w", err)
	}
	if len(notifications) == 0 {
		return notifications, ErrNotificationNotFound
	}
	return notifications, nil
}

func (nr *notificationRepository) MarkNotificationAsRead(ctx context.Context, notificationID string) error {
	query := `UPDATE notifications SET is_seen = TRUE, seen_at = NOW() WHERE id = $1`
	result, err := nr.db.ExecContext(ctx, query, notificationID)
	return CheckErrResult(result, err, ErrNotificationNotFound)
}
func (nr *notificationRepository) DeleteNotification(ctx context.Context, notificationID string) error {
	query := `DELETE FROM notifications WHERE id = $1`
	result, err := nr.db.ExecContext(ctx, query, notificationID)
	return CheckErrResult(result, err, ErrNotificationNotFound)
}
