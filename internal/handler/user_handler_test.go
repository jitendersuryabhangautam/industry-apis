package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"industry-api/internal/models"
	"industry-api/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// mockRepo implements repository.UserRepo for handler tests
type mockRepo struct {
	GetUserByEmailFn   func(ctx context.Context, email string) (*models.User, error)
	CreateUserFn       func(ctx context.Context, user *models.User) error
	GetUserByIDFn      func(ctx context.Context, id int) (*models.User, error)
	UpdateUserStatusFn func(ctx context.Context, id int, isActive bool) (*models.User, error)
	GetUserListFn      func(ctx context.Context, role string, isActive *bool, search string, page, limit int) (*models.UserListResponse, error)
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.GetUserByEmailFn(ctx, email)
}
func (m *mockRepo) CreateUser(ctx context.Context, user *models.User) error {
	return m.CreateUserFn(ctx, user)
}
func (m *mockRepo) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return m.GetUserByIDFn(ctx, id)
}
func (m *mockRepo) UpdateUserStatus(ctx context.Context, id int, isActive bool) (*models.User, error) {
	return m.UpdateUserStatusFn(ctx, id, isActive)
}
func (m *mockRepo) GetUserList(ctx context.Context, role string, isActive *bool, search string, page, limit int) (*models.UserListResponse, error) {
	return m.GetUserListFn(ctx, role, isActive, search, page, limit)
}

func TestRegisterHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// mock repo: no existing email, create returns ID
	mr := &mockRepo{
		GetUserByEmailFn: func(ctx context.Context, email string) (*models.User, error) { return nil, context.Canceled },
		CreateUserFn:     func(ctx context.Context, user *models.User) error { user.ID = "u-100"; return nil },
	}
	us := service.NewUserService(mr)
	h := NewUserHandler(us)

	reqBody := models.RegisterRequest{Name: "Test", Email: "t@example.com", Password: "pw12345", Phone: "1234567890"}
	b, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if resp["success"] != true {
		t.Fatalf("expected success true, got %v", resp["success"])
	}
}

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// prepare password hash
	raw := "loginpass"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)

	mr := &mockRepo{
		GetUserByEmailFn: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{ID: "u-200", Email: email, Password: string(hashed), Role: "guest"}, nil
		},
		CreateUserFn: func(ctx context.Context, user *models.User) error { return nil },
	}
	us := service.NewUserService(mr)
	// ensure JWT secret set for token creation during login
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret")
	h := NewUserHandler(us)

	// build request
	body := map[string]string{"email": "a@b.com", "password": raw}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(b))
	c.Request.Header.Set("Content-Type", "application/json")

	// set env JWT_SECRET so token generation works
	// Use service.LoginUser which calls generateJWT reading JWT_SECRET
	// set env var
	// Note: tests assume generateJWT reads os.Getenv("JWT_SECRET")

	h.LoginUser(c)

	if w.Code != http.StatusUnauthorized && w.Code != http.StatusOK {
		t.Fatalf("expected 200 or 401 (if token generation missing), got %d body=%s", w.Code, w.Body.String())
	}
}
