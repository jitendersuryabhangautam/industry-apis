package service

import (
	"context"
	"errors"
	"industry-api/internal/models"
	"industry-api/internal/repository"
)

type RoomMaintenanceService struct {
	repo *repository.RoomMaintenanceRepository
}

func NewRoomMaintenanceService(repo *repository.RoomMaintenanceRepository) *RoomMaintenanceService {
	return &RoomMaintenanceService{repo: repo}
}

func (s *RoomMaintenanceService) AddRoomMaintenance(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error) {
	if rm.RoomID == 0 {
		return nil, errors.New("room id is required")
	}
	if rm.StartDate.IsZero() {
		return nil, errors.New("start date is required")
	}
	if rm.EndDate.IsZero() {
		return nil, errors.New("end date is required")
	}
	if rm.Reason == "" {
		return nil, errors.New("reason is required")
	}
	if rm.Status == "" {
		return nil, errors.New("status is required")
	}
	if rm.CreatedBy == 0 {
		return nil, errors.New("created by is required")
	}
	room, err := s.repo.AddRoomMaintenance(ctx, rm)
	if err != nil {
		return nil, err
	}
	return room, nil
}
