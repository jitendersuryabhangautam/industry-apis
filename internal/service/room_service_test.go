package service

import (
	"context"
	"errors"
	"testing"

	"industry-api/internal/models"
)

// mockRoomRepo implements only what RoomService needs
type mockRoomRepo struct {
	add       func(ctx context.Context, room *models.Room) error
	list      func(ctx context.Context) ([]*models.Room, error)
	available func(ctx context.Context) ([]*models.Room, error)
}

func (m *mockRoomRepo) AddRoom(ctx context.Context, room *models.Room) error     { return m.add(ctx, room) }
func (m *mockRoomRepo) GetRoomsList(ctx context.Context) ([]*models.Room, error) { return m.list(ctx) }
func (m *mockRoomRepo) GetAvailableRooms(ctx context.Context) ([]*models.Room, error) {
	return m.available(ctx)
}

func TestAddRoom_Validation(t *testing.T) {
	svc := &RoomService{repo: &mockRoomRepo{}}

	tests := []struct {
		name string
		room *models.Room
	}{
		{"missing number", &models.Room{RoomNumber: "", RoomType: "Deluxe", Description: "d", Price: 100, Capacity: 2, Floor: 1, Amenities: []string{"a"}}},
		{"missing type", &models.Room{RoomNumber: "101", RoomType: "", Description: "d", Price: 100, Capacity: 2, Floor: 1, Amenities: []string{"a"}}},
		{"bad price", &models.Room{RoomNumber: "101", RoomType: "A", Description: "d", Price: 0, Capacity: 2, Floor: 1, Amenities: []string{"a"}}},
		{"no amenities", &models.Room{RoomNumber: "101", RoomType: "A", Description: "d", Price: 100, Capacity: 2, Floor: 1, Amenities: []string{}}},
	}

	for _, tt := range tests {
		if _, err := svc.AddRoom(context.Background(), tt.room); err == nil {
			t.Fatalf("%s: expected validation error", tt.name)
		}
	}
}

func TestAddRoom_Success(t *testing.T) {
	repo := &mockRoomRepo{
		add: func(ctx context.Context, room *models.Room) error {
			room.ID = 777
			return nil
		},
		list:      func(ctx context.Context) ([]*models.Room, error) { return nil, errors.New("not-implemented") },
		available: func(ctx context.Context) ([]*models.Room, error) { return nil, errors.New("not-implemented") },
	}
	svc := &RoomService{repo: repo}
	room := &models.Room{RoomNumber: "R1", RoomType: "T", Description: "desc", Price: 50, Capacity: 2, Floor: 1, Amenities: []string{"wifi"}}
	created, err := svc.AddRoom(context.Background(), room)
	if err != nil {
		t.Fatalf("AddRoom returned error: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected ID set by repo")
	}
}

func TestGetRoomsList_DelegatesToRepo(t *testing.T) {
	sample := []*models.Room{{ID: 1, RoomNumber: "1"}}
	repo := &mockRoomRepo{list: func(ctx context.Context) ([]*models.Room, error) { return sample, nil }}
	svc := &RoomService{repo: repo}
	got, err := svc.GetRoomsList(context.Background())
	if err != nil {
		t.Fatalf("GetRoomsList error: %v", err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("unexpected result from GetRoomsList")
	}
}
