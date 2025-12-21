// Package v1 provides HTTP handlers for API version 1 endpoints.
package v1

import "github.com/gin-gonic/gin"

// RouteRegistrar defines the interface for handlers that can register their own routes
type RouteRegistrar interface {
	RegisterRoutes(g *gin.RouterGroup)
}
