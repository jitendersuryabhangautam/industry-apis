package service

import (
	"os"
	"testing"

	"industry-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	// Ensure JWT_SECRET is set for token generation
	os.Setenv("JWT_SECRET", "test-secret")

	user := &models.User{
		ID:    "42",
		Email: "tester@example.com",
		Role:  "admin",
	}

	tokenStr, err := generateJWT(user)
	if err != nil {
		t.Fatalf("generateJWT returned error: %v", err)
	}
	if tokenStr == "" {
		t.Fatalf("expected non-empty token string")
	}

	// Parse token and validate claims
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !parsed.Valid {
		t.Fatalf("token is not valid")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected MapClaims, got %T", parsed.Claims)
	}
	if claims["email"] != user.Email {
		t.Fatalf("expected email claim %s, got %v", user.Email, claims["email"])
	}
}
