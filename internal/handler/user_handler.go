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

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "Invalid request", nil, err.Error())
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		response.JSON(c, http.StatusBadRequest, false, "all fields are required", nil, "missing required fields")
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Role:     req.Role,
	}
	createdUser, err := h.svc.CreateUser(c.Request.Context(), user)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "failed to create user", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, true, "user created successfully", createdUser, "")

}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	user, err := h.svc.LoginUser(c, req.Email, req.Password)
	if err != nil {
		response.JSON(c, http.StatusUnauthorized, false, "Login failed", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, true, "Login successful", user, "")

}

type GetUsersQuery struct {
	Role     string `form:"role"`
	IsActive *bool  `form:"is_active"`                                // Use pointer to distinguish between false and not provided
	Search   string `form:"search"`                                   // Search in name or email
	Limit    int    `form:"limit" binding:"omitempty,min=0,max=1000"` // 0 means all records
	Page     int    `form:"page" binding:"omitempty,min=1"`
}

func (h *UserHandler) GetUserList(c *gin.Context) {
	var req GetUsersQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	// Set default values
	if req.Limit == 0 {
		req.Limit = 0 // 0 means return all records
	}
	if req.Page == 0 {
		req.Page = 1 // Default to page 1
	}

	users, err := h.svc.GetUserList(c, req.Role, req.IsActive, req.Search, req.Page, req.Limit)
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, false, "failed to get user list", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, true, "user list retrieved successfully", users, "")
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.JSON(c, http.StatusBadRequest, false, "missing ID", nil, "ID is required")
		return
	}
	userID, err := strconv.Atoi(id)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid ID", nil, "ID must be a valid integer")
		return
	}
	user, err := h.svc.GetUserByID(c, userID)
	if err != nil {
		if err.Error() == "user not found" {
			response.JSON(c, http.StatusNotFound, false, "user not found", nil, err.Error())
			return
		}
		response.JSON(c, http.StatusInternalServerError, false, "failed to get user", nil, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, true, "user retrieved successfully", user, "")

}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.JSON(c, http.StatusBadRequest, false, "missing ID", nil, "ID is required")
		return
	}

	// Debug the raw request
	body, _ := c.GetRawData()
	fmt.Printf("DEBUG: Raw request body: %s\n", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req models.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("DEBUG: Binding error: %v\n", err)
		response.JSON(c, http.StatusBadRequest, false, "invalid request", nil, err.Error())
		return
	}

	fmt.Printf("DEBUG: Parsed request - IsActive: %t\n", req.IsActive)

	userID, err := strconv.Atoi(id)
	if err != nil {
		response.JSON(c, http.StatusBadRequest, false, "invalid ID", nil, "ID must be a valid integer")
		return
	}

	user, err := h.svc.UpdateUserStatus(c, userID, req.IsActive)
	if err != nil {
		if err.Error() == "user not found" {
			response.JSON(c, http.StatusNotFound, false, "user not found", nil, err.Error())
			return
		}
		response.JSON(c, http.StatusInternalServerError, false, "failed to update user status", nil, err.Error())
		return
	}

	resp := models.UpdateUserStatusResponse{
		ID:        user.ID, // ← Now converted to int
		IsActive:  user.IsActive,
		UpdatedAt: user.UpdatedAt,
	}

	response.JSON(c, http.StatusOK, true, "user status updated successfully", resp, "")
}
