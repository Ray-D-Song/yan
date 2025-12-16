// Package router provides route registration functions for HTTP endpoints.
package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ray-d-song/yan/internal/api/v1"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(rg *gin.RouterGroup, userHandler *v1.UserHandler) {
	users := rg.Group("/users")
	{
		// Public routes
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.Login)

		// User management routes
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateProfile)
		users.PUT("/:id/password", userHandler.ChangePassword)
	}
}
