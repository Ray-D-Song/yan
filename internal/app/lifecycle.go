package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/ray-d-song/yan/internal/infra"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RegisterLifecycle(lc fx.Lifecycle, engine *gin.Engine, db *sqlx.DB, logger *infra.Logger) *http.Server {
	srv := &http.Server{
		Addr:    ":18080",
		Handler: engine,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Run database migrations
			if err := infra.AutoMigrate(db, logger); err != nil {
				logger.Error("Failed to run database migrations", zap.Error(err))
				return err
			}

			// Start HTTP server
			go srv.ListenAndServe()
			logger.Info("Server started on :18080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down server...")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
