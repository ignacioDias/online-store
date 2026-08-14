package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

var ErrSessionNotFound = errors.New("Session not found")

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domains.Session) error
	GetSessionByID(ctx context.Context, id string) (*domains.Session, error)
	DeleteSessionByID(ctx context.Context, id string) error
	DeleteSessionsByUserID(ctx context.Context, userID string) error
}

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (sessRepo *sessionRepository) CreateSession(ctx context.Context, session *domains.Session) error {
	query := `
	INSERT INTO sessions (id, user_id, expires_at)
	VALUES (:id, :user_id, :expires_at)
	RETURNING created_at
	`
	stmt, err := sessRepo.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.GetContext(ctx, session, session)
}

func (sessRepo *sessionRepository) GetSessionByID(ctx context.Context, id string) (*domains.Session, error) {
	var session domains.Session
	query := "SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = $1 AND expires_at > (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')"
	if err := sessRepo.db.GetContext(ctx, &session, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (sessRepo *sessionRepository) DeleteSessionByID(ctx context.Context, id string) error {
	query := "DELETE FROM sessions WHERE id = $1"
	result, err := sessRepo.db.ExecContext(ctx, query, id)
	return CheckErrResult(result, err, ErrSessionNotFound)
}

func (sessRepo *sessionRepository) DeleteSessionsByUserID(ctx context.Context, userID string) error {
	query := "DELETE FROM sessions WHERE user_id = $1"
	result, err := sessRepo.db.ExecContext(ctx, query, userID)
	return CheckErrResult(result, err, ErrSessionNotFound)
}
