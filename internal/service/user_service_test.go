package service

import (
	"context"
	"errors"
	"os"
	"testing"

	"industry-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo implements the minimal subset of repository methods used by UserService
type mockUserRepo struct {
	getByEmail func(ctx context.Context, email string) (*models.User, error)
	create     func(ctx context.Context, user *models.User) error
}

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.getByEmail(ctx, email)
}
func (m *mockUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return m.create(ctx, user)
}
func (m *mockUserRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return nil, errors.New("not-implemented")
}
func (m *mockUserRepo) UpdateUserStatus(ctx context.Context, id int, isActive bool) (*models.User, error) {
	return nil, errors.New("not-implemented")
}
func (m *mockUserRepo) GetUserList(ctx context.Context, role string, isActive *bool, search string, page, limit int) (*models.UserListResponse, error) {
	return nil, errors.New("not-implemented")
}

func TestCreateUser_Success(t *testing.T) {
	// repo returns "not found" on email lookup, and sets ID on create
	repo := &mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (*models.User, error) {
			return nil, errors.New("not found")
		},
		create: func(ctx context.Context, user *models.User) error {
			user.ID = "u-123"
			return nil
		},
	}

	svc := &UserService{repo: repo}

	user := &models.User{Name: "Alice", Email: "a@example.com", Password: "pass123", Phone: "1234567890"}
	created, err := svc.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected ID to be set")
	}
	if created.Password != "" {
		t.Fatalf("expected password to be cleared, got %q", created.Password)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repo := &mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{ID: "existing"}, nil
		},
		create: func(ctx context.Context, user *models.User) error { return nil },
	}
	svc := &UserService{repo: repo}
	_, err := svc.CreateUser(context.Background(), &models.User{Name: "B", Email: "b@example.com", Password: "pw"})
	if err == nil {
		t.Fatalf("expected error for duplicate email")
	}
}

func TestLoginUser_SuccessAndFailure(t *testing.T) {
	// prepare hashed password
	raw := "supersecret"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)

	repo := &mockUserRepo{
		getByEmail: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{ID: "u1", Email: email, Password: string(hashed), Role: "admin"}, nil
		},
		create: func(ctx context.Context, user *models.User) error { return nil },
	}
	svc := &UserService{repo: repo}

	// ensure JWT secret exists for token generation
	os.Setenv("JWT_SECRET", "test-secret")

	// successful login
	lr, err := svc.LoginUser(context.Background(), "x@example.com", raw)
	if err != nil {
		t.Fatalf("expected successful login, got error: %v", err)
	}
	if lr.Token == "" {
		t.Fatalf("expected non-empty token")
	}

	// invalid password
	_, err = svc.LoginUser(context.Background(), "x@example.com", "wrong")
	if err == nil {
		t.Fatalf("expected error for invalid password")
	}
}
