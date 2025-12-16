// Tests replaced by integration tests. This file kept to avoid accidental
// references to pgxmock in environments where the mock is incompatible with
// the project's pgx version. See user_repo_integration_test.go for repo tests.
package repository

import "testing"

func TestUserRepo_Skipped(t *testing.T) {
	t.Skip("unit mock tests disabled; use integration tests or reintroduce pgxmock compatible version")
}
