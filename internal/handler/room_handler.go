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

// RoomHandler handles HTTP requests related to hotel room management.
type RoomHandler struct {
	svc *service.RoomService // Service layer for room business logic
}

// NewRoomHandler creates and returns a new instance of RoomHandler.
// It accepts a RoomService dependency for handling room operations.
func NewRoomHandler(svc *service.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

// AddRoom handles HTTP POST requests to create a new room.
// It validates the incoming request, converts it to a Room model, and calls the service layer.
// Returns a 201 Created status with the created room or an error response.
func (h *RoomHandler) AddRoom(c *gin.Context) {
	// Parse and validate the JSON request body
	var req models.RoomRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// Return 400 Bad Request if JSON binding fails
		response.JSON(c, http.StatusBadRequest, false, "Invalid request", nil, err.Error())
		return
	}

	// Convert request to domain model
	room := &models.Room{
		RoomNumber:  req.RoomNumber,
		RoomType:    req.RoomType,
		Description: req.Description,
		Price:       req.Price,
		Capacity:    req.Capacity,
		Floor:       req.Floor,
		Amenities:   req.Amenities,
	}

	// Call service to create the room
	createdRoom, err := h.svc.AddRoom(c, room)
	if err != nil {
		// Return 500 Internal Server Error if service call fails
		response.JSON(c, http.StatusInternalServerError, false, "failed to add room", nil, err.Error())
		return
	}

	// Return 201 Created with the newly created room
	response.JSON(c, http.StatusCreated, true, "room added successfully", createdRoom, "")

}

// GetRoomsList handles HTTP GET requests to retrieve all rooms in the system.
// It fetches the list of all rooms from the service layer and returns them.
func (h *RoomHandler) GetRoomsList(c *gin.Context) {
	// Call service to get all rooms
	rooms, err := h.svc.GetRoomsList(c)
	if err != nil {
		// Return 400 Bad Request if service call fails
		response.JSON(c, http.StatusBadRequest, false, "failed to get rooms", nil, err.Error())
		return
	}
	// Return 200 OK with the list of rooms
	response.JSON(c, http.StatusOK, true, "rooms fetched successfully", rooms, "")
}

// GetAvailableRooms handles HTTP GET requests to retrieve all available rooms.
// It fetches rooms that are available for booking (with caching support).
func (h *RoomHandler) GetAvailableRooms(c *gin.Context) {
	// Call service to get available rooms (may be cached)
	rooms, err := h.svc.GetAvailableRooms(c)
	if err != nil {
		// Return 400 Bad Request if service call fails
		response.JSON(c, http.StatusBadRequest, false, "failed to get available rooms", nil, err.Error())
		return
	}
	// Return 200 OK with the list of available rooms
	response.JSON(c, http.StatusOK, true, "available rooms fetched successfully", rooms, "")
}
