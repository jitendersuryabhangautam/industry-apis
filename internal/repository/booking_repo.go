// Package repository provides database access layer implementations.
// Repositories handle all direct database operations using SQL queries.
package repository

import (
	"context"
	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// BookingRepository provides database access for booking operations.
type BookingRepository struct {
	db *pgxpool.Pool // Database connection pool
}

// BookingRepo defines the methods used by services for booking operations.
type BookingRepo interface {
	AddBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error)
}

// NewBookingRepository creates and returns a new instance of BookingRepository.
// It accepts a database connection pool for executing database operations.
func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

// AddBooking inserts a new booking record into the database.
// It executes an INSERT query and returns the generated booking ID and creation timestamp.
// Returns the created booking with ID and CreatedAt fields populated, or an error if the operation fails.
func (b *BookingRepository) AddBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	query := `
	INSERT INTO bookings (user_id, room_id, check_in_date, check_out_date, adults, children, total_amount, status, payment_status, special_requests)
	Values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, created_at
	`
	// Execute the insert query and scan the returned ID and timestamp
	err := b.db.QueryRow(ctx, query,
		booking.UserID,
		booking.RoomID,
		booking.CheckInDate,
		booking.CheckOutDate,
		booking.Adults,
		booking.Children,
		booking.TotalAmount,
		booking.Status,
		booking.PaymentStatus,
		booking.SpecialRequests,
	).Scan(&booking.ID, &booking.CreatedAt)

	if err != nil {
		return nil, err
	}
	return booking, nil
}
