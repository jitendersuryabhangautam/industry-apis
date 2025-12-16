// Package service provides business logic layer implementations.
// This file contains user authentication logic including password verification and JWT token generation.
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

// LoginUser authenticates a user by verifying email and password credentials.
// It validates the credentials, generates a JWT token, and returns the user with token information.
// Returns a LoginResponse containing user details and access token, or an error if authentication fails.
func (s *UserService) LoginUser(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	// Retrieve user by email from database
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	// Verify the provided password matches the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	// Generate JWT token for authenticated session
	token, err := generateJWT(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Clear password before sending to client
	user.Password = ""
	// Return user details with JWT token and expiration info
	return &models.LoginResponse{
		User:      user,
		Token:     token,
		TokenType: "Bearer",       // Token type for Authorization header
		ExpiresIn: time.Hour * 24, // Token valid for 24 hours
	}, nil
}

// generateJWT creates a signed JWT token with user claims.
// It reads the JWT secret from environment variables and creates claims with user data and expiration time.
// Returns the signed token string or an error if token generation fails.
func generateJWT(user *models.User) (string, error) {
	// Get the JWT signing secret from environment variables
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	// Create JWT claims with user information
	claims := jwt.MapClaims{
		"user_id": user.ID,                               // Unique user identifier
		"email":   user.Email,                            // User's email
		"role":    user.Role,                             // User's role for authorization
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expiration time (24 hours)
		"iat":     time.Now().Unix(),                     // Issued at timestamp
	}
	// Create token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key and return
	return token.SignedString([]byte(secret))
}
