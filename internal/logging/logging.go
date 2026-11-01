package logging

import (
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(app config.AppConfig) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if app.Env == "development" {
		cfg = zap.NewDevelopmentConfig()
	}
	level := zapcore.InfoLevel
	switch app.LogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(level)
	return cfg.Build()
}
