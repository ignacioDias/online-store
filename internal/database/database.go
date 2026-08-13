package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB                     *sqlx.DB
	UserRepository         *UserRepository
	SessionRepository      *SessionRepository
	ProductRepository      *ProductRepository
	NotificationRepository *NotificationRepository
	ShoppingCartRepository *ShoppingCartRepository
	FavoriteRepository     *FavoriteRepository
	ReviewRepository       *ReviewRepository
	InventoryRepository    *InventoryRepository
	PaymentRepository      *PaymentRepository
	CategoryRepository     *CategoryRepository
	CouponRepository       *CouponRepository
	AddressRepository      *AddressRepository
	OrderRepository        *OrderRepository
	OrderItemRepository    *OrderItemRepository
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		DB:                     db,
		UserRepository:         NewUserRepository(db),
		SessionRepository:      NewSessionRepository(db),
		ProductRepository:      NewProductRepository(db),
		NotificationRepository: NewNotificationRepository(db),
		ShoppingCartRepository: NewShoppingCartRepository(db),
	}
}

func (db *Database) Init() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return fmt.Errorf("begin init db transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	tableStatements := []struct {
		name string
		ddl  string
	}{
		{name: "users", ddl: createUsersTable},
		{name: "sessions", ddl: createSessionsTable},
	}

	for _, tableStmt := range tableStatements {
		if _, err := tx.Exec(tableStmt.ddl); err != nil {
			return fmt.Errorf("create table %s: %w", tableStmt.name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit init db transaction: %w", err)
	}

	return nil
}

func (db *Database) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}
