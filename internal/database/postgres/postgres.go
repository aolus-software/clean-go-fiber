package postgres

import (
	"log"
	"os"
	"time"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	cfg := config.C.Database

	logLevel := logger.Warn
	if config.C.App.Env == "development" && config.C.App.LogLevel == "debug" {
		logLevel = logger.Info
	}
	gormLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             cfg.SlowQuery,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		Colorful:                  config.C.App.Env == "development",
	})

	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{
		Logger:         gormLogger,
		TranslateError: true,
		NowFunc:        func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	DB = db

	return DB, nil
}
