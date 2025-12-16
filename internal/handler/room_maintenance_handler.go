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

// RoomMaintenanceHandler handles HTTP requests related to room maintenance management.
type RoomMaintenanceHandler struct {
	svc *service.RoomMaintenanceService // Service layer for room maintenance business logic
}

// NewRoomMaintenanceHandler creates and returns a new instance of RoomMaintenanceHandler.
// It accepts a RoomMaintenanceService dependency for handling maintenance operations.
func NewRoomMaintenanceHandler(svc *service.RoomMaintenanceService) *RoomMaintenanceHandler {
	return &RoomMaintenanceHandler{svc: svc}
}

// AddRoomMaintenance handles HTTP POST requests to schedule room maintenance.
// It validates the incoming request, converts it to a RoomMaintenance model, and calls the service layer.
// Returns a 201 Created status with the created maintenance record or an error response.
func (h *RoomMaintenanceHandler) AddRoomMaintenance(c *gin.Context) {
	// Parse and validate the JSON request body
	var req models.RoomMaintenanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// Return 400 Bad Request if JSON binding fails
		response.JSON(c, http.StatusBadRequest, false, "Invalid request", nil, err.Error())
		return
	}
	// Additional validation for required fields
	if req.RoomID == 0 || req.StartDate.IsZero() || req.EndDate.IsZero() || req.Reason == "" || req.Status == "" || req.CreatedBy == 0 {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
		return
	}

	// Convert request to domain model
	roomMaintenance := &models.RoomMaintenance{
		RoomID:    req.RoomID,    // Room ID for the maintenance
		StartDate: req.StartDate, // When maintenance starts
		EndDate:   req.EndDate,   // When maintenance ends
		Reason:    req.Reason,    // Reason for maintenance
		Status:    req.Status,    // Current maintenance status
		CreatedBy: req.CreatedBy, // User who created the maintenance record
	}
	// Call service to create the maintenance record
	createdRoomMaintenance, err := h.svc.AddRoomMaintenance(c, roomMaintenance)
	if err != nil {
		// Return 500 Internal Server Error if service call fails
		response.JSON(c, http.StatusInternalServerError, false, "failed to add room maintenance", nil, err.Error())
		return
	}
	// Return 201 Created with the newly created maintenance record
	response.JSON(c, http.StatusCreated, true, "room maintenance added successfully", createdRoomMaintenance, "")
}
