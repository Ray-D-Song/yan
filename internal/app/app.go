// Package app provides the application initialization and dependency injection setup using Uber FX.
package app

import (
	v1 "github.com/ray-d-song/yan/internal/api/v1"
	"github.com/ray-d-song/yan/internal/infra"
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
			infra.NewRouter,

			// repo
			repo.NewUserRepo,

			// service
			service.NewUserService,

			// handler
			v1.NewUserHandler,
		),
		fx.Invoke(
			RegisterRoutes,
			RegisterLifecycle,
		),
	)
}

