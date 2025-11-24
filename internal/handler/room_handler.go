package handler

import (
	"industry-api/internal/models"
	"industry-api/internal/response"
	"industry-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	svc *service.RoomService
}

func NewRoomHandler(svc *service.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

func (h *RoomHandler) AddRoom(c *gin.Context) {
	var req models.RoomRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "Invalid request", nil, err.Error())
		return
	}

	room := &models.Room{
		RoomNumber:  req.RoomNumber,
		RoomType:    req.RoomType,
		Description: req.Description,
		Price:       req.Price,
		Capacity:    req.Capacity,
		Floor:       req.Floor,
		Amenities:   req.Amenities,
	}

	createdRoom, err := h.svc.AddRoom(c, room)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "failed to add room", nil, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, true, "room added successfully", createdRoom, "")

}

func (h *RoomHandler) GetRoomsList(c *gin.Context) {
	rooms, err := h.svc.GetRoomsList(c)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "failed to get rooms", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, true, "rooms fetched successfully", rooms, "")
}

func (h *RoomHandler) GetAvailableRooms(c *gin.Context) {
	rooms, err := h.svc.GetAvailableRooms(c)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "failed to get available rooms", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, true, "available rooms fetched successfully", rooms, "")
}
