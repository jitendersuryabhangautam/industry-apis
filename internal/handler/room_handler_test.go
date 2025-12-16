package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"industry-api/internal/models"
	"industry-api/internal/service"

	"github.com/gin-gonic/gin"
)

type mockRoomSvcRepo struct {
	add       func(ctx context.Context, room *models.Room) error
	list      func(ctx context.Context) ([]*models.Room, error)
	available func(ctx context.Context) ([]*models.Room, error)
}

func (m *mockRoomSvcRepo) AddRoom(ctx context.Context, room *models.Room) error {
	return m.add(ctx, room)
}
func (m *mockRoomSvcRepo) GetRoomsList(ctx context.Context) ([]*models.Room, error) {
	return m.list(ctx)
}
func (m *mockRoomSvcRepo) GetAvailableRooms(ctx context.Context) ([]*models.Room, error) {
	return m.available(ctx)
}

func TestAddRoomHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := &mockRoomSvcRepo{
		add:       func(ctx context.Context, room *models.Room) error { room.ID = 1; return nil },
		list:      func(ctx context.Context) ([]*models.Room, error) { return nil, errors.New("not-impl") },
		available: func(ctx context.Context) ([]*models.Room, error) { return nil, errors.New("not-impl") },
	}
	h := NewRoomHandler(service.NewRoomService(mr))

	reqBody := models.RoomRequest{RoomNumber: "101", RoomType: "Deluxe", Description: "test room", Price: 100, Capacity: 2, Floor: 1, Amenities: []string{"wifi"}}
	b, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/rooms/add", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AddRoom(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestGetRoomsListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := &mockRoomSvcRepo{
		add: func(ctx context.Context, room *models.Room) error { return errors.New("not-impl") },
		list: func(ctx context.Context) ([]*models.Room, error) {
			return []*models.Room{{ID: 1, RoomNumber: "101"}}, nil
		},
		available: func(ctx context.Context) ([]*models.Room, error) { return nil, errors.New("not-impl") },
	}
	h := NewRoomHandler(service.NewRoomService(mr))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/rooms/allRoomsList", nil)
	h.GetRoomsList(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
