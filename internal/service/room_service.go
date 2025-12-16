// Package service provides business logic layer implementations.
// Services contain validation logic and business rules before delegating to repositories.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"industry-api/internal/cache"
	"industry-api/internal/models"
	"industry-api/internal/repository"
	"time"
)

// RoomService handles all room-related business logic operations.
// It includes validation, caching logic, and delegates data access to the repository.
type RoomService struct {
	repo repository.RoomRepo // Repository interface for data access (allows mocking in tests)
}

// NewRoomService creates and returns a new instance of RoomService.
// It accepts a RoomRepository dependency for data access operations.
func NewRoomService(repo repository.RoomRepo) *RoomService {
	return &RoomService{repo: repo}
}

// AddRoom creates a new room after comprehensive validation.
// It validates all required room fields before delegating to the repository for persistence.
// Returns the created room or an error if validation fails.
func (s *RoomService) AddRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	fmt.Printf("Service: Adding room - %+v\n", room)

	// Validate room number is provided
	if room.RoomNumber == "" {
		return nil, errors.New("room number is required")
	}
	// Validate room type is provided
	if room.RoomType == "" {
		return nil, errors.New("room type is required")
	}
	// Validate room description is provided
	if room.Description == "" {
		return nil, errors.New("description is required")
	}
	// Validate price is greater than zero
	if room.Price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	// Validate room capacity is greater than zero
	if room.Capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}
	// Validate floor number is greater than zero
	if room.Floor <= 0 {
		return nil, errors.New("floor must be greater than 0")
	}
	// Validate amenities list is not empty
	if len(room.Amenities) == 0 { // Fix: Check slice length, not empty string
		return nil, errors.New("amenities are required")
	}

	fmt.Println("Service: Validation passed, calling repository...")

	// Delegate to repository to persist the room
	if err := s.repo.AddRoom(ctx, room); err != nil {
		fmt.Printf("Service: Repository error - %v\n", err)
		return nil, fmt.Errorf("failed to add room: %w", err)
	}

	fmt.Printf("Service: Room added successfully - ID: %d\n", room.ID)
	return room, nil
}

// GetRoomsList retrieves all rooms in the system without filtering or caching.
// Returns a list of all rooms or an error if the database operation fails.
func (s *RoomService) GetRoomsList(ctx context.Context) ([]*models.Room, error) {
	// Delegate to repository to fetch all rooms
	rooms, err := s.repo.GetRoomsList(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetAvailableRooms retrieves all available rooms with Redis caching.
// It first tries to fetch from cache, and if not found (cache miss),
// fetches from the database and caches the result for 10 minutes.
func (s *RoomService) GetAvailableRooms(ctx context.Context) ([]*models.Room, error) {
	// Cache key for storing available rooms
	cacheKey := "available_rooms"

	// Check if Redis client is initialized
	if cache.Client == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	// Try to get available rooms from cache
	cached, err := cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var rooms []*models.Room
		// Deserialize cached JSON into rooms slice
		if err := json.Unmarshal([]byte(cached), &rooms); err == nil {
			fmt.Println("\u2705 Cache hit - available rooms")
			return rooms, nil
		}
	}

	// If cache miss or error, proceed to database
	fmt.Println("\u274c Cache miss - fetching available rooms from database")

	// Fetch available rooms from database
	rooms, err := s.repo.GetAvailableRooms(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result for future requests (10-minute TTL)
	roomsJSON, err := json.Marshal(rooms)
	if err == nil {
		err = cache.Client.Set(ctx, cacheKey, roomsJSON, 10*time.Minute).Err()
		if err != nil {
			fmt.Printf("\u26a0\ufe0f Failed to cache rooms: %v\n", err)
		} else {
			fmt.Println("\u2705 Available rooms cached successfully")
		}
	}

	return rooms, nil
}
