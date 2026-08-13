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
	OrderRepository        *OrderRepository //orders and orderitems tables
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		DB:                     db,
		UserRepository:         NewUserRepository(db),
		SessionRepository:      NewSessionRepository(db),
		ProductRepository:      NewProductRepository(db),
		NotificationRepository: NewNotificationRepository(db),
		ShoppingCartRepository: NewShoppingCartRepository(db),
		FavoriteRepository:     NewFavoriteRepository(db),
		ReviewRepository:       NewReviewRepository(db),
		InventoryRepository:    NewInventoryRepository(db),
		PaymentRepository:      NewPaymentRepository(db),
		CategoryRepository:     NewCategoryRepository(db),
		CouponRepository:       NewCouponRepository(db),
		AddressRepository:      NewAddressRepository(db),
		OrderRepository:        NewOrderRepository(db),
	}
}

var createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);`

var createUsersTable = `CREATE TABLE IF NOT EXISTS users (
    id                UUID PRIMARY KEY,
    username          VARCHAR(50) UNIQUE NOT NULL,
    name              VARCHAR(255) NOT NULL,
    email             VARCHAR(255) UNIQUE NOT NULL,
    hashed_password   TEXT NOT NULL,
    profile_picture   TEXT,                         -- URL o path al archivo
    is_active         BOOLEAN NOT NULL DEFAULT TRUE, -- para soft-disable en vez de borrar
    email_verified    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
`

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
		{name: "products", ddl: createProductsTable},
		{name: "notifications", ddl: createNotificationsTable},
		{name: "shopping_cart", ddl: createShoppingCartTable},
		{name: "favorites", ddl: createFavoritesTable},
		{name: "reviews", ddl: createReviewsTable},
		{name: "inventory", ddl: createInventoryTable},
		{name: "payments", ddl: createPaymentsTable},
		{name: "categories", ddl: createCategoriesTable},
		{name: "coupons", ddl: createCouponsTable},
		{name: "addresses", ddl: createAddressesTable},
		{name: "orders", ddl: createOrdersTable},
		{name: "order_items", ddl: createOrderItemsTable},
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
