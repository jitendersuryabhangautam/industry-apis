// Package service provides business logic and service layer operations for the application
package service

import (
	"context"
	"errors"
	"industry-api/internal/models"
	"industry-api/internal/repository"
)

// BookingService handles all booking-related business logic operations.
// It acts as an intermediary between the handler and repository layers.
type BookingService struct {
	repo repository.BookingRepo // interface for booking data access (mockable)
}

// NewBookingService creates and returns a new instance of BookingService.
// It accepts a BookingRepo interface for data access operations.
func NewBookingService(repo repository.BookingRepo) *BookingService {
	return &BookingService{repo: repo}
}

// AddBooking creates a new booking after validating all required fields.
// It validates the booking data before delegating to the repository for persistence.
// Returns the created booking or an error if validation fails or database operation fails.
func (s *BookingService) AddBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	// Validate UserID is provided
	if booking.UserID == 0 {
		return nil, errors.New("user id is required")
	}
	// Validate RoomID is provided
	if booking.RoomID == 0 {
		return nil, errors.New("room id is required")
	}
	// Validate CheckInDate is provided
	if booking.CheckInDate.IsZero() {
		return nil, errors.New("check in date is required")
	}
	// Validate CheckOutDate is provided
	if booking.CheckOutDate.IsZero() {
		return nil, errors.New("check out date is required")
	}
	// Validate Adults count is provided
	if booking.Adults == 0 {
		return nil, errors.New("adults is required")
	}
	// Validate Children count is provided
	if booking.Children == 0 {
		return nil, errors.New("children is required")
	}
	// Validate TotalAmount is provided
	if booking.TotalAmount == 0 {
		return nil, errors.New("total amount is required")
	}
	// Validate Status is provided
	if booking.Status == "" {
		return nil, errors.New("status is required")
	}
	// Validate PaymentStatus is provided
	if booking.PaymentStatus == "" {
		return nil, errors.New("payment status is required")
	}
	// Persist the booking to the database through the repository
	booking, err := s.repo.AddBooking(ctx, booking)
	if err != nil {
		return nil, err
	}
	return booking, nil

}
