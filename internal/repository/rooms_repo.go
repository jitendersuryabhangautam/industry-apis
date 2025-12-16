// Package repository provides database access layer implementations.
// Repositories handle all direct database operations using SQL queries.
package repository

import (
	"context"
	"fmt"
	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RoomRepository provides database access for room operations.
type RoomRepository struct {
	db *pgxpool.Pool // Database connection pool
}

// RoomRepo defines the methods used by services for room data access.
type RoomRepo interface {
	AddRoom(ctx context.Context, room *models.Room) error
	GetRoomsList(ctx context.Context) ([]*models.Room, error)
	GetAvailableRooms(ctx context.Context) ([]*models.Room, error)
}

// NewRoomRepository creates and returns a new instance of RoomRepository.
// It accepts a database connection pool for executing database operations.
func NewRoomRepository(db *pgxpool.Pool) *RoomRepository {
	return &RoomRepository{db: db}
}

// AddRoom inserts a new room record into the database.
// It executes an INSERT query with room details and returns the generated room ID and timestamps.
// The room is created with is_available set to TRUE by default.
// Returns nil on success or an error if the database operation fails.
func (r *RoomRepository) AddRoom(ctx context.Context, room *models.Room) error {
	// SQL query to insert a new room record
	query := `
    INSERT INTO rooms 
        (room_number, room_type, description, price_per_night, capacity, floor, amenities, is_available)
    VALUES 
        ($1, $2, $3, $4, $5, $6, $7, TRUE)
    RETURNING id, created_at, updated_at;
    `

	fmt.Printf("Repository: Executing query with amenities: %v\n", room.Amenities)

	// Execute the insert query and scan the returned ID and timestamps
	// No conversion needed - room.Amenities is already []string
	err := r.db.QueryRow(
		ctx,
		query,
		room.RoomNumber,  // Room number identifier
		room.RoomType,    // Type of room
		room.Description, // Room description
		room.Price,       // Price per night
		room.Capacity,    // Guest capacity
		room.Floor,       // Floor number
		room.Amenities,   // List of amenities
	).Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		fmt.Printf("Repository: Database error - %v\n", err)
		return err
	}

	fmt.Printf("Repository: Room inserted successfully - ID: %d\n", room.ID)
	return nil
}

// GetRoomsList retrieves all rooms from the database.
// It returns a list of all rooms with their complete details.
// Returns a slice of room pointers or an error if the database query fails.
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

// GetAvailableRooms retrieves all rooms that are currently available for booking.
// It filters for rooms where is_available is true.
// Returns a slice of available room pointers or an error if the database query fails.
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
