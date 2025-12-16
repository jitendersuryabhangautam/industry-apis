// Package models defines all data structures used throughout the application.
package models

import "time"

// Booking represents a hotel room booking record stored in the database.
type Booking struct {
	ID              int       `json:"id"`               // Unique booking identifier
	UserID          int       `json:"user_id"`          // ID of the user making the booking
	RoomID          int       `json:"room_id"`          // ID of the room being booked
	CheckInDate     time.Time `json:"check_in_date"`    // Date and time of guest arrival
	CheckOutDate    time.Time `json:"check_out_date"`   // Date and time of guest departure
	Adults          int       `json:"adults"`           // Number of adults in the booking
	Children        int       `json:"children"`         // Number of children in the booking
	TotalAmount     float64   `json:"total_amount"`     // Total cost of the booking
	Status          string    `json:"status"`           // Booking status (e.g., "confirmed", "cancelled")
	PaymentStatus   string    `json:"payment_status"`   // Payment status (e.g., "pending", "completed")
	SpecialRequests string    `json:"special_requests"` // Any special requests from the guest
	CreatedAt       time.Time `json:"created_at"`       // Timestamp when the booking was created
	UpdatedAt       time.Time `json:"updated_at"`       // Timestamp of the last update
}

// BookingRequest represents the HTTP request body for creating a new booking.
type BookingRequest struct {
	UserID          int       `json:"user_id" binding:"required"`        // ID of the user making the booking
	RoomID          int       `json:"room_id" binding:"required"`        // ID of the room to book
	CheckInDate     time.Time `json:"check_in_date" binding:"required"`  // Date of arrival
	CheckOutDate    time.Time `json:"check_out_date" binding:"required"` // Date of departure
	Adults          int       `json:"adults" binding:"required"`         // Number of adults
	Children        int       `json:"children" binding:"required"`       // Number of children
	SpecialRequests string    `json:"special_requests"`                  // Optional special requests
	Status          string    `json:"status" binding:"required"`         // Initial booking status
	PaymentStatus   string    `json:"payment_status" binding:"required"` // Initial payment status
	TotalAmount     float64   `json:"total_amount" binding:"required"`   // Total booking amount
}
