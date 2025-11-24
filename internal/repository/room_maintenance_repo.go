package repository

import (
	"context"
	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomMaintenanceRepository struct {
	db *pgxpool.Pool
}

func NewRoomMaintenanceRepository(db *pgxpool.Pool) *RoomMaintenanceRepository {
	return &RoomMaintenanceRepository{db: db}
}

func (r *RoomMaintenanceRepository) AddRoomMaintenance(ctx context.Context, roomMaintenance *models.RoomMaintenance) (*models.RoomMaintenance, error) {
	query := `
    INSERT INTO room_maintenance (room_id, start_date, end_date, reason, status, created_by) 
    VALUES ($1, $2, $3, $4, $5, $6) 
    RETURNING id, created_at
    `

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
