package auth

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Service handles authentication operations
type Service struct {
	config *Config
}

// NewService creates a new auth service
func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

// GuestLogin creates a new guest user and returns tokens
func (s *Service) GuestLogin(c *gin.Context) {
	// Generate unique guest ID
	guestID := uuid.New()
	userID := uuid.New().String()

	// Generate random display name
	randomNum := rand.Intn(s.config.GuestIDMax-s.config.GuestIDMin+1) + s.config.GuestIDMin
	displayName := fmt.Sprintf("Guest%d", randomNum)

	now := time.Now()

	// Insert new user
	query := fmt.Sprintf(`
		INSERT INTO %s (id, guest_id, display_name, role, is_guest, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, guest_id, display_name, role, is_guest, created_at, updated_at
	`, s.config.TableName)

	var user User
	err := s.config.DB.QueryRow(
		query,
		userID, guestID, displayName, RoleUser, true, now, now,
	).Scan(
		&user.ID, &user.GuestID, &user.DisplayName,
		&user.Role, &user.IsGuest, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Role, s.config.JWTSecret, s.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken(user.ID, s.config.JWTSecret, s.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// RefreshToken validates refresh token and issues new tokens
func (s *Service) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate refresh token
	claims, err := ValidateRefreshToken(req.RefreshToken, s.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Get user from database
	query := fmt.Sprintf(`
		SELECT id, guest_id, google_id, display_name, role, is_guest, created_at, updated_at
		FROM %s WHERE id = $1
	`, s.config.TableName)

	var user User
	err = s.config.DB.QueryRow(query, claims.UserID).Scan(
		&user.ID, &user.GuestID, &user.GoogleID, &user.DisplayName,
		&user.Role, &user.IsGuest, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Update updated_at timestamp
	updateQuery := fmt.Sprintf(`UPDATE %s SET updated_at = $1 WHERE id = $2`, s.config.TableName)
	_, err = s.config.DB.Exec(updateQuery, time.Now(), user.ID)
	if err != nil {
		// Log but don't fail the request
		fmt.Printf("Warning: Failed to update updated_at: %v\n", err)
	}

	// Generate new tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Role, s.config.JWTSecret, s.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	newRefreshToken, err := GenerateRefreshToken(user.ID, s.config.JWTSecret, s.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	})
}

// GetProfile returns the authenticated user's profile
func (s *Service) GetProfile(c *gin.Context) {
	userID, exists := GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user from database
	query := fmt.Sprintf(`
		SELECT id, guest_id, google_id, display_name, role, is_guest, created_at, updated_at
		FROM %s WHERE id = $1
	`, s.config.TableName)

	var user User
	err := s.config.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.GuestID, &user.GoogleID, &user.DisplayName,
		&user.Role, &user.IsGuest, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
