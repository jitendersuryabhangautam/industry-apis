package handler

import (
	"industry-api/internal/models"
	"industry-api/internal/response"
	"industry-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomMaintenanceHandler struct {
	svc *service.RoomMaintenanceService
}

func NewRoomMaintenanceHandler(svc *service.RoomMaintenanceService) *RoomMaintenanceHandler {
	return &RoomMaintenanceHandler{svc: svc}
}

func (h *RoomMaintenanceHandler) AddRoomMaintenance(c *gin.Context) {
	var req models.RoomMaintenanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "Invalid request", nil, err.Error())
		return
	}
	if req.RoomID == 0 || req.StartDate.IsZero() || req.EndDate.IsZero() || req.Reason == "" || req.Status == "" || req.CreatedBy == 0 {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
		return
	}

	roomMaintenance := &models.RoomMaintenance{
		RoomID:    req.RoomID, // Fix field name inconsistency
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Reason:    req.Reason,
		Status:    req.Status,
		CreatedBy: req.CreatedBy,
	}
	createdRoomMaintenance, err := h.svc.AddRoomMaintenance(c, roomMaintenance)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "failed to add room maintenance", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, true, "room maintenance added successfully", createdRoomMaintenance, "")
}
