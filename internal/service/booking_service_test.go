package service

import (
	"context"
	"testing"
	"time"

	"industry-api/internal/models"
)

type mockBookingRepo struct {
	add func(ctx context.Context, b *models.Booking) (*models.Booking, error)
}

func (m *mockBookingRepo) AddBooking(ctx context.Context, b *models.Booking) (*models.Booking, error) {
	return m.add(ctx, b)
}

func TestAddBooking_ValidationErrors(t *testing.T) {
	svc := &BookingService{repo: &mockBookingRepo{add: func(ctx context.Context, b *models.Booking) (*models.Booking, error) { return b, nil }}}

	b := &models.Booking{}
	if _, err := svc.AddBooking(context.Background(), b); err == nil {
		t.Fatalf("expected validation error for missing fields")
	}

	b = &models.Booking{UserID: 1, RoomID: 1, CheckInDate: time.Now(), CheckOutDate: time.Now().Add(24 * time.Hour), Adults: 1, Children: 1, TotalAmount: 100, Status: "pending", PaymentStatus: "pending"}
	got, err := svc.AddBooking(context.Background(), b)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if got == nil || got.UserID != 1 {
		t.Fatalf("unexpected result: %+v", got)
	}
}
