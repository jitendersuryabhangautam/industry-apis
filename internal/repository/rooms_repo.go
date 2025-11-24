package repository

import (
	"context"
	"fmt"
	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomRepository struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) AddRoom(ctx context.Context, room *models.Room) error {
	query := `
    INSERT INTO rooms 
        (room_number, room_type, description, price_per_night, capacity, floor, amenities, is_available)
    VALUES 
        ($1, $2, $3, $4, $5, $6, $7, TRUE)
    RETURNING id, created_at, updated_at;
    `

	fmt.Printf("Repository: Executing query with amenities: %v\n", room.Amenities)

	// No conversion needed - room.Amenities is already []string
	err := r.db.QueryRow(
		ctx,
		query,
		room.RoomNumber,
		room.RoomType,
		room.Description,
		room.Price,
		room.Capacity,
		room.Floor,
		room.Amenities, // Pass the slice directly
	).Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		fmt.Printf("Repository: Database error - %v\n", err)
		return err
	}

	fmt.Printf("Repository: Room inserted successfully - ID: %d\n", room.ID)
	return nil
}

func (r *RoomRepository) GetRoomsList(ctx context.Context) ([]*models.Room, error) {
	query := `
    SELECT id, room_number, room_type, description, price_per_night, capacity, floor, amenities, is_available, created_at, updated_at
    FROM rooms
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roomsList []*models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomNumber,
			&room.RoomType,
			&room.Description,
			&room.Price,
			&room.Capacity,
			&room.Floor,
			&room.Amenities, // This should now work with []string
			&room.IsAvailable,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roomsList = append(roomsList, &room)
	}

	return roomsList, nil
}

func (r *RoomRepository) GetAvailableRooms(ctx context.Context) ([]*models.Room, error) {
	query := `
	SELECT id, room_number, room_type, description, price_per_night, capacity, floor, amenities, is_available
	FROM rooms
	WHERE is_available = true
	`
	rooms, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rooms.Close()

	var roomsList []*models.Room
	for rooms.Next() {
		var room models.Room
		err := rooms.Scan(&room.ID, &room.RoomNumber, &room.RoomType, &room.Description, &room.Price, &room.Capacity, &room.Floor, &room.Amenities, &room.IsAvailable)
		if err != nil {
			return nil, err
		}
		roomsList = append(roomsList, &room)
	}
	return roomsList, nil
}
