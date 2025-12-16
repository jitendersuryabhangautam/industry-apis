// Package handler contains HTTP request handlers for all API endpoints.
// Handlers receive HTTP requests, validate input, call services, and send responses.
package handler

import (
	"bytes"
	"fmt"
	"industry-api/internal/models"
	"industry-api/internal/response"
	"industry-api/internal/service"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests related to user management and authentication.
type UserHandler struct {
	svc *service.UserService // Service layer for user business logic
}

// NewUserHandler creates and returns a new instance of UserHandler.
// It accepts a UserService dependency for handling user operations.
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register handles HTTP POST requests for user registration.
// It validates the registration request, creates a new user account, and returns the created user.
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	// Parse and validate the JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "Invalid request", nil, err.Error())
		return
	}
	// Additional validation for required fields
	if req.Name == "" || req.Email == "" || req.Password == "" {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
		return
	}

	// Convert request to domain model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Role:     req.Role,
	}
	// Call service to create the user account
	createdUser, err := h.svc.CreateUser(c.Request.Context(), user)
	if err != nil {
		// Return 500 Internal Server Error if service call fails
		response.JSON(c, http.StatusInternalServerError, false, "failed to create user", nil, err.Error())
		return
	}
	// Return 201 Created with the newly created user
	response.JSON(c, http.StatusCreated, true, "user created successfully", createdUser, "")

}

// LoginUser handles HTTP POST requests for user authentication.
// It validates credentials and returns a JWT access token upon successful authentication.
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req models.LoginRequest
	// Parse and validate the JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	// Call service to authenticate user and generate token
	user, err := h.svc.LoginUser(c, req.Email, req.Password)
	if err != nil {
		// Return 401 Unauthorized if authentication fails
		response.JSON(c, http.StatusUnauthorized, false, "Login failed", nil, err.Error())
		return
	}
	// Return 200 OK with user details and JWT token
	response.JSON(c, http.StatusOK, true, "Login successful", user, "")

}

// GetUsersQuery represents the query parameters for listing users with filtering and pagination.
type GetUsersQuery struct {
	Role     string `form:"role"`                                     // Filter by user role
	IsActive *bool  `form:"is_active"`                                // Filter by active status (pointer to distinguish false from not provided)
	Search   string `form:"search"`                                   // Search in name or email
	Limit    int    `form:"limit" binding:"omitempty,min=0,max=1000"` // Records per page (0 means all records)
	Page     int    `form:"page" binding:"omitempty,min=1"`           // Page number for pagination
}

// GetUserList handles HTTP GET requests to retrieve a paginated list of users with optional filtering.
// It supports filtering by role, active status, and search term (name/email).
// Supports pagination with page and limit query parameters.
func (h *UserHandler) GetUserList(c *gin.Context) {
	var req GetUsersQuery
	// Parse and validate query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	// Set default values for pagination
	if req.Limit == 0 {
		req.Limit = 0 // 0 means return all records
	}
	if req.Page == 0 {
		req.Page = 1 // Default to page 1
	}

	// Call service to get filtered user list
	users, err := h.svc.GetUserList(c, req.Role, req.IsActive, req.Search, req.Page, req.Limit)
	if err != nil {
		// Return 500 Internal Server Error if service call fails
		response.JSON(c, http.StatusInternalServerError, false, "failed to get user list", nil, err.Error())
		return
	}
	// Return 200 OK with the paginated user list
	response.JSON(c, http.StatusOK, true, "user list retrieved successfully", users, "")
}

// GetUserByID handles HTTP GET requests to retrieve a specific user by their ID.
// Returns 404 if the user is not found, or 500 if a database error occurs.
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Extract user ID from URL parameter
	id := c.Param("id")
	if id == "" {
		response.JSON(c, http.StatusBadRequest, false, "missing ID", nil, "ID is required")
		return
	}
	// Convert string ID to integer
	userID, err := strconv.Atoi(id)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid ID", nil, "ID must be a valid integer")
		return
	}
	// Call service to retrieve user by ID
	user, err := h.svc.GetUserByID(c, userID)
	if err != nil {
		// Return 404 Not Found if user doesn't exist
		if err.Error() == "user not found" {
			response.JSON(c, http.StatusNotFound, false, "user not found", nil, err.Error())
			return
		}
		// Return 500 Internal Server Error for other errors
		response.JSON(c, http.StatusInternalServerError, false, "failed to get user", nil, err.Error())
		return
	}
	// Return 200 OK with the user details
	response.JSON(c, http.StatusOK, true, "user retrieved successfully", user, "")

}

// UpdateUserStatus handles HTTP PUT requests to update a user's active status.
// It accepts the user ID as a URL parameter and the new status in the request body.
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	// Extract user ID from URL parameter
	id := c.Param("id")
	if id == "" {
		response.JSON(c, http.StatusBadRequest, false, "missing ID", nil, "ID is required")
		return
	}

	// Debug: Read raw request body
	body, _ := c.GetRawData()
	fmt.Printf("DEBUG: Raw request body: %s\n", string(body))
	// Reset body for JSON binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Parse and validate the JSON request body
	var req models.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("DEBUG: Binding error: %v\n", err)
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	fmt.Printf("DEBUG: Parsed request - IsActive: %t\n", req.IsActive)

	// Convert string ID to integer
	userID, err := strconv.Atoi(id)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid ID", nil, "ID must be a valid integer")
		return
	}

	// Call service to update user status
	user, err := h.svc.UpdateUserStatus(c, userID, req.IsActive)
	if err != nil {
		// Return 404 Not Found if user doesn't exist
		if err.Error() == "user not found" {
			response.JSON(c, http.StatusNotFound, false, "user not found", nil, err.Error())
			return
		}
		// Return 500 Internal Server Error for other errors
		response.JSON(c, http.StatusInternalServerError, false, "failed to update user status", nil, err.Error())
		return
	}

	// Build response with updated status
	resp := models.UpdateUserStatusResponse{
		ID:        user.ID,
		IsActive:  user.IsActive,
		UpdatedAt: user.UpdatedAt,
	}

	// Return 200 OK with the updated status
	response.JSON(c, http.StatusOK, true, "user status updated successfully", resp, "")
}
