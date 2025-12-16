package service

import (
	"context"
	"testing"
	"time"

	"industry-api/internal/models"
)

type mockPaymentRepo struct {
	init   func(ctx context.Context, p *models.Payment) (*models.Payment, error)
	update func(ctx context.Context, p *models.Payment, id int) (*models.Payment, error)
}

func (m *mockPaymentRepo) InitiatePayment(ctx context.Context, p *models.Payment) (*models.Payment, error) {
	return m.init(ctx, p)
}
func (m *mockPaymentRepo) UpdatePayment(ctx context.Context, p *models.Payment, id int) (*models.Payment, error) {
	return m.update(ctx, p, id)
}

func TestInitiatePayment_Validation(t *testing.T) {
	svc := &PaymentService{repo: &mockPaymentRepo{init: func(ctx context.Context, p *models.Payment) (*models.Payment, error) { p.ID = 1; return p, nil }}}
	p := &models.Payment{BookingID: 0, Amount: 0}
	if _, err := svc.InitiatePaymet(context.Background(), p); err == nil {
		t.Fatalf("expected validation error")
	}

	p = &models.Payment{BookingID: 1, Amount: 100, PaymentMethod: "test"}
	got, err := svc.InitiatePaymet(context.Background(), p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID == 0 {
		t.Fatalf("expected ID set")
	}
}

func TestUpdatePayment_Validation(t *testing.T) {
	now := time.Now()
	card := "1234"
	brand := "Visa"
	svc := &PaymentService{repo: &mockPaymentRepo{update: func(ctx context.Context, p *models.Payment, id int) (*models.Payment, error) {
		p.ID = id
		p.ProcessedAt = &now
		return p, nil
	}}}
	p := &models.Payment{CardLastFour: nil, CardBrand: nil}
	if _, err := svc.UpdatePayment(context.Background(), p, 1); err == nil {
		t.Fatalf("expected validation error for missing card info")
	}

	p.CardLastFour = &card
	p.CardBrand = &brand
	got, err := svc.UpdatePayment(context.Background(), p, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 2 {
		t.Fatalf("expected updated id 2")
	}
}
