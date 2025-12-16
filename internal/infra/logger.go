package infra

import (
	"go.uber.org/zap"
)

func NewLogger(config *Config) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if err := cfg.Level.UnmarshalText([]byte(config.Log.Level)); err != nil {
		return nil, err
	}
	return cfg.Build()
}
