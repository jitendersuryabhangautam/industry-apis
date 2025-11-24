package service

import (
	"context"
	"errors"
	"industry-api/internal/models"
	"industry-api/internal/repository"
	"strings"

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
