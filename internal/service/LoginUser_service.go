package service

import (
	"context"
	"errors"
	"industry-api/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) LoginUser(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	token, err := generateJWT(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	user.Password = ""
	return &models.LoginResponse{
		User:      user,
		Token:     token,
		TokenType: "Bearer",       // Add token type
		ExpiresIn: time.Hour * 24, // Add expiration info
	}, nil
}

func generateJWT(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(), // issued at
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
