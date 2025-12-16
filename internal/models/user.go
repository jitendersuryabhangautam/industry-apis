// Package models defines all data structures used throughout the application.
package models

import "time"

// User represents a user account in the system.
type User struct {
	ID        string    `json:"id"`         // Unique user identifier
	Name      string    `json:"name"`       // User's full name
	Email     string    `json:"email"`      // User's email address (unique)
	Phone     string    `json:"phone"`      // User's phone number
	Password  string    `json:"-"`          // Password hash (NOT sent in API responses for security)
	Role      string    `json:"role"`       // User role (e.g., "admin", "guest", "staff")
	IsActive  bool      `json:"is_active"`  // Whether the user account is active
	CreatedAt time.Time `json:"created_at"` // Account creation timestamp
	UpdatedAt time.Time `json:"updated_at"` // Last update timestamp
}

// RegisterRequest represents the HTTP request body for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,max=100"`           // User's full name
	Email    string `json:"email" binding:"required,email,max=255"`    // User's email
	Password string `json:"password" binding:"required,min=6,max=255"` // User's password (min 6 chars)
	Phone    string `json:"phone" binding:"required,len=10"`           // User's phone (exactly 10 digits)
	Role     string `json:"role"`                                      // User's role
}

// UserListResponse represents the HTTP response for listing users with pagination.
type UserListResponse struct {
	Page       int    `json:"page"`        // Current page number
	Limit      int    `json:"limit"`       // Items per page
	Total      int    `json:"total"`       // Total number of users
	TotalPages int    `json:"total_pages"` // Total number of pages
	Users      []User `json:"users"`       // List of users on this page
}

// UpdateUserStatusRequest represents the HTTP request body for updating user status.
type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active" binding:"required"` // New active status
}

// UpdateUserStatusResponse represents the HTTP response when updating user status.
type UpdateUserStatusResponse struct {
	ID        string    `json:"id"`         // User ID
	IsActive  bool      `json:"is_active"`  // Updated active status
	UpdatedAt time.Time `json:"updated_at"` // Update timestamp
}
