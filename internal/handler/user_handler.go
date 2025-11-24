package handler

import (
	"industry-api/internal/models"
	"industry-api/internal/response"
	"industry-api/internal/service"
	"net/http"

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
