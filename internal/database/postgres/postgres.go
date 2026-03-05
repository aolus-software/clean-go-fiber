package postgres

import (
	"github.com/aolus-software/clean-go-fiber/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	databseUrl := config.C.Database.URL
	logLevel := logger.Silent
	if config.C.App.Env == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(databseUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, err
	}

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	return nil, err
	// }

	DB = db

	return DB, nil
}
