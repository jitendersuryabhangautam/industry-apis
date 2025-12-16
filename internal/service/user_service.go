// Package service provides business logic layer implementations.
// Services contain validation logic and business rules before delegating to repositories.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"industry-api/internal/cache"
	"industry-api/internal/models"
	"industry-api/internal/repository"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles all user-related business logic operations.
// It includes user creation, retrieval, caching, and password management.
type UserService struct {
	repo repository.UserRepo // Repository interface for data access (allows mocking in tests)
}

// NewUserService creates and returns a new instance of UserService.
// It accepts a UserRepository dependency for data access operations.
func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user account after validation and password hashing.
// It validates all required fields, checks for duplicate emails, hashes the password,
// and delegates to the repository for persistence.
// Returns the created user without the password hash or an error if validation fails.
func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Validate required fields are not empty
	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		return nil, errors.New("all fields required")
	}
	// Check if user with this email already exists
	existing, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user with this email already exists")
	}
	// Hash the password for secure storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("falied to hash password")
	}

	// Store hashed password instead of plain text
	user.Password = string(hashedPassword)
	// Delegate to repository to persist the user
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	// Clear password before returning to client for security
	user.Password = ""
	return user, nil

}

// GetUserList retrieves a paginated list of users with optional filtering.
// It delegates to the repository for database query execution.
func (s *UserService) GetUserList(ctx context.Context, role string, isActive *bool, search string, page, limit int) (*models.UserListResponse, error) {
	return s.repo.GetUserList(ctx, role, isActive, search, page, limit)
}

// GetUserByID retrieves a user by ID with Redis caching.
// It first tries to fetch from cache, and if not found (cache miss),
// fetches from the database, caches the result, and returns it.
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	// Create a unique cache key for this user
	cacheKey := fmt.Sprintf("user:%d", id)
	// Try to get user from cache
	cachedUser, err := cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		// Deserialize cached JSON into user model
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			fmt.Println("cache hit- user details")
			return &user, nil
		}
	}

	// Cache miss - fetch from database
	fmt.Println("cache miss- fetching from database")
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Create a copy of the user for caching (without password)
	userToCache := *user
	userToCache.Password = "" // Remove password before caching
	// Serialize user to JSON for caching
	userJSON, err := json.Marshal(userToCache)
	if err == nil {
		// Cache the user data for 10 minutes
		err = cache.Client.Set(ctx, cacheKey, userJSON, 10*time.Minute).Err()
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

// UpdateUserStatus updates a user's active status.
// It validates the user exists before updating and invalidates the cache.
// Returns the updated user or an error if the user is not found.
func (s *UserService) UpdateUserStatus(ctx context.Context, id int, isActive bool) (*models.User, error) {
	// Verify user exists before attempting update
	_, err := s.GetUserByID(ctx, id)
	if err != nil {
		return nil, err // This will already be "user not found" if applicable
	}
	// Delegate to repository to update the status
	user, err := s.repo.UpdateUserStatus(ctx, id, isActive)
	if err != nil {
		return nil, err
	}
	// Invalidate cache for this user since it has been updated
	s.invalidateUserCache(ctx, id)

	return user, nil
}

// invalidateUserCache removes cached data for a specific user and user lists.
// This is called whenever user data is modified to ensure fresh data on next fetch.
func (s *UserService) invalidateUserCache(ctx context.Context, userID int) {
	// Invalidate the specific user's cache
	cacheKey := fmt.Sprintf("user:%d", userID)
	cache.Client.Del(ctx, cacheKey)

	// Also invalidate user lists if you have them cached
	userListCacheKey := "users:list"
	cache.Client.Del(ctx, userListCacheKey)
}
