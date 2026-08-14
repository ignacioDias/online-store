package repositories

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"sports-store/internal/domains"

	"github.com/DATA-DOG/go-sqlmock"
)

func notificationRows(notification *domains.Notification) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "user_id", "type", "title", "message", "metadata", "is_seen", "seen_at", "created_at"}).AddRow(
		notification.ID, notification.UserID, notification.Type, notification.Title, notification.Message,
		[]byte(notification.Metadata), notification.IsSeen, notification.SeenAt, notification.CreatedAt,
	)
}

func TestNotificationRepositoryCreateAndGet(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewNotificationRepository(db)
	seenAt := time.Now().UTC()
	notification := domains.NewNotification("user-1", domains.OrderPlacedNotification, "Order placed", "Your order was placed", []byte(`{"order_id":"order-1"}`))
	notification.IsSeen = true
	notification.SeenAt = &seenAt
	notification.CreatedAt = seenAt

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO notifications (id, user_id, type, title, message, metadata, is_seen, seen_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")).WithArgs(notification.ID, notification.UserID, notification.Type, notification.Title, notification.Message, notification.Metadata, notification.IsSeen, notification.SeenAt).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, type, title, message, metadata, is_seen, seen_at, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3")).WithArgs(notification.UserID, 10, 0).WillReturnRows(notificationRows(notification))

	if err := repo.CreateNotification(context.Background(), notification); err != nil {
		t.Fatalf("CreateNotification returned error: %v", err)
	}
	got, err := repo.GetNotificationsByUserID(context.Background(), notification.UserID, 10, 0)
	if err != nil {
		t.Fatalf("GetNotificationsByUserID returned error: %v", err)
	}
	if len(got) != 1 || got[0].ID != notification.ID || !got[0].CreatedAt.Equal(seenAt) {
		t.Fatalf("GetNotificationsByUserID returned %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestNotificationRepositoryEmptyAndNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewNotificationRepository(db)
	selectQuery := regexp.QuoteMeta("SELECT id, user_id, type, title, message, metadata, is_seen, seen_at, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3")
	mock.ExpectQuery(selectQuery).WithArgs("user-1", 10, 0).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "type", "title", "message", "metadata", "is_seen", "seen_at", "created_at"}))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE notifications SET is_seen = TRUE, seen_at = NOW() WHERE id = $1")).WithArgs("missing").WillReturnResult(sqlmock.NewResult(0, 0))

	items, err := repo.GetNotificationsByUserID(context.Background(), "user-1", 10, 0)
	if len(items) != 0 || err == nil || err.Error() != "No notifications" {
		t.Fatalf("empty result = (%v, %v), want empty slice and No notifications", items, err)
	}
	if err := repo.MarkNotificationAsRead(context.Background(), "missing"); !errors.Is(err, ErrNotificationNotFound) {
		t.Fatalf("error = %v, want ErrNotificationNotFound", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestNotificationRepositoryMarkReadAndDelete(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewNotificationRepository(db)
	mock.ExpectExec(regexp.QuoteMeta("UPDATE notifications SET is_seen = TRUE, seen_at = NOW() WHERE id = $1")).WithArgs("notification-1").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM notifications WHERE id = $1")).WithArgs("notification-1").WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.MarkNotificationAsRead(context.Background(), "notification-1"); err != nil {
		t.Fatalf("MarkNotificationAsRead returned error: %v", err)
	}
	if err := repo.DeleteNotification(context.Background(), "notification-1"); err != nil {
		t.Fatalf("DeleteNotification returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
