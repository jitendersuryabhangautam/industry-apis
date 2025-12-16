// Tests replaced by integration tests. This file kept to avoid accidental
// references to pgxmock in environments where the mock is incompatible with
// the project's pgx version. See user_repo_integration_test.go for repo tests.
package repository

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestUserRepo_Skipped(t *testing.T) {
	// Load .env to allow local docker-compose to supply DB_URL
	_ = godotenv.Load()
	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		dsn = os.Getenv("DB_URL")
	}
	if dsn == "" {
		t.Skip("unit mock tests disabled; use integration tests or reintroduce pgxmock compatible version")
		return
	}

	// If DB is configured, run the integration test
	TestUserRepository_CreateAndGet(t)
}
