package models

import "time"

// Room maps directly to your table
type Room struct {
	ID          int       `json:"id"`
	RoomNumber  string    `json:"room_number"`
	RoomType    string    `json:"room_type"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Capacity    int       `json:"capacity"`
	Floor       int       `json:"floor"`
	Amenities   []string  `json:"amenities"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// For incoming POST /room requests
type RoomRequest struct {
	RoomNumber  string   `json:"room_number" binding:"required"`
	RoomType    string   `json:"room_type" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       float64  `json:"price" binding:"required"` // ✅ changed from int → float64
	Capacity    int      `json:"capacity" binding:"required"`
	Floor       int      `json:"floor" binding:"required"`
	Amenities   []string `json:"amenities" binding:"required"`
}
