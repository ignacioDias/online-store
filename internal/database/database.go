package database

import (
	"fmt"
	"sports-store/internal/database/repositories"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB                     *sqlx.DB
	UserRepository         repositories.UserRepository
	SessionRepository      repositories.SessionRepository
	ProductRepository      repositories.ProductRepository
	NotificationRepository repositories.NotificationRepository
	ShoppingCartRepository repositories.ShoppingCartRepository
	FavoriteRepository     repositories.FavoriteRepository
	ReviewRepository       repositories.ReviewRepository
	InventoryRepository    repositories.InventoryRepository
	PaymentRepository      repositories.PaymentRepository
	CategoryRepository     repositories.CategoryRepository
	CouponRepository       repositories.CouponRepository
	AddressRepository      repositories.AddressRepository
	OrderRepository        repositories.OrderRepository //orders and orderitems tables
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		DB:                     db,
		UserRepository:         repositories.NewUserRepository(db),
		SessionRepository:      repositories.NewSessionRepository(db),
		ProductRepository:      repositories.NewProductRepository(db),
		NotificationRepository: repositories.NewNotificationRepository(db),
		ShoppingCartRepository: repositories.NewShoppingCartRepository(db),
		FavoriteRepository:     repositories.NewFavoriteRepository(db),
		ReviewRepository:       repositories.NewReviewRepository(db),
		InventoryRepository:    repositories.NewInventoryRepository(db),
		PaymentRepository:      repositories.NewPaymentRepository(db),
		CategoryRepository:     repositories.NewCategoryRepository(db),
		CouponRepository:       repositories.NewCouponRepository(db),
		AddressRepository:      repositories.NewAddressRepository(db),
		OrderRepository:        repositories.NewOrderRepository(db),
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

var createNotificationsTable = `CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,       
    title VARCHAR(255) NOT NULL,    
    message TEXT NOT NULL,          
    metadata JSONB,                   -- data extra según el tipo (ej: order_id)
    is_seen BOOLEAN DEFAULT FALSE,
    seen_at TIMESTAMPTZ,                
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_notifications_user_unseen ON notifications(user_id, is_seen);
`

var createFavoritesTable = `CREATE TABLE IF NOT EXISTS favorites (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, product_id)
);
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
