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
		{name: "categories", ddl: createCategoriesTable},
		{name: "products", ddl: createProductsTable},
		{name: "notifications", ddl: createNotificationsTable},
		{name: "shopping_cart", ddl: createShoppingCartTable},
		{name: "favorites", ddl: createFavoritesTable},
		{name: "reviews", ddl: createReviewsTable},
		{name: "coupons", ddl: createCouponsTable},
		{name: "addresses", ddl: createAddressesTable},
		{name: "orders", ddl: createOrdersTable},
		{name: "order_items", ddl: createOrderItemsTable},
		{name: "coupon_usages", ddl: createCouponUsagesTable},
		{name: "payments", ddl: createPaymentsTable},
		{name: "inventory", ddl: createInventoryTable},
		{name: "review_votes", ddl: createReviewVotesTable},
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
