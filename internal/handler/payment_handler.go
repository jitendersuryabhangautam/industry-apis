// Package handler contains HTTP request handlers for all API endpoints.
// Handlers receive HTTP requests, validate input, call services, and send responses.
package handler

import (
	"industry-api/internal/models"
	"industry-api/internal/response"
	"industry-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaymentHandler handles HTTP requests related to payment processing.
type PaymentHandler struct {
	svc *service.PaymentService // Service layer for payment business logic
}

// NewPaymentHandler creates and returns a new instance of PaymentHandler.
// It accepts a PaymentService dependency for handling payment operations.
func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// InitiatePayment handles HTTP POST requests to initiate a new payment.
// It validates the incoming payment request and initiates payment processing.
// Returns a 201 Created status with the initiated payment or an error response.
func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	// Parse and validate the JSON request body
	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Return 400 Bad Request if JSON binding fails
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	// Additional validation for required fields
	if req.Amount == 0 || req.BookingID == 0 || req.PaymentMethod == "" || req.TransactionID == "" {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
		return
	}

	// Convert request to domain model with initial status
	payment := &models.Payment{
		BookingID:     req.BookingID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		TransactionID: &req.TransactionID,
		Status:        "pending", // Initial status is pending
	}

	// Call service to initiate payment
	pmt, err := h.svc.InitiatePaymet(c, payment)

	if err != nil {
		// Return 500 Internal Server Error if service call fails
		response.JSON(c, http.StatusInternalServerError, false, "failed to initiate payment", nil, err.Error())
		return
	}
	// Return 201 Created with the initiated payment
	response.JSON(c, http.StatusCreated, true, "payment initiated successfully", pmt, "")

}

// UpdatePayment handles HTTP PUT requests to update an existing payment with card details.
// It accepts a payment ID as a query parameter and updates payment status and card information.
// Returns a 200 OK status with the updated payment or an error response.
func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	// Parse request body for payment details
	var req models.UpdatePaymentRequest
	// Get payment ID from query parameters
	idStr := c.Query("id")
	if idStr == "" {
		response.JSON(c, http.StatusBadRequest, false, "missing required fields", nil, "missing ID query parameter")
		return
	}

	// Convert the string ID to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, "ID must be a valid integer")
		return
	}
	// Parse and validate the JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	// Validate required fields
	if req.CardBrand == "" || req.CardLastFour == "" {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
		return
	}
	// Create payment model with updated information
	pmt := &models.Payment{
		CardBrand:    &req.CardBrand,
		CardLastFour: &req.CardLastFour,
		ReceitURL:    &req.ReceiptURL,
		Status:       "paid", // Mark payment as paid
	}
	// Call service to update payment
	updatedPmt, err := h.svc.UpdatePayment(c, pmt, id)
	if err != nil {
		// Return 500 Internal Server Error if service call fails
		response.JSON(c, http.StatusInternalServerError, false, "failed to update payment", nil, err.Error())
		return
	}
	// Return 200 OK with the updated payment
	response.JSON(c, http.StatusOK, true, "payment updated successfully", updatedPmt, "")
}
