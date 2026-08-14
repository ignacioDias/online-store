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

func TestFavoriteRepositoryAddDuplicate(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewFavoriteRepository(db)
	favorite := domains.NewFavorite("user-1", "product-1")

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO favorites (user_id, product_id) VALUES ($1, $2) ON CONFLICT (user_id, product_id) DO NOTHING RETURNING created_at")).
		WithArgs(favorite.UserID, favorite.ProductID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}))

	if err := repo.AddFavorite(context.Background(), favorite); !errors.Is(err, ErrFavoriteAlreadyExists) {
		t.Fatalf("error = %v, want ErrFavoriteAlreadyExists", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFavoriteRepositoryAdd(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewFavoriteRepository(db)
	favorite := domains.NewFavorite("user-1", "product-1")
	createdAt := time.Now().UTC()

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO favorites (user_id, product_id) VALUES ($1, $2) ON CONFLICT (user_id, product_id) DO NOTHING RETURNING created_at")).
		WithArgs(favorite.UserID, favorite.ProductID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(createdAt))

	if err := repo.AddFavorite(context.Background(), favorite); err != nil {
		t.Fatalf("AddFavorite returned error: %v", err)
	}
	if !favorite.CreatedAt.Equal(createdAt) {
		t.Fatalf("CreatedAt = %v, want %v", favorite.CreatedAt, createdAt)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFavoriteRepositoryRemove(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewFavoriteRepository(db)
	favorite := domains.NewFavorite("user-1", "product-1")

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM favorites WHERE user_id = $1 AND product_id = $2")).
		WithArgs(favorite.UserID, favorite.ProductID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.RemoveFavorite(context.Background(), favorite); err != nil {
		t.Fatalf("RemoveFavorite returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFavoriteRepositoryRemoveNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewFavoriteRepository(db)
	favorite := domains.NewFavorite("user-1", "missing-product")

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM favorites WHERE user_id = $1 AND product_id = $2")).
		WithArgs(favorite.UserID, favorite.ProductID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := repo.RemoveFavorite(context.Background(), favorite); !errors.Is(err, ErrFavoriteNotFound) {
		t.Fatalf("error = %v, want ErrFavoriteNotFound", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFavoriteRepositoryGetByUserID(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewFavoriteRepository(db)
	query := regexp.QuoteMeta("SELECT product_id FROM favorites WHERE user_id = $1")

	mock.ExpectQuery(query).
		WithArgs("user-1").
		WillReturnRows(sqlmock.NewRows([]string{"product_id"}).AddRow("product-1").AddRow("product-2"))

	productIDs, err := repo.GetFavoritesByUserID(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("GetFavoritesByUserID returned error: %v", err)
	}
	if len(productIDs) != 2 || productIDs[0] != "product-1" || productIDs[1] != "product-2" {
		t.Fatalf("product IDs = %v, want [product-1 product-2]", productIDs)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFavoriteRepositoryGetByUserIDEmpty(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewFavoriteRepository(db)
	query := regexp.QuoteMeta("SELECT product_id FROM favorites WHERE user_id = $1")

	mock.ExpectQuery(query).
		WithArgs("user-1").
		WillReturnRows(sqlmock.NewRows([]string{"product_id"}))

	productIDs, err := repo.GetFavoritesByUserID(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("GetFavoritesByUserID returned error: %v", err)
	}
	if productIDs == nil || len(productIDs) != 0 {
		t.Fatalf("product IDs = %v, want a non-nil empty slice", productIDs)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestFavoriteRepositoryGetByUserIDDatabaseError(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewFavoriteRepository(db)
	query := regexp.QuoteMeta("SELECT product_id FROM favorites WHERE user_id = $1")
	dbErr := errors.New("database unavailable")

	mock.ExpectQuery(query).WithArgs("user-1").WillReturnError(dbErr)

	_, err := repo.GetFavoritesByUserID(context.Background(), "user-1")
	if !errors.Is(err, dbErr) {
		t.Fatalf("error = %v, want wrapped database error", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
