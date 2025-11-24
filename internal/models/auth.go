package models

import "time"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=255"`
}

type LoginResponse struct {
	User      *User         `json:"user"`
	Token     string        `json:"access_token"`
	TokenType string        `json:"token_type"`
	ExpiresIn time.Duration `json:"expires_in"`
}
