package infra

import (
	"go.uber.org/zap"
)

// Logger wraps zap.Logger to provide a more convenient API
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger(config *Config) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	if err := cfg.Level.UnmarshalText([]byte(config.Log.Level)); err != nil {
		return nil, err
	}
	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}

// Convenience methods for common logging patterns

// Infof logs a formatted info message
func (l *Logger) Infof(template string, args ...interface{}) {
	l.Sugar().Infof(template, args...)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Sugar().Errorf(template, args...)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.Sugar().Debugf(template, args...)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Sugar().Warnf(template, args...)
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.With(fields...)}
}

// WithContext returns a logger with context information
func (l *Logger) WithContext(key string, value string) *Logger {
	return &Logger{Logger: l.With(zap.String(key, value))}
}
