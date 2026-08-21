package domains

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusApproved   PaymentStatus = "approved"
	PaymentStatusRejected   PaymentStatus = "rejected"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)

type Payment struct {
	ID                string        `json:"id" db:"id"`
	OrderID           string        `json:"order_id" db:"order_id"`
	UserID            string        `json:"user_id" db:"user_id"`
	Provider          string        `json:"provider" db:"provider"`                       // 'mercadopago', 'stripe'
	ProviderPaymentID string        `json:"provider_payment_id" db:"provider_payment_id"` // el ID que te devuelve el proveedor
	Status            PaymentStatus `json:"status" db:"status"`                           // 'pending', 'approved', 'rejected', 'refunded'
	Amount            float64       `json:"amount" db:"amount"`
	Currency          string        `json:"currency" db:"currency"`                       // default: 'ARS'
	PaymentMethod     *string       `json:"payment_method,omitempty" db:"payment_method"` // 'credit_card', 'debit_card', 'cash', etc. (te lo da el proveedor)
	Metadata          *string       `json:"metadata,omitempty" db:"metadata"`             // data extra que te devuelva el proveedor
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
}

func NewPayment(orderID, userID, provider, providerPaymentID string, status PaymentStatus, amount float64, currency string, paymentMethod, metadata *string) *Payment {
	return &Payment{
		ID:                uuid.NewString(),
		OrderID:           orderID,
		UserID:            userID,
		Provider:          provider,
		ProviderPaymentID: providerPaymentID,
		Status:            status,
		Amount:            amount,
		Currency:          currency,
		PaymentMethod:     paymentMethod,
		Metadata:          metadata,
	}
}
