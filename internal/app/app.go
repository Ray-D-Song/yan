// Package app provides the application initialization and dependency injection setup using Uber FX.
package app

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ray-d-song/yan/internal/api/v1"
	"github.com/ray-d-song/yan/internal/infra"
	"github.com/ray-d-song/yan/internal/mdw"
	"github.com/ray-d-song/yan/internal/repo"
	"github.com/ray-d-song/yan/internal/service"
	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			// infra
			infra.LoadConfig,
			infra.NewLogger,
			infra.NewDB,
			infra.NewGin,
			infra.NewAPIV1Group,

			// session
			repo.NewSessionRepo,
			infra.NewSessionStore,

			// repo
			repo.NewUserRepo,
			repo.NewNoteRepo,

			// service
			service.NewUserService,
			service.NewNoteService,

			// handler
			v1.NewUserHandler,
			v1.NewNoteHandler,
		),
		fx.Invoke(
			RegisterLifecycle,
			RegisterRoutes,
		),
	)
}

// RegisterRoutes registers all application routes with appropriate middleware
func RegisterRoutes(
	apiV1 *gin.RouterGroup,
	userHandler *v1.UserHandler,
	noteHandler *v1.NoteHandler,
	store *infra.DBStore,
	userService service.UserService,
) {
	// Register user routes (public endpoints)
	userHandler.RegisterRoutes(apiV1)

	// Create auth middleware
	authMiddleware := mdw.AuthMiddleware(store, userService)

	// Register note routes with auth protection
	notesGroup := apiV1.Group("/notes")
	notesGroup.Use(authMiddleware)
	noteHandler.RegisterRoutes(notesGroup)
}
