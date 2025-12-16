// Package repository provides database access layer implementations.
// Repositories handle all direct database operations using SQL queries.
package repository

import (
	"context"
	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RoomMaintenanceRepository provides database access for room maintenance operations.
type RoomMaintenanceRepository struct {
	db *pgxpool.Pool // Database connection pool
}

// RoomMaintenanceRepo defines methods used by services for room maintenance.
type RoomMaintenanceRepo interface {
	AddRoomMaintenance(ctx context.Context, roomMaintenance *models.RoomMaintenance) (*models.RoomMaintenance, error)
}

// NewRoomMaintenanceRepository creates and returns a new instance of RoomMaintenanceRepository.
// It accepts a database connection pool for executing database operations.
func NewRoomMaintenanceRepository(db *pgxpool.Pool) *RoomMaintenanceRepository {
	return &RoomMaintenanceRepository{db: db}
}

// AddRoomMaintenance inserts a new room maintenance record into the database.
// It schedules maintenance for a room with start date, end date, reason, and status.
// Returns the created maintenance record with ID and creation timestamp, or an error if the operation fails.
func (r *RoomMaintenanceRepository) AddRoomMaintenance(ctx context.Context, roomMaintenance *models.RoomMaintenance) (*models.RoomMaintenance, error) {
	query := `
    INSERT INTO room_maintenance (room_id, start_date, end_date, reason, status, created_by) 
    VALUES ($1, $2, $3, $4, $5, $6) 
    RETURNING id, created_at
    `

	// Execute the insert query and scan the returned ID and timestamp
	err := r.db.QueryRow(ctx, query,
		roomMaintenance.RoomID,
		roomMaintenance.StartDate,
		roomMaintenance.EndDate,
		roomMaintenance.Reason,
		roomMaintenance.Status,
		roomMaintenance.CreatedBy,
	).Scan(&roomMaintenance.ID, &roomMaintenance.CreatedAt)

	if err != nil {
		return nil, err
	}

	return roomMaintenance, nil
}
