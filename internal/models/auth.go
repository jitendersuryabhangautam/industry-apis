// Package models defines all data structures used throughout the application.
package models

import "time"

// LoginRequest represents the HTTP request body for user login.
// It contains the user's email and password for authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`    // User's email address (required, must be valid email format)
	Password string `json:"password" binding:"required,min=6,max=255"` // User's password (required, 6-255 characters)
}

// LoginResponse represents the HTTP response body after successful user authentication.
// It contains the authenticated user's information and a JWT access token.
type LoginResponse struct {
	User      *User         `json:"user"`         // The authenticated user's details
	Token     string        `json:"access_token"` // JWT token for authenticated requests
	TokenType string        `json:"token_type"`   // Type of token (typically "Bearer")
	ExpiresIn time.Duration `json:"expires_in"`   // Token expiration time duration
}
