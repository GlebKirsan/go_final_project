package database

import (
	"database/sql"

	"github.com/GlebKirsan/go-final-project/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
)

func runMigration(db *sql.DB) error {
	cfg := config.Get()

	if cfg.MigrationPath == "" {
		return nil
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	f := file.File{}
	d, err := f.Open("file://" + cfg.MigrationPath)
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithInstance(
		"file", d, "sqlite3", driver,
	)

	if err != nil {
		return err
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
