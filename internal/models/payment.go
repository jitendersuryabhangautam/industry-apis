// Package models defines all data structures used throughout the application.
package models

import "time"

// Payment represents a payment transaction record in the system.
type Payment struct {
	ID            int        `json:"id"`             // Unique payment identifier
	BookingID     int        `json:"booking_id"`     // ID of the associated booking
	Amount        int64      `json:"amount"`         // Payment amount in cents (to avoid floating point issues)
	PaymentMethod string     `json:"payment_method"` // Payment method used (e.g., "credit_card", "debit_card")
	TransactionID *string    `json:"transaction_id"` // External transaction ID from payment gateway
	Status        string     `json:"status"`         // Payment status (e.g., "pending", "completed", "failed")
	CardLastFour  *string    `json:"card_last4"`     // Last four digits of the card used
	CardBrand     *string    `json:"card_brand"`     // Brand of the card (e.g., "Visa", "Mastercard")
	ReceitURL     *string    `json:"receipt_url"`    // URL to the payment receipt
	ProcessedAt   *time.Time `json:"processed_at"`   // Timestamp when the payment was processed
	CreatedAt     time.Time  `json:"created_at"`     // Timestamp when the payment was created
}

// PaymentRequest represents the HTTP request body for initiating a payment.
type PaymentRequest struct {
	BookingID     int    `json:"booking_id" binding:"required"`     // ID of the booking to pay for
	Amount        int64  `json:"amount" binding:"required"`         // Amount to pay in cents
	PaymentMethod string `json:"payment_method" binding:"required"` // Payment method
	TransactionID string `json:"transaction_id" binding:"required"` // Transaction ID from payment gateway
}

// UpdatePaymentRequest represents the HTTP request body for updating payment details.
type UpdatePaymentRequest struct {
	CardLastFour string `json:"card_last4" binding:"required"` // Last 4 digits of card
	CardBrand    string `json:"card_brand" binding:"required"` // Card brand name
	ReceiptURL   string `json:"receipt_url"`                   // URL to payment receipt
}
