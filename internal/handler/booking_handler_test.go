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

type mockBookingSvcRepo struct {
	add func(ctx context.Context, b *models.Booking) (*models.Booking, error)
}

func (m *mockBookingSvcRepo) AddBooking(ctx context.Context, b *models.Booking) (*models.Booking, error) {
	return m.add(ctx, b)
}

func TestAddBookingHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := &mockBookingSvcRepo{
		add: func(ctx context.Context, b *models.Booking) (*models.Booking, error) { b.ID = 1; return b, nil },
	}
	h := NewBookingHandler(service.NewBookingService(mr))

	reqBody := models.BookingRequest{
		UserID:        1,
		RoomID:        1,
		CheckInDate:   time.Now(),
		CheckOutDate:  time.Now().Add(24 * time.Hour),
		Adults:        2,
		Children:      1,
		TotalAmount:   200,
		Status:        "pending",
		PaymentStatus: "pending",
	}
	b, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/bookings/add", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AddBooking(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAddBookingHandler_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := &mockBookingSvcRepo{
		add: func(ctx context.Context, b *models.Booking) (*models.Booking, error) { return b, nil },
	}
	h := NewBookingHandler(service.NewBookingService(mr))

	// missing required fields
	reqBody := models.BookingRequest{UserID: 0, RoomID: 0}
	b, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/bookings/add", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AddBooking(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
