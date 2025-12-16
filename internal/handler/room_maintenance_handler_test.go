package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"industry-api/internal/models"
	"industry-api/internal/service"

	"github.com/gin-gonic/gin"
)

type mockRoomMaintSvcRepo struct {
	add func(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error)
}

func (m *mockRoomMaintSvcRepo) AddRoomMaintenance(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error) {
	return m.add(ctx, rm)
}

func TestAddRoomMaintenanceHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := &mockRoomMaintSvcRepo{
		add: func(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error) {
			rm.ID = 1
			return rm, nil
		},
	}
	h := NewRoomMaintenanceHandler(service.NewRoomMaintenanceService(mr))

	reqBody := models.RoomMaintenanceRequest{
		RoomID:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(2 * time.Hour),
		Reason:    "inspection",
		Status:    "scheduled",
		CreatedBy: 1,
	}
	b, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/roomMaintenance/add", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AddRoomMaintenance(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAddRoomMaintenanceHandler_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := &mockRoomMaintSvcRepo{
		add: func(ctx context.Context, rm *models.RoomMaintenance) (*models.RoomMaintenance, error) { return rm, nil },
	}
	h := NewRoomMaintenanceHandler(service.NewRoomMaintenanceService(mr))

	// missing required fields
	reqBody := models.RoomMaintenanceRequest{RoomID: 0}
	b, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/roomMaintenance/add", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AddRoomMaintenance(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
