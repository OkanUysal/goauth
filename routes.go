package goauth

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all auth routes to the provided router
// baseURL: e.g., "/api/v1/auth"
func (s *Service) RegisterRoutes(router *gin.Engine, baseURL string) {
	authGroup := router.Group(baseURL)
	{
		// Public routes
		authGroup.POST("/guest", s.GuestLogin)
		authGroup.POST("/refresh", s.RefreshToken)
		
		// Protected routes
		authGroup.GET("/profile", s.AuthMiddleware(), s.GetProfile)
	}
}

// RegisterRoutesWithGroup registers routes to an existing router group
func (s *Service) RegisterRoutesWithGroup(group *gin.RouterGroup) {
	// Public routes
	group.POST("/guest", s.GuestLogin)
	group.POST("/refresh", s.RefreshToken)
	
	// Protected routes
	group.GET("/profile", s.AuthMiddleware(), s.GetProfile)
}
