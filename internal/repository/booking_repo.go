package repository

import (
	"context"
	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}

func (b *BookingRepository) AddBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	query := `
	INSERT INTO bookings (user_id, room_id, check_in_date, check_out_date, adults, children, total_amount, status, payment_status, special_requests)
	Values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, created_at
	`
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
