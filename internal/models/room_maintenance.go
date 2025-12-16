// Package models defines all data structures used throughout the application.
package models

import "time"

// RoomMaintenance represents a room maintenance record in the system.
type RoomMaintenance struct {
	ID        int       `json:"id"`         // Unique maintenance record identifier
	RoomID    int       `json:"room_id"`    // ID of the room undergoing maintenance
	StartDate time.Time `json:"start_date"` // Date and time when maintenance starts
	EndDate   time.Time `json:"end_date"`   // Date and time when maintenance ends
	Reason    string    `json:"reason"`     // Reason for maintenance
	Status    string    `json:"status"`     // Current status of maintenance (e.g., "scheduled", "in-progress", "completed")
	CreatedBy int       `json:"created_by"` // ID of the user who created the maintenance record
	CreatedAt time.Time `json:"created_at"` // Timestamp when the record was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp of the last update
}

// RoomMaintenanceRequest represents the HTTP request body for creating a room maintenance record.
type RoomMaintenanceRequest struct {
	RoomID    int       `json:"room_id" binding:"required"`    // Room ID
	StartDate time.Time `json:"start_date" binding:"required"` // Maintenance start date
	EndDate   time.Time `json:"end_date" binding:"required"`   // Maintenance end date
	Reason    string    `json:"reason" binding:"required"`     // Reason for maintenance
	Status    string    `json:"status" binding:"required"`     // Maintenance status
	CreatedBy int       `json:"created_by" binding:"required"` // ID of user creating the record
}
