package service

import (
	"context"
	"errors"
	"industry-api/internal/models"
	"industry-api/internal/repository"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) AddBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	if booking.UserID == 0 {
		return nil, errors.New("user id is required")
	}
	if booking.RoomID == 0 {
		return nil, errors.New("room id is required")
	}
	if booking.CheckInDate.IsZero() {
		return nil, errors.New("check in date is required")
	}
	if booking.CheckOutDate.IsZero() {
		return nil, errors.New("check out date is required")
	}
	if booking.Adults == 0 {
		return nil, errors.New("adults is required")
	}
	if booking.Children == 0 {
		return nil, errors.New("children is required")
	}
	if booking.TotalAmount == 0 {
		return nil, errors.New("total amount is required")
	}
	if booking.Status == "" {
		return nil, errors.New("status is required")
	}
	if booking.PaymentStatus == "" {
		return nil, errors.New("payment status is required")
	}
	booking, err := s.repo.AddBooking(ctx, booking)
	if err != nil {
		return nil, err
	}
	return booking, nil

}
