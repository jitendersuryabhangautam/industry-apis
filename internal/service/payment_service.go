// Package service provides business logic layer implementations.
// Services contain validation logic and business rules before delegating to repositories.
package service

import (
	"context"
	"errors"
	"industry-api/internal/models"
	"industry-api/internal/repository"
)

// PaymentService handles all payment-related business logic operations.
// It acts as an intermediary between handlers and repository layers.
type PaymentService struct {
	repo *repository.PaymentRepository // Repository layer for data access
}

// NewPaymentRepository creates and returns a new instance of PaymentService.
// It accepts a PaymentRepository dependency for data access operations.
func NewPaymentRepository(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

// InitiatePaymet initiates a new payment transaction after validation.
// It validates all required payment fields before delegating to the repository.
// Returns the created payment or an error if validation fails.
func (s *PaymentService) InitiatePaymet(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	// Validate payment amount is provided and greater than zero
	if payment.Amount == 0 {
		return nil, errors.New("amount is required")
	}
	// Validate booking ID is provided
	if payment.BookingID == 0 {
		return nil, errors.New("booking id is required")
	}
	// Validate payment method is specified
	if payment.PaymentMethod == "" {
		return nil, errors.New("payment method is required")
	}

	// Delegate to repository to persist the payment
	pmt, err := s.repo.InitiatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}
	return pmt, nil

}

// UpdatePayment updates an existing payment with card details and mark it as paid.
// It validates card information before delegating to the repository.
// Returns the updated payment or an error if validation fails.
func (s *PaymentService) UpdatePayment(ctx context.Context, payment *models.Payment, id int) (*models.Payment, error) {
	// Validate card last 4 digits are provided
	if payment.CardLastFour == nil {
		return nil, errors.New("card last 4 is required")
	}
	// Validate card brand is specified
	if payment.CardBrand == nil {
		return nil, errors.New("card brand is required")
	}
	// Delegate to repository to update the payment
	pmt, err := s.repo.UpdatePayment(ctx, payment, id)
	if err != nil {
		return nil, err
	}
	return pmt, nil
}
