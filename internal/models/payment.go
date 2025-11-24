package models

import "time"

type Payment struct {
	ID            int        `json:"id"`
	BookingID     int        `json:"booking_id"`
	Amount        int64      `json:"amount"`
	PaymentMethod string     `json:"payment_method"`
	TransactionID *string    `json:"transaction_id"`
	Status        string     `json:"status"`
	CardLastFour  *string    `json:"card_last4"`
	CardBrand     *string    `json:"card_brand"`
	ReceitURL     *string    `json:"receipt_url"`
	ProcessedAt   *time.Time `json:"processed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type PaymentRequest struct {
	BookingID     int    `json:"booking_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	TransactionID string `json:"transaction_id" binding:"required"`
}

type UpdatePaymentRequest struct {
	CardLastFour string `json:"card_last4" binding:"required"`
	CardBrand    string `json:"card_brand" binding:"required"`
	ReceiptURL   string `json:"receipt_url"`
}
