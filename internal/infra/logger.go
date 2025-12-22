package infra

import (
	"sync"

	"go.uber.org/zap"
)

var (
	globalLogger *Logger
	loggerMu     sync.RWMutex
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

	logger := &Logger{Logger: zapLogger}

	// Set as global logger
	SetGlobalLogger(logger)

	return logger, nil
}

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	globalLogger = logger
}

// L returns the global logger instance
func L() *Logger {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	if globalLogger == nil {
		// Return a no-op logger if not initialized
		return &Logger{Logger: zap.NewNop()}
	}
	return globalLogger
}

// Global convenience functions that use the global logger

// Infof logs a formatted info message using the global logger
func Infof(template string, args ...interface{}) {
	L().Infof(template, args...)
}

// Errorf logs a formatted error message using the global logger
func Errorf(template string, args ...interface{}) {
	L().Errorf(template, args...)
}

// Debugf logs a formatted debug message using the global logger
func Debugf(template string, args ...interface{}) {
	L().Debugf(template, args...)
}

// Warnf logs a formatted warning message using the global logger
func Warnf(template string, args ...interface{}) {
	L().Warnf(template, args...)
}

// Instance methods

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
