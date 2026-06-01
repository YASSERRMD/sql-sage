package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func RunMigrations(gdb *gorm.DB, path string, log *slog.Logger) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}
	return runMigrationsWith(sqlDB, path, log)
}

func runMigrationsWith(sqlDB *sql.DB, path string, log *slog.Logger) error {
	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("migrate driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+path, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate new: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %w", err)
	}
	if log != nil {
		log.Info("database migrations applied")
	}
	return nil
}
