package service

import (
	"context"
	"errors"
	"industry-api/internal/models"
	"industry-api/internal/repository"
)

type PaymentService struct {
	repo *repository.PaymentRepository
}

func NewPaymentRepository(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) InitiatePaymet(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	if payment.Amount == 0 {
		return nil, errors.New("amount is required")
	}
	if payment.BookingID == 0 {
		return nil, errors.New("booking id is required")
	}
	if payment.PaymentMethod == "" {
		return nil, errors.New("payment method is required")
	}

	pmt, err := s.repo.InitiatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}
	return pmt, nil

}

func (s *PaymentService) UpdatePayment(ctx context.Context, payment *models.Payment, id int) (*models.Payment, error) {
	if payment.CardLastFour == nil {
		return nil, errors.New("card last 4 is required")
	}
	if payment.CardBrand == nil {
		return nil, errors.New("card brand is required")
	}
	pmt, err := s.repo.UpdatePayment(ctx, payment, id)
	if err != nil {
		return nil, err
	}
	return pmt, nil
}
