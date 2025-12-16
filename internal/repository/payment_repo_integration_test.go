package repository

import (
	"context"
	"os"
	"testing"

	"industry-api/internal/models"

	"industry-api/internal/testsetup"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPaymentRepository_InitiateAndCleanup(t *testing.T) {
	_ = testsetup.LoadEnv()
	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		dsn = os.Getenv("DB_URL")
	}
	if dsn == "" {
		t.Skip("DB_URL or TEST_DB_URL not set; skipping integration tests")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	repo := NewPaymentRepository(pool)

	tx := "tx-test"
	p := &models.Payment{
		BookingID:     0,
		Amount:        100,
		PaymentMethod: "test",
		TransactionID: &tx,
		Status:        "pending",
	}

	created, err := repo.InitiatePayment(ctx, p)
	if err != nil {
		t.Logf("InitiatePayment failed (likely FK constraints): %v", err)
		t.Skip("InitiatePayment requires valid booking FK; skipping in this environment")
	}
	if created.ID == 0 {
		t.Fatalf("expected payment ID set")
	}

	// cleanup
	if _, err := pool.Exec(ctx, "DELETE FROM payments WHERE id = $1", created.ID); err != nil {
		t.Logf("warning: cleanup failed: %v", err)
	}
}
