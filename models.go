package auth

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents user role types
type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

// User represents a user in the system
type User struct {
	ID          string    `json:"id" db:"id"`
	GuestID     uuid.UUID `json:"guest_id" db:"guest_id"`
	GoogleID    *string   `json:"google_id,omitempty" db:"google_id"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Role        UserRole  `json:"role" db:"role"`
	IsGuest     bool      `json:"is_guest" db:"is_guest"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// LoginRequest represents a login request body
type LoginRequest struct {
	// For guest login, this can be empty
	// For Google login, this would contain the Google token
	Token string `json:"token,omitempty"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
