// Package handler contains HTTP request handlers for all API endpoints.
// Handlers receive HTTP requests, validate input, call services, and send responses.
package handler

import (
	"industry-api/internal/models"
	"industry-api/internal/response"
	"industry-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BookingHandler handles HTTP requests related to hotel bookings.
type BookingHandler struct {
	svc *service.BookingService // Service layer for business logic
}

// NewBookingHandler creates and returns a new instance of BookingHandler.
// It accepts a BookingService dependency for handling booking operations.
func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

// AddBooking handles HTTP POST requests to create a new booking.
// It validates the incoming request, converts it to a Booking model, and calls the service layer.
// Returns a 201 Created status with the created booking or an error response.
func (h *BookingHandler) AddBooking(c *gin.Context) {
	// Parse and validate the JSON request body
	var req models.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Return 400 Bad Request if JSON binding fails
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}
	// Additional validation for required fields
	if req.UserID == 0 || req.RoomID == 0 || req.CheckInDate.IsZero() || req.CheckOutDate.IsZero() || req.Adults == 0 || req.TotalAmount == 0 {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
	}

	// Convert request to domain model
	booking := &models.Booking{
		UserID:          req.UserID,
		RoomID:          req.RoomID,
		CheckInDate:     req.CheckInDate,
		CheckOutDate:    req.CheckOutDate,
		Adults:          req.Adults,
		Children:        req.Children,
		SpecialRequests: req.SpecialRequests,
		Status:          req.Status,
		PaymentStatus:   req.PaymentStatus,
		TotalAmount:     req.TotalAmount,
	}

	// Call service to create the booking
	createdBooking, err := h.svc.AddBooking(c, booking)
	if err != nil {
		// Return 500 Internal Server Error if service call fails
		response.JSON(c, http.StatusInternalServerError, false, "failed to add booking", nil, err.Error())
		return
	}
	// Return 201 Created with the newly created booking
	response.JSON(c, http.StatusCreated, true, "booking added successfully", createdBooking, "")

}
