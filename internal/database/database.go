package database

import "github.com/jmoiron/sqlx"

type Database struct {
	DB                     *sqlx.DB
	UserRepository         *UserRepository
	SessionRepository      *SessionRepository
	ProductRepository      *ProductRepository
	NotificationRepository *NotificationRepository
	CartRepository         *CartRepository
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		DB: db,
	}
}
