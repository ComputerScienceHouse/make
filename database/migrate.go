package database

import (
	"database/sql"
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(db *sql.DB) error {
	log.Println("[DB] Beginning database migration")

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Printf("[DB] Failed to load migrations: %v\n", err)
		return err
	}

	log.Println("[DB] Running migrations...")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("[DB] Failed to create postgres driver: %v\n", err)
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		source,
		"postgres",
		driver,
	)
	if err != nil {
		log.Printf("[DB] Failed to create migrator: %v\n", err)
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("[DB] Database already up to date")
		return nil
	}

	if err != nil {
		log.Printf("[DB] Migration failed: %v\n", err)
		return err
	}

	version, dirty, verr := m.Version()
	if verr == nil {
		log.Printf("[DB] Migration complete. Current version: %d (dirty=%v)\n", version, dirty)
	} else {
		log.Println("[DB] Migration complete.")
	}

	return nil
}
