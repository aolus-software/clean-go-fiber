package postgres

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate runs all pending up migrations. It is a no-op when the schema is
// already up to date. DB must be connected before calling this.
func Migrate() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("migrate: get sql.DB: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("migrate: create iofs source: %w", err)
	}

	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate: create postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate: init migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate: apply migrations: %w", err)
	}

	return nil
}
