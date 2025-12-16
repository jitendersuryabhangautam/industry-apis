// Package repository provides database access layer implementations.
// Repositories handle all direct database operations using SQL queries.
package repository

import (
	"context"
	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentRepository provides database access for payment operations.
type PaymentRepository struct {
	db *pgxpool.Pool // Database connection pool
}

// PaymentRepo defines the methods used by services for payment operations.
type PaymentRepo interface {
	InitiatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	UpdatePayment(ctx context.Context, payment *models.Payment, id int) (*models.Payment, error)
}

// NewPaymentRepository creates and returns a new instance of PaymentRepository.
// It accepts a database connection pool for executing database operations.
func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// InitiatePayment inserts a new payment record into the database.
// It creates a new payment entry with the initial status and returns the payment with ID and creation timestamp.
// Returns the created payment with ID and CreatedAt fields populated, or an error if the operation fails.
func (r *PaymentRepository) InitiatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	query := `
		INSERT INTO payments (booking_id, amount, payment_method, transaction_id, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	// Execute the insert query and scan the returned ID and timestamp
	err := r.db.QueryRow(ctx, query, payment.BookingID, payment.Amount, payment.PaymentMethod, payment.TransactionID, payment.Status).Scan(&payment.ID, &payment.CreatedAt)

	if err != nil {
		return nil, err
	}
	return payment, nil

}

// UpdatePayment updates an existing payment record with card details and payment status.
// It updates the card information and marks the payment as paid.
// Returns the updated payment with all transaction details, or an error if the operation fails.
func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment, id int) (*models.Payment, error) {
	query := `
	UPDATE payments
	SET card_last4 = $1, card_brand=$2, receipt_url=$3, status=$4
	WHERE id = $5
	RETURNING id, booking_id, payment_method, amount, status, processed_at
	`
	// Execute the update query and scan the returned payment details
	pmt := r.db.QueryRow(ctx, query, payment.CardLastFour, payment.CardBrand, payment.ReceitURL, payment.Status, id)
	err := pmt.Scan(&payment.ID, &payment.BookingID, &payment.PaymentMethod, &payment.Amount, &payment.Status, &payment.ProcessedAt)
	if err != nil {
		return nil, err
	}
	return payment, nil

}
