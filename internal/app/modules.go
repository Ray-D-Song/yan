package app

import (
	"github.com/ray-d-song/yan/internal/log"
	"github.com/ray-d-song/yan/internal/repo"
	"github.com/ray-d-song/yan/internal/service"
	"go.uber.org/fx"
)

func Modules() fx.Option {
	return fx.Options(
		// infra
		fx.Provide(log.NewLogger),

		// db
		fx.Provide(NewDB),

		// repo
		fx.Provide(repo.NewUserRepo),

		// service
		fx.Provide(service.NewUserService),
	)
}
