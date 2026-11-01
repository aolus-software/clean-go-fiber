package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate() error {
	m, err := newMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: apply migrations: %w", err)
	}
	return nil
}

func MigrateDown() error {
	m, err := newMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: roll back migrations: %w", err)
	}
	return nil
}

func newMigrator() (*migrate.Migrate, error) {
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("migrate: get sql.DB: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("migrate: create iofs source: %w", err)
	}

	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		return nil, fmt.Errorf("migrate: get connection: %w", err)
	}

	driver, err := migratePostgres.WithConnection(context.Background(), conn, &migratePostgres.Config{})
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("migrate: create postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		_ = driver.Close()
		return nil, fmt.Errorf("migrate: init migrator: %w", err)
	}
	return m, nil
}
