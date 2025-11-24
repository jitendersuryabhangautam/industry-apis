package handler

import (
	"industry-api/internal/models"
	"industry-api/internal/response"
	"industry-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

func (h *BookingHandler) AddBooking(c *gin.Context) {
	var req models.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}
	if req.UserID == 0 || req.RoomID == 0 || req.CheckInDate.IsZero() || req.CheckOutDate.IsZero() || req.Adults == 0 || req.TotalAmount == 0 {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
	}

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

	createdBooking, err := h.svc.AddBooking(c, booking)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "failed to add booking", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, true, "booking added successfully", createdBooking, "")

}
