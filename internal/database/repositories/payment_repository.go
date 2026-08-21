package repositories

import (
	"context"
	"errors"
	"sports-store/internal/domains"

	"github.com/jmoiron/sqlx"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *domains.Payment) error
	GetPaymentByID(ctx context.Context, paymentID string) (*domains.Payment, error)
	GetPaymentsByUserID(ctx context.Context, userID string) ([]*domains.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID string) (*domains.Payment, error)
	GetPaymentsByProvider(ctx context.Context, provider string) ([]*domains.Payment, error)
	GetPaymentByProviderID(ctx context.Context, provider string, providerPaymentID string) (*domains.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID string, status domains.PaymentStatus) error
}

type paymentRepository struct {
	db *sqlx.DB
}

var ErrPaymentNotFound = errors.New("Payment not found")

const paymentColumns = `id, order_id, user_id, provider, provider_payment_id, status, amount, currency, payment_method, metadata, created_at, updated_at`

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (p *paymentRepository) CreatePayment(ctx context.Context, payment *domains.Payment) error {
	query := `INSERT INTO payments (id, order_id, user_id, provider, provider_payment_id, status, amount, currency, payment_method, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING updated_at, created_at`
	return p.db.QueryRowContext(ctx, query, payment.ID, payment.OrderID, payment.UserID, payment.Provider, payment.ProviderPaymentID, payment.Status, payment.Amount, payment.Currency, payment.PaymentMethod, payment.Metadata).Scan(&payment.UpdatedAt, &payment.CreatedAt)
}

func (p *paymentRepository) GetPaymentByID(ctx context.Context, paymentID string) (*domains.Payment, error) {
	return p.getPayment(ctx, "id = $1", paymentID)
}

func (p *paymentRepository) GetPaymentByOrderID(ctx context.Context, orderID string) (*domains.Payment, error) {
	return p.getPayment(ctx, "order_id = $1", orderID)
}

func (p *paymentRepository) GetPaymentByProviderID(ctx context.Context, provider string, providerPaymentID string) (*domains.Payment, error) {
	return p.getPayment(ctx, "provider = $1 AND provider_payment_id = $2", provider, providerPaymentID)
}

func (p *paymentRepository) getPayment(ctx context.Context, condition string, args ...interface{}) (*domains.Payment, error) {
	query := `SELECT ` + paymentColumns + ` FROM payments WHERE ` + condition
	var payment domains.Payment
	err := p.db.GetContext(ctx, &payment, query, args...)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	return &payment, nil
}

func (p *paymentRepository) GetPaymentsByProvider(ctx context.Context, provider string) ([]*domains.Payment, error) {
	return p.getPayments(ctx, "provider = $1", provider)
}

func (p *paymentRepository) GetPaymentsByUserID(ctx context.Context, userID string) ([]*domains.Payment, error) {
	return p.getPayments(ctx, "user_id = $1", userID)
}

func (p *paymentRepository) getPayments(ctx context.Context, condition string, args ...any) ([]*domains.Payment, error) {
	query := `SELECT ` + paymentColumns + ` FROM payments WHERE ` + condition
	var payments []*domains.Payment
	err := p.db.SelectContext(ctx, &payments, query, args...)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (p *paymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID string, status domains.PaymentStatus) error {
	query := `UPDATE payments SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := p.db.ExecContext(ctx, query, status, paymentID)
	return CheckErrResult(result, err, ErrPaymentNotFound)
}

func (p *paymentRepository) DeletePayment(ctx context.Context, paymentID string) error {
	query := `DELETE FROM payments WHERE id = $1`
	result, err := p.db.ExecContext(ctx, query, paymentID)
	return CheckErrResult(result, err, ErrPaymentNotFound)
}
