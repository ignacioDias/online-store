package repositories

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"sports-store/internal/domains"

	"github.com/DATA-DOG/go-sqlmock"
)

func sessionRows(session *domains.Session) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "user_id", "created_at", "expires_at"}).AddRow(
		session.ID, session.UserID, session.CreatedAt, session.ExpiresAt,
	)
}

func TestSessionRepositoryCreateAndGet(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSessionRepository(db)
	session := domains.NewSession("user-1")
	createdAt := time.Now().UTC()

	prepare := mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?) RETURNING created_at"))
	prepare.ExpectQuery().WithArgs(session.ID, session.UserID, session.ExpiresAt).WillReturnRows(
		sqlmock.NewRows([]string{"created_at"}).AddRow(createdAt),
	)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = $1 AND expires_at > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')")).WithArgs(session.ID).WillReturnRows(sessionRows(&domains.Session{ID: session.ID, UserID: session.UserID, CreatedAt: createdAt, ExpiresAt: session.ExpiresAt}))

	if err := repo.CreateSession(context.Background(), session); err != nil {
		t.Fatalf("CreateSession returned error: %v", err)
	}
	if !session.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %v, want %v", session.CreatedAt, createdAt)
	}
	got, err := repo.GetSessionByID(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("GetSessionByID returned error: %v", err)
	}
	if got.ID != session.ID || got.UserID != session.UserID {
		t.Fatalf("GetSessionByID returned %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSessionRepositoryGetNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSessionRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = $1 AND expires_at > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')")).WithArgs("missing").WillReturnError(sql.ErrNoRows)

	_, err := repo.GetSessionByID(context.Background(), "missing")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("error = %v, want ErrSessionNotFound", err)
	}
}

func TestSessionRepositoryDeleteMethods(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewSessionRepository(db)
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM sessions WHERE id = $1")).WithArgs("session-1").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM sessions WHERE user_id = $1")).WithArgs("user-1").WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.DeleteSessionByID(context.Background(), "session-1"); err != nil {
		t.Fatalf("DeleteSessionByID returned error: %v", err)
	}
	if err := repo.DeleteSessionsByUserID(context.Background(), "user-1"); err != nil {
		t.Fatalf("DeleteSessionsByUserID returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
