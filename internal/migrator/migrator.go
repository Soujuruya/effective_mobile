package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Migrator struct {
	srcDriver source.Driver
}

func MustGetNewMigrator(sqlFiles embed.FS, dirName string) *Migrator {

	d, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		panic(err)
	}
	return &Migrator{
		srcDriver: d,
	}
}

func (m *Migrator) ApplyMigrations(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("unable to open database for migration: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %v", err)
	}

	migrator, err := migrate.NewWithInstance("embed_migrations", m.srcDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}
	defer migrator.Close()

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations: %v", err)
	}

	return nil
}
