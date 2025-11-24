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

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) AddRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	fmt.Printf("Service: Adding room - %+v\n", room)

	// Validation
	if room.RoomNumber == "" {
		return nil, errors.New("room number is required")
	}
	if room.RoomType == "" {
		return nil, errors.New("room type is required")
	}
	if room.Description == "" {
		return nil, errors.New("description is required")
	}
	if room.Price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	if room.Capacity <= 0 {
		return nil, errors.New("capacity must be greater than 0")
	}
	if room.Floor <= 0 {
		return nil, errors.New("floor must be greater than 0")
	}
	if len(room.Amenities) == 0 { // Fix: Check slice length, not empty string
		return nil, errors.New("amenities are required")
	}

	fmt.Println("Service: Validation passed, calling repository...")

	if err := s.repo.AddRoom(ctx, room); err != nil {
		fmt.Printf("Service: Repository error - %v\n", err)
		return nil, fmt.Errorf("failed to add room: %w", err)
	}

	fmt.Printf("Service: Room added successfully - ID: %d\n", room.ID)
	return room, nil
}

func (s *RoomService) GetRoomsList(ctx context.Context) ([]*models.Room, error) {
	rooms, err := s.repo.GetRoomsList(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *RoomService) GetAvailableRooms(ctx context.Context) ([]*models.Room, error) {
	cacheKey := "available_rooms" // Fixed: removed fmt.Sprintln

	// Check if Redis client is initialized
	if cache.Client == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}

	cached, err := cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var rooms []*models.Room
		if err := json.Unmarshal([]byte(cached), &rooms); err == nil {
			fmt.Println("✅ Cache hit - available rooms")
			return rooms, nil
		}
	}

	// If cache miss or error, proceed to database
	fmt.Println("❌ Cache miss - fetching available rooms from database")

	rooms, err := s.repo.GetAvailableRooms(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result
	roomsJSON, err := json.Marshal(rooms)
	if err == nil {
		err = cache.Client.Set(ctx, cacheKey, roomsJSON, 10*time.Minute).Err()
		if err != nil {
			fmt.Printf("⚠️ Failed to cache rooms: %v\n", err)
		} else {
			fmt.Println("✅ Available rooms cached successfully")
		}
	}

	return rooms, nil
}
