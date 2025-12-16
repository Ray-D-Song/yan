package app

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ray-d-song/yan/internal/api/v1"
	"github.com/ray-d-song/yan/internal/router"
	"go.uber.org/fx"
)

// Handlers aggregates all HTTP handlers
type Handlers struct {
	fx.In

	User *v1.UserHandler
}

// RegisterRoutes registers all API routes
func RegisterRoutes(engine *gin.Engine, handlers Handlers) {
	apiV1 := engine.Group("/api/v1")
	router.RegisterUserRoutes(apiV1, handlers.User)
}
