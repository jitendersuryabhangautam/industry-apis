package models

import "time"

type RoomMaintenance struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"` // Consistent naming
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoomMaintenanceRequest struct {
	RoomID    int       `json:"room_id" binding:"required"` // Changed to RoomID for consistency
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
	Status    string    `json:"status" binding:"required"` // Added required binding
	CreatedBy int       `json:"created_by" binding:"required"`
}
