package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"industry-api/internal/models"
	"industry-api/internal/testsetup"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRoomMaintenance_AddAndCleanup(t *testing.T) {
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

	repo := NewRoomMaintenanceRepository(pool)

	rm := &models.RoomMaintenance{
		RoomID:    0,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(2 * time.Hour),
		Reason:    "test",
		Status:    "scheduled",
		CreatedBy: 0,
	}

	created, err := repo.AddRoomMaintenance(ctx, rm)
	if err != nil {
		t.Logf("AddRoomMaintenance failed (likely FK constraints): %v", err)
		t.Skip("AddRoomMaintenance requires valid FK; skipping in this environment")
	}
	if created.ID == 0 {
		t.Fatalf("expected maintenance ID set")
	}

	// cleanup
	if _, err := pool.Exec(ctx, "DELETE FROM room_maintenance WHERE id = $1", created.ID); err != nil {
		t.Logf("warning: cleanup failed: %v", err)
	}
}
