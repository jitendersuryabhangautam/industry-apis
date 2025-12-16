// Package models defines all data structures used throughout the application.
package models

import "time"

// Room represents a hotel room record stored in the database.
type Room struct {
	ID          int       `json:"id"`           // Unique room identifier
	RoomNumber  string    `json:"room_number"`  // Room number or identifier (e.g., "101", "Suite-A")
	RoomType    string    `json:"room_type"`    // Type of room (e.g., "single", "double", "suite")
	Description string    `json:"description"`  // Detailed description of the room
	Price       float64   `json:"price"`        // Price per night
	Capacity    int       `json:"capacity"`     // Maximum number of guests the room can accommodate
	Floor       int       `json:"floor"`        // Floor number where the room is located
	Amenities   []string  `json:"amenities"`    // List of amenities available in the room
	IsAvailable bool      `json:"is_available"` // Whether the room is currently available for booking
	CreatedAt   time.Time `json:"created_at"`   // Timestamp when the room record was created
	UpdatedAt   time.Time `json:"updated_at"`   // Timestamp of the last update
}

// RoomRequest represents the HTTP request body for creating or updating a room.
type RoomRequest struct {
	RoomNumber  string   `json:"room_number" binding:"required"` // Room number
	RoomType    string   `json:"room_type" binding:"required"`   // Type of room
	Description string   `json:"description" binding:"required"` // Room description
	Price       float64  `json:"price" binding:"required"`       // Price per night
	Capacity    int      `json:"capacity" binding:"required"`    // Guest capacity
	Floor       int      `json:"floor" binding:"required"`       // Floor number
	Amenities   []string `json:"amenities" binding:"required"`   // List of amenities
}
