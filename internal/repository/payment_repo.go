package repository

import (
	"context"
	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) InitiatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	query := `
		INSERT INTO payments (booking_id, amount, payment_method, transaction_id, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query, payment.BookingID, payment.Amount, payment.PaymentMethod, payment.TransactionID, payment.Status).Scan(&payment.ID, &payment.CreatedAt)

	if err != nil {
		return nil, err
	}
	return payment, nil

}

func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment, id int) (*models.Payment, error) {
	query := `
	UPDATE payments
	SET card_last4 = $1, card_brand=$2, receipt_url=$3, status=$4
	WHERE id = $5
	RETURNING id, booking_id, payment_method, amount, status, processed_at
	`
	pmt := r.db.QueryRow(ctx, query, payment.CardLastFour, payment.CardBrand, payment.ReceitURL, payment.Status, id)
	err := pmt.Scan(&payment.ID, &payment.BookingID, &payment.PaymentMethod, &payment.Amount, &payment.Status, &payment.ProcessedAt)
	if err != nil {
		return nil, err
	}
	return payment, nil

}
