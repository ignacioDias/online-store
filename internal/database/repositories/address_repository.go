package repositories

import (
	"context"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type AddressRepository interface {
	CreateAddress(ctx context.Context, address *domains.Address) error
	GetAddress(ctx context.Context, id string) (*domains.Address, error)
	GetAddressesByUserID(ctx context.Context, userID string) ([]*domains.Address, error)
	UpdateAddress(ctx context.Context, address *domains.Address) error
	DeleteAddress(ctx context.Context, id string) error
}

type addressRepository struct {
	db *sqlx.DB
}

var ErrAddressNotFound = errors.New("Address not found")

func NewAddressRepository(db *sqlx.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) CreateAddress(ctx context.Context, address *domains.Address) error {
	query := `INSERT INTO addresses (id, user_id, street, apartment, city, state, postal_code, country, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING created_at`
	return r.db.GetContext(ctx, &address.CreatedAt, query, address.ID, address.UserID, address.Street, address.Apartment, address.City, address.State, address.PostalCode, address.Country, address.IsDefault)

}
func (r *addressRepository) GetAddress(ctx context.Context, id string) (*domains.Address, error) {
	query := `SELECT id, user_id, street, apartment, city, state, postal_code, country, is_default, created_at FROM addresses WHERE id = $1`
	var address domains.Address
	err := r.db.GetContext(ctx, &address, query, id)
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) GetAddressesByUserID(ctx context.Context, userID string) ([]*domains.Address, error) {
	query := `SELECT id, user_id, street, apartment, city, state, postal_code, country, is_default, created_at FROM addresses WHERE user_id = $1`
	var addresses []*domains.Address
	err := r.db.SelectContext(ctx, &addresses, query, userID)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, address *domains.Address) error {
	query := `UPDATE addresses SET street = $1, apartment = $2, city = $3, state = $4, postal_code = $5, country = $6, is_default = $7 WHERE id = $8`
	result, err := r.db.ExecContext(ctx, query, address.Street, address.Apartment, address.City, address.State, address.PostalCode, address.Country, address.IsDefault, address.ID)
	return CheckErrResult(result, err, ErrAddressNotFound)
}

func (r *addressRepository) DeleteAddress(ctx context.Context, id string) error {
	query := `DELETE FROM addresses WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	return CheckErrResult(result, err, ErrAddressNotFound)
}
