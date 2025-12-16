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

func TestRoomsRepo_CreateAndList(t *testing.T) {
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

	repo := NewRoomRepository(pool)

	ts := time.Now().UnixNano()
	roomNum := "T-" + time.Unix(0, ts).Format("150405")
	room := &models.Room{
		RoomNumber:  roomNum,
		RoomType:    "Test",
		Description: "Integration test room",
		Price:       10,
		Capacity:    1,
		Floor:       1,
		Amenities:   []string{"test"},
	}

	if err := repo.AddRoom(ctx, room); err != nil {
		t.Fatalf("AddRoom failed: %v", err)
	}
	if room.ID == 0 {
		t.Fatalf("expected room ID to be set")
	}

	rooms, err := repo.GetRoomsList(ctx)
	if err != nil {
		t.Fatalf("GetRoomsList failed: %v", err)
	}
	found := false
	for _, r := range rooms {
		if r.RoomNumber == roomNum {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created room not found in list")
	}

	// cleanup
	if _, err := pool.Exec(ctx, "DELETE FROM rooms WHERE room_number = $1", roomNum); err != nil {
		t.Logf("warning: cleanup failed: %v", err)
	}
}
