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
	"github.com/jmoiron/sqlx"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return sqlx.NewDb(db, "sqlmock"), mock
}

func userRows(user *domains.User) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "username", "name", "email", "hashed_password", "profile_picture",
		"is_active", "email_verified", "created_at", "updated_at",
	}).AddRow(
		user.ID, user.Username, user.Name, user.Email, user.HashedPassword, user.ProfilePicture,
		user.IsActive, user.EmailVerified, user.CreatedAt, user.UpdatedAt,
	)
}

func TestUserRepositoryCreateUser(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	user := domains.NewUser("alex", "Alex Doe", "alex@example.com", "hashed")

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id,  username, name, email, hashed_password, profile_picture, is_active, email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")).WithArgs(user.ID, user.Username, user.Name, user.Email, user.HashedPassword, user.ProfilePicture, user.IsActive, user.EmailVerified).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepositoryGetUserByID(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	now := time.Now().UTC()
	want := &domains.User{ID: "user-1", Username: "alex", Name: "Alex Doe", Email: "alex@example.com", HashedPassword: "hashed", IsActive: true, CreatedAt: now, UpdatedAt: now}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, name, email, hashed_password, profile_picture, is_active, email_verified, created_at, updated_at FROM users WHERE id = $1")).WithArgs(want.ID).WillReturnRows(userRows(want))

	got, err := repo.GetUserByID(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}
	if got.ID != want.ID || got.Email != want.Email || !got.CreatedAt.Equal(now) {
		t.Fatalf("GetUserByID returned %+v, want %+v", got, want)
	}
}

func TestUserRepositoryGetUserNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, username, name, email, hashed_password, profile_picture, is_active, email_verified, created_at, updated_at FROM users WHERE email = $1")).WithArgs("missing@example.com").WillReturnError(sql.ErrNoRows)

	_, err := repo.GetUserByEmail(context.Background(), "missing@example.com")
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("error = %v, want ErrUserNotFound", err)
	}
}

func TestUserRepositoryUpdateAndDelete(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewUserRepository(db)
	user := domains.NewUser("alex", "Alex Doe", "alex@example.com", "hashed")

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET username = $1, name = $2, email = $3, hashed_password = $4, profile_picture = $5, is_active = $6, email_verified = $7, updated_at = NOW() WHERE id = $8")).WithArgs(user.Username, user.Name, user.Email, user.HashedPassword, user.ProfilePicture, user.IsActive, user.EmailVerified, user.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).WithArgs(user.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.UpdateUser(context.Background(), user); err != nil {
		t.Fatalf("UpdateUser returned error: %v", err)
	}
	if err := repo.DeleteUserByID(context.Background(), user.ID); err != nil {
		t.Fatalf("DeleteUserByID returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
