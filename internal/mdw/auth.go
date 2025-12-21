// Package mdw provides middleware functions for the Gin framework.
package mdw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/ray-d-song/yan/internal/infra"
	"github.com/ray-d-song/yan/internal/service"
)

const (
	SessionName = "yan_session"
	SessionKeyUserID = "user_id"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(store sessions.Store, userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session
		session, err := store.Get(c.Request, SessionName)
		if err != nil {
			c.String(http.StatusUnauthorized, "invalid session")
			c.Abort()
			return
		}

		// Get user ID from session
		userID, ok := session.Values[SessionKeyUserID].(int64)
		if !ok || userID == 0 {
			c.String(http.StatusUnauthorized, "not authenticated")
			c.Abort()
			return
		}

		// Get user from database
		user, err := userService.GetByID(c.Request.Context(), userID)
		if err != nil {
			c.String(http.StatusUnauthorized, "user not found")
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive() {
			c.String(http.StatusForbidden, "user account is disabled")
			c.Abort()
			return
		}

		// Store user in context
		infra.SetUserInContext(c, user)

		c.Next()
	}
}

// OptionalAuthMiddleware is like AuthMiddleware but doesn't abort if no auth
// Useful for endpoints that work with or without authentication
func OptionalAuthMiddleware(store sessions.Store, userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session
		session, err := store.Get(c.Request, SessionName)
		if err != nil {
			c.Next()
			return
		}

		// Get user ID from session
		userID, ok := session.Values[SessionKeyUserID].(int64)
		if !ok || userID == 0 {
			c.Next()
			return
		}

		// Get user from database
		user, err := userService.GetByID(c.Request.Context(), userID)
		if err != nil {
			c.Next()
			return
		}

		// Check if user is active
		if !user.IsActive() {
			c.Next()
			return
		}

		// Store user in context
		infra.SetUserInContext(c, user)

		c.Next()
	}
}

// RequireAdminMiddleware checks if the authenticated user is an admin
// Must be used after AuthMiddleware
func RequireAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := infra.UserFromCtx(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "not authenticated")
			c.Abort()
			return
		}

		if !user.IsAdministrator() {
			c.String(http.StatusForbidden, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
