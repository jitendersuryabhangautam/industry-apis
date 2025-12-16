// Package response provides utility functions for sending JSON responses to HTTP clients.
// It standardizes the response format across all API endpoints.
package response

import "github.com/gin-gonic/gin"

// JSON sends a standardized JSON response to the HTTP client.
// It wraps the response in a consistent format containing success status, message, data, and error information.
// This ensures all API endpoints return data in a uniform structure.
//
// Parameters:
// - c: The Gin context containing the HTTP response writer
// - status: HTTP status code (e.g., 200 for OK, 400 for Bad Request, 500 for Internal Server Error)
// - success: Boolean indicating if the operation was successful (true) or failed (false)
// - message: Human-readable message describing the response (e.g., "User created successfully")
// - data: The actual response data payload (can be nil if not applicable)
// - err: Error message string (empty string "" if no error occurred)
func JSON(c *gin.Context, status int, success bool, message string, data interface{}, err string) {
	// Send JSON response with standardized format
	c.JSON(status, gin.H{
		"success": success, // Indicates operation success/failure
		"message": message, // Human-readable message
		"data":    data,    // Response data payload
		"error":   err,     // Error message if any
	})
}
