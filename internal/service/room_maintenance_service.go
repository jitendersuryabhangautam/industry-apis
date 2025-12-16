// Package service provides business logic layer implementations.
// Services contain validation logic and business rules before delegating to repositories.
package service

import (
	"context"
	"errors"
	"industry-api/internal/models"
	"industry-api/internal/repository"
)

// RoomMaintenanceService handles all room maintenance-related business logic operations.
// It acts as an intermediary between handlers and repository layers.
type RoomMaintenanceService struct {
	repo repository.RoomMaintenanceRepo // Repository interface for testing
}

// NewRoomMaintenanceService creates and returns a new instance of RoomMaintenanceService.
// It accepts a RoomMaintenanceRepo interface for data access operations.
func NewRoomMaintenanceService(repo repository.RoomMaintenanceRepo) *RoomMaintenanceService {
	return &RoomMaintenanceService{repo: repo}
}

// AddRoomMaintenance schedules room maintenance after validating all required fields.
// It validates all maintenance data before delegating to the repository for persistence.
// Returns the created maintenance record or an error if validation fails.
func (s *RoomMaintenanceService) AddRoomMaintenance(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error) {
	// Validate room ID is provided
	if rm.RoomID == 0 {
		return nil, errors.New("room id is required")
	}
	// Validate maintenance start date is provided
	if rm.StartDate.IsZero() {
		return nil, errors.New("start date is required")
	}
	// Validate maintenance end date is provided
	if rm.EndDate.IsZero() {
		return nil, errors.New("end date is required")
	}
	// Validate maintenance reason is provided
	if rm.Reason == "" {
		return nil, errors.New("reason is required")
	}
	// Validate maintenance status is provided
	if rm.Status == "" {
		return nil, errors.New("status is required")
	}
	// Validate user ID who created the record is provided
	if rm.CreatedBy == 0 {
		return nil, errors.New("created by is required")
	}
	// Delegate to repository to persist the maintenance record
	room, err := s.repo.AddRoomMaintenance(ctx, rm)
	if err != nil {
		return nil, err
	}
	return room, nil
}
