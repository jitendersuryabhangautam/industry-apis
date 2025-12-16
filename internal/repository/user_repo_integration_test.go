package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"industry-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Integration tests for UserRepository. These run against the DB specified by
// TEST_DB_URL or DB_URL. Tests create and then delete their test rows to keep
// the same DB usable for other development tasks.
func TestUserRepository_CreateAndGet(t *testing.T) {
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

	repo := NewUserRepository(pool)

	// unique email to avoid conflicts
	ts := time.Now().UnixNano()
	email := "test+" + time.Unix(0, ts).Format("20060102150405") + "@example.com"

	user := &models.User{
		Name:     "Integration Tester",
		Email:    email,
		Password: "hash-placeholder",
		Phone:    "0000000000",
		Role:     "guest",
	}

	// Create
	if err := repo.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if user.ID == "" {
		t.Fatalf("expected created user to have ID assigned")
	}

	// Fetch by email
	fetched, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if fetched.Email != email {
		t.Fatalf("expected email %s, got %s", email, fetched.Email)
	}

	// Cleanup: delete created user
	if _, err := pool.Exec(ctx, "DELETE FROM users WHERE email = $1", email); err != nil {
		t.Logf("warning: failed to cleanup test user: %v", err)
	}
}
