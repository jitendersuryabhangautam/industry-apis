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

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if strings.TrimSpace(user.Name) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		return nil, errors.New("all fields required")
	}
	existing, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user with this email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("falied to hash password")
	}

	user.Password = string(hashedPassword)
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil

}

func (s *UserService) GetUserList(ctx context.Context, role string, isActive *bool, search string, page, limit int) (*models.UserListResponse, error) {
	return s.repo.GetUserList(ctx, role, isActive, search, page, limit)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {

	cacheKey := fmt.Sprintf("user:%d", id)
	cachedUser, err := cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			fmt.Println("cache hit- user details")
			return &user, nil
		}
	}

	fmt.Println("cache miss- fetching from database")
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	userToCache := *user
	userToCache.Password = ""
	userJSON, err := json.Marshal(userToCache)
	if err == nil {
		err = cache.Client.Set(ctx, cacheKey, userJSON, 10*time.Minute).Err()
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, id int, isActive bool) (*models.User, error) {
	_, err := s.GetUserByID(ctx, id)
	if err != nil {
		return nil, err // This will already be "user not found" if applicable
	}
	user, err := s.repo.UpdateUserStatus(ctx, id, isActive)
	if err != nil {
		return nil, err
	}
	s.invalidateUserCache(ctx, id)

	return user, nil
}
func (s *UserService) invalidateUserCache(ctx context.Context, userID int) {
	cacheKey := fmt.Sprintf("user:%d", userID)
	cache.Client.Del(ctx, cacheKey)

	// Also invalidate user lists if you have them cached
	userListCacheKey := "users:list"
	cache.Client.Del(ctx, userListCacheKey)
}
