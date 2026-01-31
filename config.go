package goauth

import (
	"database/sql"
	"time"
)

// Config holds the configuration for the auth service
type Config struct {
	// Database connection
	DB *sql.DB
	
	// Table name for users (default: "outcome_user")
	TableName string
	
	// JWT secret key
	JWTSecret string
	
	// Access token duration (default: 2 hours)
	AccessTokenDuration time.Duration
	
	// Refresh token duration (default: 30 days)
	RefreshTokenDuration time.Duration
	
	// Guest ID range for random generation (default: 1-100000)
	GuestIDMin int
	GuestIDMax int
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig(db *sql.DB, jwtSecret string) *Config {
	return &Config{
		DB:                   db,
		TableName:            "outcome_user",
		JWTSecret:            jwtSecret,
		AccessTokenDuration:  2 * time.Hour,
		RefreshTokenDuration: 30 * 24 * time.Hour,
		GuestIDMin:           1,
		GuestIDMax:           100000,
	}
}
