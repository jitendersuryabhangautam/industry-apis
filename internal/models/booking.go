package models

import "time"

type Booking struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	RoomID          int       `json:"room_id"`
	CheckInDate     time.Time `json:"check_in_date"`
	CheckOutDate    time.Time `json:"check_out_date"`
	Adults          int       `json:"adults"`
	Children        int       `json:"children"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status"`
	PaymentStatus   string    `json:"payment_status"`
	SpecialRequests string    `json:"special_requests"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type BookingRequest struct {
	UserID          int       `json:"user_id" binding:"required"`
	RoomID          int       `json:"room_id" binding:"required"`
	CheckInDate     time.Time `json:"check_in_date" binding:"required"`
	CheckOutDate    time.Time `json:"check_out_date" binding:"required"`
	Adults          int       `json:"adults" binding:"required"`
	Children        int       `json:"children" binding:"required"`
	SpecialRequests string    `json:"special_requests"`
	Status          string    `json:"status" binding:"required"`
	PaymentStatus   string    `json:"payment_status" binding:"required"`
	TotalAmount     float64   `json:"total_amount" binding:"required"`
}
