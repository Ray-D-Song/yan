package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func RegisterLifecycle(lc fx.Lifecycle, engine *gin.Engine) *http.Server {
	srv := &http.Server{
		Addr:    ":18080",
		Handler: engine,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
