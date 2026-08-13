package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("User not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user *domains.User) error
	GetUserByID(ctx context.Context, id string) (*domains.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domains.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domains.User, error)
	UpdateUser(ctx context.Context, user *domains.User) error
	DeleteUserByID(ctx context.Context, id string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *domains.User) error {
	query := `INSERT INTO users (id,  username, name, email, hashed_password, profile_picture, is_active, email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := ur.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Name,
		user.Email,
		user.HashedPassword,
		user.ProfilePicture,
		user.IsActive,
		user.EmailVerified,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (ur *userRepository) GetUserByID(ctx context.Context, id string) (*domains.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	return ur.getUser(ctx, query, id)
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*domains.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	return ur.getUser(ctx, query, email)
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*domains.User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	return ur.getUser(ctx, query, username)
}

func (ur *userRepository) getUser(ctx context.Context, query string, args ...interface{}) (*domains.User, error) {
	var user domains.User
	if err := ur.db.GetContext(ctx, &user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, user *domains.User) error {
	query := `UPDATE users SET username = $1, name = $2, email = $3, hashed_password = $4, profile_picture = $5, is_active = $6, email_verified = $7, updated_at = NOW() WHERE id = $8`
	result, err := ur.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Name,
		user.Email,
		user.HashedPassword,
		user.ProfilePicture,
		user.IsActive,
		user.EmailVerified,
		user.ID,
	)
	return CheckErrResult(result, err, ErrUserNotFound)
}

func (ur *userRepository) DeleteUserByID(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := ur.db.ExecContext(ctx, query, id)
	return CheckErrResult(result, err, ErrUserNotFound)
}
