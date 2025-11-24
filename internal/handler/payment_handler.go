package handler

import (
	"industry-api/internal/models"
	"industry-api/internal/response"
	"industry-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	if req.Amount == 0 || req.BookingID == 0 || req.PaymentMethod == "" || req.TransactionID == "" {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
		return
	}

	payment := &models.Payment{
		BookingID:     req.BookingID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		TransactionID: &req.TransactionID,
		Status:        "pending",
	}

	pmt, err := h.svc.InitiatePaymet(c, payment)

	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "failed to initiate payment", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, true, "payment initiated successfully", pmt, "")

}

func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	var req models.UpdatePaymentRequest
	idStr := c.Query("id")
	if idStr == "" {
		response.JSON(c, http.StatusBadRequest, false, "missing required fields", nil, "missing ID query parameter")
		return
	}

	// Convert the string ID to an integer ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, "ID must be a valid integer")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	if req.CardBrand == "" || req.CardLastFour == "" {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
		return
	}
	pmt := &models.Payment{
		CardBrand:    &req.CardBrand,
		CardLastFour: &req.CardLastFour,
		ReceitURL:    &req.ReceiptURL,
		Status:       "paid",
	}
	updatedPmt, err := h.svc.UpdatePayment(c, pmt, id)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "failed to update payment", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, true, "payment updated successfully", updatedPmt, "")
}
