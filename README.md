# GoAuth - Reusable Authentication Library for Go

A flexible, plug-and-play authentication library for Go web applications using Gin framework.

## Features

- 🔐 Guest Login with JWT tokens
- 🔄 Refresh Token rotation
- 👤 User Profile management
- 🛡️ Role-based access control (USER/ADMIN)
- 🔌 Easy integration - just provide DB connection and config
- 📦 Zero boilerplate - routes are automatically registered

## Installation

```bash
go get github.com/OkanUysal/goauth
```

## Quick Start

### 1. Database Setup

Create a users table in your PostgreSQL database:

```sql
CREATE TABLE outcome_user (
    id UUID PRIMARY KEY,
    guest_id UUID UNIQUE NOT NULL,
    google_id TEXT UNIQUE,
    display_name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'USER',
    is_guest BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 2. Initialize the Library

```go
package main

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "github.com/OkanUysal/goauth"
    _ "github.com/lib/pq"
)

func main() {
    // Connect to database
    db, _ := sql.Open("postgres", "your-connection-string")
    
    // Create auth service with default config
    authConfig := goauth.DefaultConfig(db, "your-jwt-secret-key")
    authService := goauth.NewService(authConfig)
    
    // Setup Gin router
    router := gin.Default()
    
    // Register auth routes automatically
    authService.RegisterRoutes(router, "/api/v1/auth")
    
    // Your other routes...
    
    router.Run(":8080")
}
```

### 3. Use the Endpoints

**Guest Login:**
```bash
POST /api/v1/auth/guest
```

Response:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "display_name": "Guest12345",
    "role": "USER",
    "is_guest": true
  }
}
```

**Refresh Token:**
```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

**Get Profile:**
```bash
GET /api/v1/auth/profile
Authorization: Bearer eyJhbGc...
```

## Configuration

### Custom Configuration

```go
authConfig := &goauth.Config{
    DB:                   db,
    TableName:            "my_users_table",  // Custom table name
    JWTSecret:            "my-secret",
    AccessTokenDuration:  1 * time.Hour,     // Custom duration
    RefreshTokenDuration: 7 * 24 * time.Hour,
    GuestIDMin:           1000,
    GuestIDMax:           9999,
}
```

## Middleware Usage

### Protect Routes with Authentication

```go
// Require authentication
protectedRoutes := router.Group("/api/v1")
protectedRoutes.Use(authService.AuthMiddleware())
{
    protectedRoutes.GET("/dashboard", dashboardHandler)
}
```

### Require Admin Role

```go
// Require admin role
adminRoutes := router.Group("/api/v1/admin")
adminRoutes.Use(authService.AuthMiddleware(), authService.AdminMiddleware())
{
    adminRoutes.GET("/users", listUsersHandler)
}
```

### Access User Info in Handlers

```go
func myHandler(c *gin.Context) {
    userID, _ := goauth.GetUserID(c)
    userRole, _ := goauth.GetUserRole(c)
    
    // Use userID and userRole...
}
```

## Advanced Usage

### Register Routes to Existing Group

```go
api := router.Group("/api/v1")
authGroup := api.Group("/auth")
authService.RegisterRoutesWithGroup(authGroup)
```

### Custom Table Name

If you want to use a different table name:

```go
authConfig := goauth.DefaultConfig(db, jwtSecret)
authConfig.TableName = "users" // Use "users" instead of "outcome_user"
```

## API Reference

### Types

- `Config`: Configuration for the auth service
- `Service`: Main authentication service
- `User`: User model
- `UserRole`: Role type (USER/ADMIN)
- `LoginResponse`: Login response with tokens
- `RefreshTokenRequest`: Refresh token request

### Functions

- `NewService(config *Config) *Service`: Create new auth service
- `DefaultConfig(db *sql.DB, jwtSecret string) *Config`: Get default config
- `RegisterRoutes(router *gin.Engine, baseURL string)`: Register routes
- `AuthMiddleware() gin.HandlerFunc`: Authentication middleware
- `AdminMiddleware() gin.HandlerFunc`: Admin authorization middleware
- `GetUserID(c *gin.Context) (string, bool)`: Get user ID from context
- `GetUserRole(c *gin.Context) (UserRole, bool)`: Get user role from context

## License

MIT
