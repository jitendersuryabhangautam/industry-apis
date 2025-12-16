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

type mockPaymentSvcRepo struct {
	init   func(ctx context.Context, p *models.Payment) (*models.Payment, error)
	update func(ctx context.Context, p *models.Payment, id int) (*models.Payment, error)
}

func (m *mockPaymentSvcRepo) InitiatePayment(ctx context.Context, p *models.Payment) (*models.Payment, error) {
	return m.init(ctx, p)
}
func (m *mockPaymentSvcRepo) UpdatePayment(ctx context.Context, p *models.Payment, id int) (*models.Payment, error) {
	return m.update(ctx, p, id)
}

func TestInitiatePaymentHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := &mockPaymentSvcRepo{
		init: func(ctx context.Context, p *models.Payment) (*models.Payment, error) { p.ID = 1; return p, nil },
		update: func(ctx context.Context, p *models.Payment, id int) (*models.Payment, error) {
			return nil, errors.New("not-impl")
		},
	}
	h := NewPaymentHandler(service.NewPaymentService(mr))

	reqBody := models.PaymentRequest{BookingID: 1, Amount: 100, PaymentMethod: "card", TransactionID: "tx-1"}
	b, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/payments/initiate", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	h.InitiatePayment(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestInitiatePaymentHandler_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mr := &mockPaymentSvcRepo{
		init: func(ctx context.Context, p *models.Payment) (*models.Payment, error) { return p, nil },
		update: func(ctx context.Context, p *models.Payment, id int) (*models.Payment, error) {
			return nil, errors.New("not-impl")
		},
	}
	h := NewPaymentHandler(service.NewPaymentService(mr))

	// missing required fields
	reqBody := models.PaymentRequest{BookingID: 0, Amount: 0}
	b, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/payments/initiate", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	h.InitiatePayment(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
